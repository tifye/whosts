package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := root.ExecuteContext(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
