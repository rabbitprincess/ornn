package gen

import (
	"os"

	"github.com/gokch/ornn/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
