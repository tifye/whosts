package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tifye/whosts/pkg"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			hostsFile, err := cmd.Flags().GetString("hosts")
			if err != nil {
				return err
			}

			file, err := os.Open(hostsFile)
			if err != nil {
				return fmt.Errorf("open file: %s", err)
			}
			defer file.Close()

			hosts, err := pkg.ParseEntries(file)
			if err != nil {
				return err
			}

			fmt.Println(hosts.String())
			return nil
		},
	}
	return cmd
}
