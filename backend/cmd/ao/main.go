package main

import (
	"fmt"
	"os"

	"github.com/sriram687/GIT_WIZ/backend/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitCode(err))
	}
}
