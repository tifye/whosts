package pkg

import (
	"io"
	"strings"
)

type Hosts struct {
	entries []Entry
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
