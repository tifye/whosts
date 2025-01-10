package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LexerItemTypes(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected []itemType
	}{
		{
			name:     "empty",
			input:    "",
			expected: []itemType{itemEOF},
		},
		{
			name:  "only comments",
			input: "#meep\n#mino",
			expected: []itemType{
				itemComment,
				itemNewline,
				itemComment,
				itemEOF,
			},
		},
		{
			name:  "only newlines",
			input: "\n\n\n\n",
			expected: []itemType{
				itemNewline,
				itemNewline,
				itemNewline,
				itemNewline,
				itemEOF,
			},
		},
		{
			name:  "comment + space + newline",
			input: "#mino  \n",
			expected: []itemType{
				itemComment,
				itemNewline,
				itemEOF,
			},
		},
		{
			name: "default",
			input: ` # End of section # foo
					109.94.209.70	*.fitgirl-repacks.xyz	# Fake FitGirl site
					192.168.18.175	host.docker.internal`,
			expected: []itemType{
				itemComment,
				itemNewline,
				itemIP,
				itemHost,
				itemComment,
				itemNewline,
				itemIP,
				itemHost,
				itemEOF,
			},
		},
		{
			name:  "entry, no host",
			input: `192.168.18.175 `,
			expected: []itemType{
				itemIP,
				itemError,
			},
		},
		{
			name:  "entry, no comment",
			input: `0.0.0.0 host`,
			expected: []itemType{
				itemIP,
				itemHost,
				itemEOF,
			},
		},
		{
			name:  "entry, with comment",
			input: `0.0.0.0 host #comment#`,
			expected: []itemType{
				itemIP,
				itemHost,
				itemComment,
				itemEOF,
			},
		},
		{
			name:  "entry, with invalid host + comment",
			input: `0.0.0.0 host# comment#`,
			expected: []itemType{
				itemIP,
				itemError,
			},
		},
		{
			name:  "entry, no comment, \\w new line",
			input: "0.0.0.0 host\n",
			expected: []itemType{
				itemIP,
				itemHost,
				itemNewline,
				itemEOF,
			},
		},
		{
			name:  "two entries, new line sperated",
			input: "0.0.0.0 host\n0.0.0.0 host",
			expected: []itemType{
				itemIP,
				itemHost,
				itemNewline,
				itemIP,
				itemHost,
				itemEOF,
			},
		},
		{
			name:     "single digit host",
			input:    "0.0.0.0 1",
			expected: []itemType{itemIP, itemError},
		},
		{
			name:     "single asterisk host",
			input:    "0.0.0.0 *",
			expected: []itemType{itemIP, itemError},
		},
		{
			name:     "asterisk followed by non '.'",
			input:    "0.0.0.0 *d",
			expected: []itemType{itemIP, itemError},
		},
		{
			name:     "asterisk following character",
			input:    "0.0.0.0 d.*",
			expected: []itemType{itemIP, itemError},
		},
		{
			name:     "larger than 255 IP",
			input:    "999.0.0.0 host",
			expected: []itemType{itemError},
		},
		{
			name:     "IP leading zero number",
			input:    "099.0.0.0 host",
			expected: []itemType{itemError},
		},
		{
			name:     "IP 256",
			input:    "256.0.0.0 host",
			expected: []itemType{itemError},
		},
		{
			name:     "IP 192.168.18.180",
			input:    "192.168.18.180 host",
			expected: []itemType{itemIP, itemHost, itemEOF},
		},
	}

	for _, td := range tt {
		t.Run(td.name, func(t *testing.T) {
			_, itemch := lex(td.input)

			items := make([]item, 0, len(td.expected))
			for i := range itemch {
				items = append(items, i)
			}

			got := make([]itemType, 0, len(td.expected))
			for _, i := range items {
				got = append(got, i.typ)
			}

			if !assert.Equal(t, td.expected, got) {
				t.Logf("last item: %q", items[len(items)-1])
			}
		})
	}
}

func Test_Lexer(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "entry",
			input:    "0.0.0.0 host",
			expected: []string{"0.0.0.0", "host", ""},
		},
		{
			name:  "entry with comment",
			input: "0.0.0.0 host #comment",
			expected: []string{
				"0.0.0.0",
				"host",
				"#comment",
				"",
			},
		},
		{
			name:  "entry with comment and newline",
			input: "0.0.0.0 host #comment\n",
			expected: []string{"0.0.0.0",
				"host",
				"#comment",
				"\n",
				"",
			},
		},
		{
			name:  "entry with comment and newline",
			input: "0.0.0.0 host\n0.0.0.0 host",
			expected: []string{"0.0.0.0",
				"host",
				"\n",
				"0.0.0.0",
				"host",
				"",
			},
		},
		{
			name:     "single letter host",
			input:    "0.0.0.0 a",
			expected: []string{"0.0.0.0", "a", ""},
		},
	}

	for _, td := range tt {
		t.Run(td.name, func(t *testing.T) {
			_, itemch := lex(td.input)

			items := make([]item, 0, len(td.expected))
			for i := range itemch {
				items = append(items, i)
			}

			got := make([]string, 0, len(td.expected))
			for _, i := range items {
				got = append(got, i.val)
			}

			if !assert.Equal(t, td.expected, got) {
				t.Logf("last item: %q", items[len(items)-1])
			}
		})
	}
}
