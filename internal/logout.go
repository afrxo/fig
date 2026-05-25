package internal

import (
	"fmt"

	"github.com/afrxo/fig/auth"
)

var LogoutCommand = Command{
	Name:    "logout",
	Usage:   fmt.Sprintf("%s logout", CLI),
	Summary: "Clear stored credentials",
	Run: func(flags CommandFlags, args []string) error {
		if err := auth.Delete(); err != nil {
			return fmt.Errorf("failed to remove credentials: %w", err)
		}

		fmt.Println("✓ Logged out.")

		return nil
	},
}
