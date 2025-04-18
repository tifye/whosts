package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func newOpenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Opens the hosts file in notepad",
		RunE: func(cmd *cobra.Command, args []string) error {
			hostsFile, err := cmd.Flags().GetString("hosts")
			if err != nil {
				return err
			}

			execCmd := exec.Command("notepad", hostsFile)
			err = execCmd.Start()
			if err, ok := err.(*exec.ExitError); ok {
				return fmt.Errorf("failed to open: %s", err)
			}

			err = execCmd.Process.Release()
			if err != nil {
				return fmt.Errorf("err proc release: %s", err)
			}

			return nil
		},
	}

	return cmd
}
