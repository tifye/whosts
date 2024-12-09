package cmd

import (
	"fmt"
	"net"

	"golang.org/x/net/idna"

	"github.com/spf13/cobra"
)

type addOptions struct {
	ip   net.IP
	host string
}

func newAddCommand() *cobra.Command {
	opts := addOptions{}
	cmd := &cobra.Command{
		Use: "add",
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
			fmt.Println("add", opts.ip, opts.host)
			return nil
		},
	}
	return cmd
}
