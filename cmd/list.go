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

			for _, e := range hosts.Entries() {
				fmt.Printf("%s %s\n", e.IP.String(), e.Host)
			}

			return nil
		},
	}
	return cmd
}
