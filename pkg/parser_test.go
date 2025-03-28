package pkg

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
