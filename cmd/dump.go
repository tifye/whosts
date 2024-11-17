package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tifye/whosts/pkg"
)

func newDumpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dump",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := os.Open(pkg.DefaultHostsPath)
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
