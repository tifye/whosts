package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var (
	defaultHostsPath = "C:\\windows\\system32\\drivers\\etc\\hosts"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "whosts",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
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
	root.AddCommand(newListCommand())
}
