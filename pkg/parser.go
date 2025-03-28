package pkg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	ErrInvalidIP    = errors.New("invalid IP address")
	ErrInvalidEntry = errors.New("invalid host entry")
)

type Entry struct {
	IP      net.IP
	Host    string
	Comment string
}

func ParseEntries(r io.Reader) ([]Entry, error) {
	entries := make([]Entry, 0)
	buf := bufio.NewReader(r)
	ln := -1
	run := true
	for run {
		ln += 1
		b, err := buf.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				run = false
			} else {
				return nil, lineParseErr(ln, fmt.Errorf("read bytes: %s", err))
			}
		}

		if bytes.HasPrefix(b, []byte{'#'}) {
			continue
		}

		entry, err := parseEntry(b)
		if err != nil {
			return nil, lineParseErr(ln, err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func parseEntry(b []byte) (Entry, error) {
	b = bytes.TrimSpace(b)
	parts := bytes.Fields(b)
	if len(parts) < 2 {
		return Entry{}, invalidEntryErr(b)
	}

	ip, host := parts[0], parts[1]
	if len(host) == 0 { // Should not happen after len(parts) < 2
		return Entry{}, fmt.Errorf("empty host")
	}

	var commentStr string
	if len(parts) > 2 {
		commentParts := parts[2:]
		comment := bytes.Join(commentParts, []byte{' '})
		if comment[0] != '#' {
			invalidEntryErr(b)
		}
		commentStr = string(comment)
	}

	var m net.IP
	err := m.UnmarshalText(ip)
	if err != nil {
		return Entry{}, invalidIPErr(ip)
	}

	return Entry{
		IP:      m,
		Host:    string(host),
		Comment: commentStr,
	}, nil
}

func invalidIPErr(d []byte) error {
	return fmt.Errorf("%w:%s", ErrInvalidIP, d)
}

func invalidEntryErr(d []byte) error {
	return fmt.Errorf("%w:%s", ErrInvalidEntry, d)
}

func lineParseErr(line int, err error) error {
	return fmt.Errorf("err on line %d, %w", line, err)
}
