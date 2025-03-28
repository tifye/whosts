package pkg

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	t.Run("empty file: parse without error", func(t *testing.T) {
		hosts, err := ParseEntries(strings.NewReader(""))
		if assert.NoError(t, err) {
			assert.Len(t, hosts.entries, 0)
		}
	})

	t.Run("empty lines: parse without error", func(t *testing.T) {
		hosts, err := ParseEntries(strings.NewReader("\n\n\n\n\n         \n           \n\n"))
		if assert.NoError(t, err) {
			assert.Len(t, hosts.entries, 0)
		}
	})

	t.Run("entry followed by newline", func(t *testing.T) {
		hosts, err := ParseEntries(strings.NewReader(
			`
			109.94.209.70   fitgirl-repack.net      # Fake FitGirl site


			`,
		))
		if assert.NoError(t, err) {
			require.Len(t, hosts.entries, 1)
			assert.Equal(t, "# Fake FitGirl site", hosts.entries[0].Comment)
			assert.Equal(t, "fitgirl-repack.net", hosts.entries[0].Host)
		}
	})
}

func TestParseEntry(t *testing.T) {
	tests := []struct {
		input    string
		expected Entry
	}{
		{
			"109.94.209.70   fitgirl-repack.net      # Fake FitGirl site",
			Entry{
				IP:      net.IPv4(109, 94, 209, 70),
				Host:    "fitgirl-repack.net",
				Comment: "# Fake FitGirl site",
			},
		},
		{
			"127.0.0.1       kubernetes.docker.internal      #",
			Entry{
				IP:      net.IPv4(127, 0, 0, 1),
				Host:    "kubernetes.docker.internal",
				Comment: "#",
			},
		},
		{
			"127.0.0.1 k#",
			Entry{
				IP:   net.IPv4(127, 0, 0, 1),
				Host: "k#",
			},
		},
		{
			"127.0.0.1 #k#",
			Entry{
				IP:   net.IPv4(127, 0, 0, 1),
				Host: "#k#",
			},
		},
		{
			"127.0.0.1 mino",
			Entry{
				IP:   net.IPv4(127, 0, 0, 1),
				Host: "mino",
			},
		},
		{
			"188.245.227.222 meep-mino",
			Entry{
				IP:   net.IPv4(188, 245, 227, 222),
				Host: "meep-mino",
			},
		},
		{
			"127.0.0.1 localhost.com",
			Entry{
				IP:   net.IPv4(127, 0, 0, 1),
				Host: "localhost.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			entry, err := parseEntry([]byte(tt.input))
			require.NoError(t, err)
			assert.Equal(t, entry.Host, tt.expected.Host)
			assert.Equal(t, entry.IP, tt.expected.IP)
		})
	}
}
