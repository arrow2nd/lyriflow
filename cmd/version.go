package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func ShowVersion(version string) cli.ActionFunc {
	return func(c *cli.Context) error {
		fmt.Printf("lyriflow version %s\n", version)
		return nil
	}
}
