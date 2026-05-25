// Package internal/commands.go
package internal

import (
	"fmt"
	"os"
	"slices"
	"sort"
)

type Command struct {
	Name    string
	Aliases []string
	Summary string
	Usage   string
	Run     func(flags CommandFlags, args []string) error
}

type CommandFlags map[string]string

var Commands = []Command{}

func Run() {
	flags, values := Parse(os.Args)

	values = values[1:]

	if v := GetFirstFlag(flags, "false", "V", "version"); v != "false" {
		fmt.Printf("%s %s\n", CLI, Version)
		os.Exit(0)
		return
	}

	argCount, flagCount := len(values), len(flags)

	if (argCount == 0 && flagCount == 0) || flags["help"] == "true" {
		PrintBanner()
		os.Exit(0)
		return
	}

	target := ""
	if argCount != 0 {
		target = values[0]
	}

	command := FindCommand(target)

	needsHelp := GetFirstFlag(flags, "", "h", "help")

	if command == nil {
		command = FindCommand(needsHelp)
	}

	if command == nil {
		fmt.Printf("%s %s: unknown command\n", CLI, needsHelp)
		fmt.Printf("Run '%s help' for usage.\n", CLI)
		os.Exit(0)
		return
	} else if command != nil && needsHelp != "" {
		DisplayHelp(command)
		os.Exit(0)
		return
	}

	err := command.Run(flags, values)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}
}

func RegisterCommands(cmds *[]Command) {
	Commands = *cmds
}

func FindCommand(name string) *Command {
	for i := range Commands {
		if Commands[i].Name == name {
			return &Commands[i]
		}
		if slices.Contains(Commands[i].Aliases, name) {
			return &Commands[i]
		}
	}
	return nil
}

func PrintBanner() {
	fmt.Fprintln(os.Stdout, CLI, Version)
	fmt.Println("A Go CLI for vendoring packages from any Git repo.")
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, ("USAGE:"))
	fmt.Printf("    %s <SUBCOMMAND> [OPTIONS]\n", CLI)
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, ("OPTIONS:"))
	fmt.Fprintln(os.Stdout, "    -h, --help       Print help information")
	fmt.Fprintln(os.Stdout, "    -V, --version    Print version information")
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "SUBCOMMANDS:")
	names := make([]string, 0, len(Commands))
	width := 0
	for _, c := range Commands {
		names = append(names, c.Name)
		if len(c.Name) > width {
			width = len(c.Name)
		}
	}
	sort.Strings(names)
	for _, n := range names {
		c := FindCommand(n)
		fmt.Fprintf(os.Stdout, "    %s    %s\n", padRight(c.Name, width), c.Summary)
	}
	fmt.Fprintln(os.Stdout)
	fmt.Printf("Run `%s help <subcommand>` for more information on a specific subcommand.\n", CLI)
}
