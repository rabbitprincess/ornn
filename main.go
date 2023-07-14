package main

import (
	"os"

	"github.com/gokch/ornn/cmd/ornn"
)

func main() {
	if err := ornn.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
