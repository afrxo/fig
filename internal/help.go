package internal

import (
	"fmt"
	"os"
	"strings"
)

var HelpCommand = Command{
	Name:    "help",
	Summary: "Display help for a command",
	Run:     runHelp,
}

func runHelp(flags CommandFlags, args []string) error {
	args = args[2:]

	if len(args) <= 0 {
		PrintBanner()
		return nil
	}

	command := args[0]

	for _, cmd := range Commands {
		if cmd.Name != command {
			continue
		}

		DisplayHelp(&cmd)
		return nil
	}

	fmt.Printf("%s: unknown command\n", command)

	var available []string

	for _, cmd := range Commands {
		available = append(available, cmd.Name)
	}

	fmt.Printf("list of available commands:\n\n%s\n", strings.Join(available, ", "))

	return nil
}

func DisplayHelp(cmd *Command) {
	fmt.Fprintln(os.Stdout, cmd.Summary)
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, cmd.Usage)
}

func padRight(s string, n int) string {
	for len(s) < n {
		s += " "
	}
	return s
}
