package internal

import (
	"fmt"
)

var Version = "dev"

var VersionCommand = Command{
	Name:    "version",
	Usage:   fmt.Sprintf("%s version", CLI),
	Summary: fmt.Sprintf("Outputs %ss current version", CLI),
	Run: func(flags CommandFlags, args []string) error {
		fmt.Printf("%s %s\n", CLI, Version)
		return nil
	},
}
