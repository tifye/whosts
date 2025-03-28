package pkg

type Hosts struct {
	entries []Entry
}

func (h *Hosts) AddEntry(entry Entry) {
	h.entries = append(h.entries, entry)
}

func (h Hosts) Entries() []Entry {
	return h.entries
}
