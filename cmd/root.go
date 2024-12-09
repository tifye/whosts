package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/tifye/whosts/pkg"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "whosts",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.PersistentFlags().String("hosts", pkg.DefaultHostsPath, "Path to hosts file to target")
	cmd.MarkPersistentFlagFilename("hosts")

	return cmd
}

func Execute() {
	root := newRootCommand()
	addCommands(root)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := root.ExecuteContext(ctx)
	if err != nil {
		fmt.Println(err)
	}
}

func addCommands(root *cobra.Command) {
	root.AddCommand(
		newListCommand(),
		newDumpCommand(),
		newAddCommand(),
		newOpenCommand(),
	)
}
