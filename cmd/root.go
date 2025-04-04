package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

	root.SetHelpFunc(rootHelp)

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
		newRemoveCommand(),
	)
}

func rootHelp(cmd *cobra.Command, _ []string) {
	fmt.Println("\nUsage")
	fmt.Printf("  %s [command]\n", cmd.Use)

	fmt.Println("\nAvailable Commands:")
	for _, c := range cmd.Commands() {
		fmt.Printf("  %-10s %s\n", c.Name(), c.Short)
		c.Flags().VisitAll(func(f *pflag.Flag) {
			fmt.Printf("    --%-10s\t%s\n", f.Name, f.Usage)
		})
	}

	fmt.Printf("\nUse \"%s [command] --help\" for more information about a command.", cmd.Name())
}
