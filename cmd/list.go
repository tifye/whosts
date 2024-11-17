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
			file, err := os.Open(pkg.DefaultHostsPath)
			if err != nil {
				return fmt.Errorf("open file: %s", err)
			}
			defer file.Close()

			entries, err := pkg.ParseEntries(file)
			if err != nil {
				return err
			}

			for _, e := range entries {
				fmt.Printf("%s %s\n", e.IP.String(), e.Host)
			}

			return nil
		},
	}
	return cmd
}