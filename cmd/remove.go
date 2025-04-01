package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/tifye/whosts/pkg"
)

type removeOptions struct {
	ip             net.IP
	host           string
	comment        string
	noComment      bool
	duplicatesOnly bool
	dryRun         bool
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

			// Avoid any accidents.
			// Could also ensure filemode but too lazy to look that up
			var perms int
			if opts.dryRun {
				perms = os.O_RDONLY
			} else {
				perms = os.O_RDWR
			}
			file, err := os.OpenFile(hostsFile, perms, 0755)
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
			if opts.duplicatesOnly && len(filters) == 0 {
				filters = append(filters, pkg.WithAll())
			}

			removed := hosts.Remove(opts.duplicatesOnly, filters...)

			if !opts.dryRun {
				terr := file.Truncate(0)
				if terr != nil {
					return fmt.Errorf("truncate: %s", terr)
				}
				_, serr := file.Seek(0, 0)
				if serr != nil {
					return fmt.Errorf("seek: %s", serr)
				}
				if _, err = hosts.WriteTo(file); err != nil {
					return err
				}
			}

			fmt.Printf("Updated:\n%s\n\nRemoved:\n%s", hosts, pkg.NewHosts(removed).String())

			return nil
		},
	}

	cmd.Flags().IPVar(&opts.ip, "ip", nil, "Remove entries with matching IP")
	cmd.Flags().StringVar(&opts.host, "host", "", "Remove entries with matching host name")
	cmd.Flags().StringVar(&opts.comment, "comment", "", "Remove entries with matching comment")
	cmd.Flags().BoolVar(&opts.noComment, "no-comment", false, "Remove entries without comments")
	cmd.Flags().BoolVar(&opts.duplicatesOnly, "duplicates-only", false, "Remove entry duplicates that match passed filters. If no filters are passed then remove any duplicate.")
	cmd.MarkFlagsOneRequired("ip", "host", "comment")

	cmd.Flags().BoolVar(&opts.dryRun, "dry", false, "Dry run command and print out which entries would have been removed")

	return cmd
}
