package pkg

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterOptions(t *testing.T) {
	tests := []struct {
		entry       Entry
		filter      filterOptions
		shouldMatch bool
		name        string
	}{
		{
			name:  "happy path: IP",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: "localhost"},
			filter: newFilterOptions(
				WithIPs(net.IPv4(127, 0, 0, 1)),
			),
			shouldMatch: true,
		},
		{
			name:  "happy path: Host",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: "localhost"},
			filter: newFilterOptions(
				WithHosts("localhost"),
			),
			shouldMatch: true,
		},
		{
			name:  "happy path: Comment",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: "localhost", Comment: "Izu"},
			filter: newFilterOptions(
				WithComments("Izu"),
			),
			shouldMatch: true,
		},
		{
			name:  "happy path: No comment",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: "localhost", Comment: ""},
			filter: newFilterOptions(
				WithNoComment(),
			),
			shouldMatch: true,
		},
		{
			name:  "",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: "localhost", Comment: ""},
			filter: newFilterOptions(
				WithNoComment(),
				WithHosts(),
				WithIPs(),
				WithNoComment(),
			),
			shouldMatch: true,
		},
		{
			name:  "edge case: nil IP in entry",
			entry: Entry{IP: nil, Host: "localhost", Comment: "test"},
			filter: newFilterOptions(
				WithIPs(net.IPv4(127, 0, 0, 1)),
			),
			shouldMatch: false,
		},
		{
			name:  "edge case: empty host string in filter",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: ""},
			filter: newFilterOptions(
				WithHosts(""),
			),
			shouldMatch: true, // or false depending on intended behavior
		},
		{
			name:  "edge case: IP in entry is IPv6, filter expects IPv4",
			entry: Entry{IP: net.ParseIP("::1"), Host: "localhost"},
			filter: newFilterOptions(
				WithIPs(net.IPv4(127, 0, 0, 1)),
			),
			shouldMatch: false,
		},
		{
			name:  "edge case: IP in entry matches IPv6 ::1",
			entry: Entry{IP: net.ParseIP("::1"), Host: "localhost"},
			filter: newFilterOptions(
				WithIPs(net.ParseIP("::1")),
			),
			shouldMatch: true,
		},
		{
			name:  "edge case: filter has multiple IPs, one matches",
			entry: Entry{IP: net.IPv4(10, 0, 0, 1), Host: "example.com"},
			filter: newFilterOptions(
				WithIPs(net.IPv4(127, 0, 0, 1), net.IPv4(10, 0, 0, 1)),
			),
			shouldMatch: true,
		},
		{
			name:  "edge case: comment is whitespace only",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: "localhost", Comment: "  "},
			filter: newFilterOptions(
				WithNoComment(),
			),
			shouldMatch: true,
		},
		{
			name:        "edge case: empty filters (match all)",
			entry:       Entry{IP: net.IPv4(8, 8, 8, 8), Host: "dns.google", Comment: "Google DNS"},
			filter:      newFilterOptions(),
			shouldMatch: true,
		},
		{
			name:  "extreme case: invalid IP (empty slice)",
			entry: Entry{IP: net.IP{}, Host: "localhost", Comment: "invalid"},
			filter: newFilterOptions(
				WithIPs(net.IPv4(127, 0, 0, 1)),
			),
			shouldMatch: false,
		},
		{
			name:  "extreme case: invalid IP string parsing",
			entry: Entry{IP: net.ParseIP("999.999.999.999"), Host: "localhost"},
			filter: newFilterOptions(
				WithIPs(net.IPv4(127, 0, 0, 1)),
			),
			shouldMatch: false,
		},
		{
			name:  "extreme case: long hostname",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: strings.Repeat("a", 255)},
			filter: newFilterOptions(
				WithHosts(strings.Repeat("a", 255)),
			),
			shouldMatch: true,
		},
		{
			name:  "extreme case: overly long comment (1MB)",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: "localhost", Comment: strings.Repeat("x", 1024*1024)},
			filter: newFilterOptions(
				WithComments(strings.Repeat("x", 1024*1024)),
			),
			shouldMatch: true,
		},
		{
			name:  "extreme case: comment contains newline and special characters",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: "localhost", Comment: "Line1\nLine2\t\u2603"},
			filter: newFilterOptions(
				WithComments("Line1\nLine2\t\u2603"),
			),
			shouldMatch: true,
		},
		{
			name:  "extreme case: filter has no values and expects no comment",
			entry: Entry{IP: net.IPv4(1, 2, 3, 4), Host: "example", Comment: ""},
			filter: newFilterOptions(
				WithNoComment(),
				WithIPs(),
				WithHosts(),
				WithComments(),
			),
			shouldMatch: true,
		},
		{
			name:  "extreme case: all fields empty",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: "localhost", Comment: "Izu"},
			filter: newFilterOptions(
				WithIPs(nil),
				WithHosts(""),
				WithComments(""),
			),
			shouldMatch: false,
		},
		{
			name:  "extreme case: mismatching Unicode comment (NFC vs NFD)",
			entry: Entry{IP: net.IPv4(127, 0, 0, 1), Host: "localhost", Comment: "é"},
			filter: newFilterOptions(
				WithComments("e\u0301"), // decomposed é
			),
			shouldMatch: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.shouldMatch, test.filter.Match(test.entry))
		})
	}
}
