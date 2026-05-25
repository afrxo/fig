package internal

import (
	"fmt"
)

var PickCommand = Command{
	Name:    "pick",
	Summary: "Add a package from a remote repository",
	Usage:   fmt.Sprintf(`%s add <name> <repo> [-o out] [-p path] [-r ref]`, CLI),
	Run:     runPick,
}

func runPick(flags CommandFlags, args []string) error {
	args = args[1:]

	argCount := len(args)

	if argCount < 2 {
		return fmt.Errorf("expected command usage: %s add <name> <repo> [-o out] [-p path] [-r ref]", CLI)
	}

	name := args[0]

	ref := GetFirstFlag(flags, "main", "r", "ref")
	path := GetFirstFlag(flags, fmt.Sprintf("packages/%s", name), "p", "path")
	out := GetFirstFlag(flags, fmt.Sprintf("src/%s", name), "o", "out")

	target := TargetPackage{
		Name: args[0],
		Repo: args[1],
		Out:  out,
		Ref:  ref,
		Path: path,
	}

	lockfile, err := ReadLockfile()
	if err != nil {
		return err
	}

	err = SyncRefs(&lockfile, []TargetPackage{target})
	if err != nil {
		return err
	}

	return err
}
