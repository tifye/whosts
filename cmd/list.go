package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tifye/whosts/pkg"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			bytes, err := os.ReadFile(pkg.DefaultHostsPath)
			if err != nil {
				return fmt.Errorf("read file: %s", err)
			}

			fmt.Println(string(bytes))

			return nil
		},
	}
	return cmd
}
