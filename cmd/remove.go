package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/tifye/whosts/pkg"
)

type removeOptions struct {
	ip        net.IP
	host      string
	comment   string
	noComment bool
	dryRun    bool
}

func newRemoveCommand() *cobra.Command {
	opts := &removeOptions{}
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove entries matching passed filters. Filters are stacked.",
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

			filters := make([]pkg.FilterOption, 0)
			if opts.ip != nil {
				filters = append(filters, pkg.WithIPs(opts.ip))
			}
			if opts.host != "" {
				filters = append(filters, pkg.WithHosts(opts.host))
			}
			if opts.comment != "" {
				filters = append(filters, pkg.WithComments(opts.comment))
			}
			if opts.noComment {
				filters = append(filters, pkg.WithNoComment())
			}

			removed := hosts.Remove(filters...)

			if !opts.dryRun {
				panic("not implemented")
			}

			fmt.Printf("Updated:\n%s\n\nRemoved:\n%s", hosts, pkg.NewHosts(removed).String())

			return nil
		},
	}

	cmd.Flags().IPVar(&opts.ip, "ip", nil, "Remove entries with matching IP")
	cmd.Flags().StringVar(&opts.host, "host", "", "Remove entries with matching host name")
	cmd.Flags().StringVar(&opts.comment, "comment", "", "Remove entries with matching comment")
	cmd.Flags().BoolVar(&opts.noComment, "no-comment", false, "Remove entries without comments")
	cmd.MarkFlagsOneRequired("ip", "host", "comment")

	cmd.Flags().BoolVar(&opts.dryRun, "dry", false, "Dry run command and print out which entries would have been removed")

	return cmd
}
