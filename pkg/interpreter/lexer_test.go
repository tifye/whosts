package interpreter

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
			name:  "entry, with invalid host+comment",
			input: `0.0.0.0 host# comment#`,
			expected: []itemType{
				itemIP,
				itemHost,
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

func Test_LexComments(t *testing.T) {
	input := `
	# End of section # foo

	# End of section # foo`
	expected := []itemType{
		itemComment,
		itemNewline,
		itemComment,
		itemEOF,
	}

	_, items := lex(input)

	got := make([]itemType, 0, len(expected))
	for i := range items {
		got = append(got, i.typ)
	}

	_ = assert.Equal(t, expected, got)
}
