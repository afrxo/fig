package internal

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/afrxo/fig/auth"

	"github.com/otiai10/copy"
)

var SyncCommand = Command{
	Name:    "sync",
	Summary: "Sync packages from the lockfile",
	Usage:   fmt.Sprintf(`%s sync [...package]`, CLI),
	Run:     runSync,
}

type RefPackageMap map[string][]TargetPackage

type TargetPackage struct {
	Ref  string
	Out  string
	Path string
	Repo string
	Name string
}

func SyncRefs(lockfile *Lockfile, packages []TargetPackage) error {
	fmt.Printf("adding packages...\n\n")

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	creds, err := auth.Load()
	if err != nil {
		return err
	}

	dname, err := os.MkdirTemp("", ".gkit")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dname)

	repos := map[string]RefPackageMap{}
	var pkgCount int
	start := time.Now()

	for _, pkg := range packages {
		// Reconcile repo scope
		_, ok := repos[pkg.Repo]
		if !ok {
			repos[pkg.Repo] = RefPackageMap{}
		}

		// Reconcile ref scope
		_, ok = repos[pkg.Repo][pkg.Ref]
		if !ok {
			repos[pkg.Repo][pkg.Ref] = []TargetPackage{}
		}

		repos[pkg.Repo][pkg.Ref] = append(repos[pkg.Repo][pkg.Ref], pkg)
	}

	for r, refs := range repos {
		repoPath := path.Join(fmt.Sprintf(dname, r))
		repoURL := fmt.Sprintf("https://github.com/%s.git", r)

		repo, err := CloneRepository(repoURL, repoPath, creds)
		if err != nil {
			return err
		}

		for ref, pkgs := range refs {
			err := CheckoutRepository(repo, ref)
			if err != nil {
				return err
			}

			headRef, err := repo.Head()
			if err != nil {
				return err
			}

			sha := headRef.Hash().String()

			w, err := repo.Worktree()
			if err != nil {
				return err
			}

			fs := w.Filesystem()

			for _, pkg := range pkgs {
				p := path.Join(fs.Root(), pkg.Path)
				dest := path.Join(wd, pkg.Out)

				_, err = os.Stat(p)
				if err != nil {
					return err
				}

				err = copy.Copy(p, dest)
				if err != nil {
					return err
				}

				if lockfile.Packages == nil {
					lockfile.Packages = make(map[string]Package)
				}

				lockfile.Packages[pkg.Name] = Package{
					Ref:  ref,
					Out:  pkg.Out,
					Sha:  sha,
					Path: pkg.Path,
					Repo: r,
				}

				pkgCount++

				fmt.Printf("added %s %s from %s\n", pkg.Name, pkg.Ref, pkg.Repo)
			}
		}
	}

	lockfile.Synced = GetFormattedTime()

	if err := WriteLockfile(lockfile); err != nil {
		return err
	}

	fmt.Printf("\nadded %d packages in %v\n", pkgCount, time.Since(start))

	return nil
}

func runSync(flags CommandFlags, args []string) error {
	args = args[2:]

	lockfile, err := ReadLockfile()
	if err != nil {
		return err
	}

	packages := []TargetPackage{}

	if len(args) <= 0 {
		for name, pkg := range lockfile.Packages {
			packages = append(packages, TargetPackage{
				Path: pkg.Path,
				Repo: pkg.Repo,
				Out:  pkg.Out,
				Ref:  pkg.Ref,
				Name: name,
			})
		}
	} else {
		for _, name := range args {
			pkg, ok := lockfile.Packages[name]
			if !ok {
				return fmt.Errorf("%s: unknown package", name)
			}

			packages = append(packages, TargetPackage{
				Path: pkg.Path,
				Repo: pkg.Repo,
				Out:  pkg.Out,
				Ref:  pkg.Ref,
				Name: name,
			})
		}
	}

	return SyncRefs(&lockfile, packages)
}
