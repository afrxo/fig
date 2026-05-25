package internal

import (
	"fmt"
	"strings"
)

var ListCommand = Command{
	Name:    "list",
	Usage:   fmt.Sprintf("%s list", CLI),
	Summary: "List all vendored packages",
	Run: func(flags CommandFlags, args []string) error {
		lockfile, err := ReadLockfile()
		if err != nil {
			return err
		}

		if lockfile.Synced == "" || len(lockfile.Packages) == 0 {
			fmt.Println("No packages installed yet")
			return nil
		}

		t, err := ParseTime(lockfile.Synced)
		if err != nil {
			return err
		}

		fmt.Printf("Last updated at %v\n", t)

		c := []string{}

		for name, pkg := range lockfile.Packages {
			c = append(c, fmt.Sprintf("%s %s %s", name, pkg.Repo, pkg.Ref))
		}

		fmt.Println()
		fmt.Println(strings.Join(c, "\n"))

		return nil
	},
}
