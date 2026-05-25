package internal

import (
	"fmt"

	"github.com/afrxo/fig/auth"
)

var UserCommand = Command{
	Name:    "user",
	Summary: "Show the current authenticated user",
	Run: func(flags CommandFlags, values []string) error {
		creds, err := auth.Load()
		if err != nil {
			fmt.Printf("Not logged in, use %s login\n", CLI)
			return nil
		}

		fmt.Printf("✓ Logged in as %s\n", creds.Username)

		return nil
	},
}
