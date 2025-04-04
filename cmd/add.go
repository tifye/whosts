package cmd

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/net/idna"

	"github.com/spf13/cobra"
	"github.com/tifye/whosts/pkg"
)

type addOptions struct {
	ip   net.IP
	host string
}

func newAddCommand() *cobra.Command {
	opts := addOptions{}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an entry",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(2)(cmd, args); err != nil {
				return fmt.Errorf("accepts 2 positional args: <ip> <host>")
			}

			ipStr, host := args[0], args[1]
			ip := net.ParseIP(ipStr)
			if ip == nil {
				return fmt.Errorf("provided IP is not a valid textual represetation of an IP address")
			}

			_, err := idna.Lookup.ToASCII(host)
			if err != nil {
				return fmt.Errorf("invalid host name format, expected format as defined in RFC 5891")
			}

			opts.ip = ip
			opts.host = host
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			hostsFile, err := cmd.Flags().GetString("hosts")
			if err != nil {
				return err
			}

			file, err := os.OpenFile(hostsFile, os.O_RDWR, 0755)
			if err != nil {
				return fmt.Errorf("open file: %s", err)
			}
			defer file.Close()

			hosts, err := pkg.ParseEntries(file)
			if err != nil {
				return err
			}

			hosts.AddEntry(pkg.Entry{
				IP:   opts.ip,
				Host: opts.host,
			})

			if _, err = hosts.WriteTo(file); err != nil {
				return err
			}

			fmt.Println(hosts.String())

			return nil
		},
	}
	return cmd
}
