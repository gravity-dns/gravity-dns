package main

import (
	"os"

	"github.com/gravity-dns/gravity-dns/cmd/ctl/commands"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
