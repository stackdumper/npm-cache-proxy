package cli

import (
	"fmt"
	"os"
)

// Run starts the CLI
func Run() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(purgeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
