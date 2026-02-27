package main

import (
	"os"

	"github.com/angolovin/yougile-cli/internal/errors"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(errors.ExitCodeError)
	}
}
