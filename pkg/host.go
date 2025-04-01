package pkg

import (
	"io"
	"net"
	"strings"
)

type Hosts struct {
	entries []Entry
}

func NewHosts(entries []Entry) Hosts {
	return Hosts{entries: entries}
}

func (h *Hosts) AddEntry(entry Entry) {
	h.entries = append(h.entries, entry)
}

func (h Hosts) Entries() []Entry {
	return h.entries
}

func (h Hosts) String() string {
	builder := strings.Builder{}
	for _, entry := range h.entries {
		_, _ = builder.WriteString(entry.String() + "\n")
	}
	return builder.String()
}

func (h Hosts) WriteTo(w io.Writer) (n int64, err error) {
	_n, err := w.Write([]byte(h.String()))
	return int64(_n), err
}

type filterOptions struct {
	ips       []net.IP
	hosts     []string
	comments  []string
	noComment bool
	matchAll  bool
}

func newFilterOptions(filters ...FilterOption) filterOptions {
	var filterOpts filterOptions
	for _, f := range filters {
		f(&filterOpts)
	}
	return filterOpts
}

func (fo filterOptions) Match(e Entry) bool {
	if fo.matchAll {
		return true
	}

	var ipMatch bool
	for _, ip := range fo.ips {
		if ip.Equal(e.IP) {
			ipMatch = true
		}
	}
	if !ipMatch && len(fo.ips) > 0 {
		return false
	}

	var hostMatch bool
	for _, host := range fo.hosts {
		if host == e.Host {
			hostMatch = true
		}
	}
	if !hostMatch && len(fo.hosts) > 0 {
		return false
	}

	if fo.noComment && strings.TrimSpace(e.Comment) == "" {
		return true
	}

	var commentMatch bool
	for _, comment := range fo.comments {
		if comment == e.Comment {
			commentMatch = true
		}
	}
	if !commentMatch && len(fo.comments) > 0 {
		return false
	}

	return true
}

type FilterOption func(opts *filterOptions)

func WithAll() FilterOption {
	return func(opts *filterOptions) {
		opts.matchAll = true
	}
}

// Filter entries matching any one of the passed host names.
func WithHosts(hosts ...string) FilterOption {
	return func(opts *filterOptions) {
		if opts.hosts == nil {
			opts.hosts = make([]string, 0, len(hosts))
		}
		opts.hosts = append(opts.hosts, hosts...)
	}
}

// Filter entries matching any one of the passed IPs.
func WithIPs(ips ...net.IP) FilterOption {
	return func(opts *filterOptions) {
		if opts.ips == nil {
			opts.ips = make([]net.IP, 0, len(ips))
		}
		opts.ips = append(opts.ips, ips...)
	}
}

// Filter entries containing any one of the passed comments.
func WithComments(comments ...string) FilterOption {
	return func(opts *filterOptions) {
		if opts.comments == nil {
			opts.comments = make([]string, 0, len(comments))
		}
		opts.comments = append(opts.comments, comments...)
	}
}

// Filter entries without comments.
func WithNoComment() FilterOption {
	return func(opts *filterOptions) {
		opts.noComment = true
	}
}

func (h *Hosts) Remove(duplicateOnly bool, filters ...FilterOption) []Entry {
	filterOpts := newFilterOptions(filters...)
	keptEntries := make([]Entry, 0)
	removedEntries := make([]Entry, 0)
	duplicatesCheck := map[string]struct{}{}
	for _, e := range h.entries {
		hasMatch := filterOpts.Match(e)
		if !hasMatch {
			keptEntries = append(keptEntries, e)
			continue
		}

		if !duplicateOnly {
			removedEntries = append(removedEntries, e)
			continue
		}

		eStr := e.String()
		if _, ok := duplicatesCheck[eStr]; !ok {
			duplicatesCheck[eStr] = struct{}{}
			keptEntries = append(keptEntries, e)
		} else {
			removedEntries = append(removedEntries, e)
		}
	}
	h.entries = keptEntries
	return removedEntries
}
