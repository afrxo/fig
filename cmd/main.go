// Package cmd/fig/main.go
package main

import (
	"github.com/afrxo/fig/internal"
)

func main() {
	internal.RegisterCommands(&[]internal.Command{
		internal.PickCommand,
		internal.SyncCommand,
		internal.HelpCommand,
		internal.UserCommand,
		internal.ListCommand,
		internal.LoginCommand,
		internal.RemoveCommand,
		internal.LogoutCommand,
		internal.VersionCommand,
	})

	internal.Run()
}
