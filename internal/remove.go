package internal

import (
	"fmt"
	"os"
	"path"
	"time"
)

var RemoveCommand = Command{
	Name:    "remove",
	Usage:   fmt.Sprintf("%s remove <...package>", CLI),
	Summary: "Remove a vendored package",
	Run: func(flags CommandFlags, args []string) error {
		if len(args) <= 1 {
			return fmt.Errorf("expected command usage: %s remove <...package>", CLI)
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		args = args[1:]

		start := time.Now()

		lockfile, err := ReadLockfile()
		if err != nil {
			return err
		}

		if len(lockfile.Packages) <= 0 {
			fmt.Printf("%s is empty, did you forget to run %s add?\n", LockfileName, CLI)
			return nil
		}

		fmt.Printf("Removing packages...\n\n")

		var count int

		for _, name := range args {
			pkg, ok := lockfile.Packages[name]

			if !ok {
				return fmt.Errorf("could not find package '%s'\n\nRun '%s list' to see a list of vendored packages", name, CLI)
			}

			if err := os.RemoveAll(path.Join(wd, pkg.Out)); err != nil {
				return err
			}

			delete(lockfile.Packages, name)

			fmt.Printf("Removed %s\n", name)

			count++
		}

		fmt.Printf("\nRemoved %d packages in %v\n", count, time.Since(start))

		if count > 0 {
			lockfile.Synced = GetFormattedTime()
			err := WriteLockfile(&lockfile)
			if err != nil {
				return err
			}
		}

		return nil
	},
}
