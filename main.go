package main

import (
	"fmt"
	"os"

	"github.com/thatnerdjosh/kvit/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "kvit: %v\n", err)
		os.Exit(1)
	}
}
