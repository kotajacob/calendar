package keyword

import (
	"reflect"
	"strings"
	"testing"
)

func TestMatch(t *testing.T) {
	type test struct {
		description string
		keywords    Keywords
		file        string
		want        bool
		wantKeyword Keyword
	}

	tests := []test{
		{
			description: "match ATTN on first line",
			keywords: Keywords{
				Keyword{Keyword: "ATTN", Color: "2"},
			},
			file:        "ATTN\n\nDo something.",
			want:        true,
			wantKeyword: Keyword{Keyword: "ATTN", Color: "2"},
		},
		{
			description: "match ATTN on last line",
			keywords: Keywords{
				Keyword{Keyword: "ATTN", Color: "2"},
			},
			file:        "Do something.\nATTN",
			want:        true,
			wantKeyword: Keyword{Keyword: "ATTN", Color: "2"},
		},
		{
			description: "dont match ATTN across lines",
			keywords: Keywords{
				Keyword{Keyword: "ATTN", Color: "2"},
			},
			file:        "AT\nTN",
			want:        false,
			wantKeyword: Keyword{},
		},
	}

	for _, tc := range tests {
		r := strings.NewReader(tc.file)
		gotKeyword, gotOK := tc.keywords.Match(r)
		if gotOK != tc.want {
			t.Fatalf(
				"got: %v, want: %v, for: %v\n",
				gotOK,
				tc.want,
				tc.description,
			)
		}
		if !reflect.DeepEqual(gotKeyword, tc.wantKeyword) {
			t.Fatalf(
				"got: %v, want: %v, for: %v\n",
				gotKeyword,
				tc.wantKeyword,
				tc.description,
			)
		}
	}
}
