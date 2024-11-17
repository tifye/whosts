package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDumpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dump",
		RunE: func(cmd *cobra.Command, args []string) error {
			hostsFile, err := cmd.Flags().GetString("hosts")
			if err != nil {
				return err
			}

			file, err := os.Open(hostsFile)
			if err != nil {
				return fmt.Errorf("open file: %s", err)
			}

			buf := bufio.NewReader(file)
			_, err = buf.WriteTo(os.Stdout)
			if err != nil {
				return fmt.Errorf("err reading file: %s", err)
			}

			return nil
		},
	}

	return cmd
}
