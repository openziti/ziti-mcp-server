package main

import (
	"os"

	"github.com/openziti/ziti-mcp-server-go/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
