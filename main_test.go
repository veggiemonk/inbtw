package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractTags(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    map[string]string
		wantErr error
	}{
		{
			name: "simple",
			in: `
#// [START xxx]
echo "hello"
#// [END xxx]
`,
			want: map[string]string{"xxx": `echo "hello"`},
		},
		{
			name: "double",
			in: `
// [START xxx]
something here
// [END xxx]
something something
something
// [START xxx]
again something here
// [END xxx]

`,
			want: map[string]string{
				"xxx": "something here\nagain something here",
			},
		},
		{
			name: "err: no start tag",
			in: `
something here
// [END xxx]
`,
			want: map[string]string{},
		},
		{
			name: "no end tag",
			in: `
hello
// [START xxx]
something
here
no end tag
`,
			want: map[string]string{"xxx": "something\nhere\nno end tag"},
		},
		{
			name: "multi",
			in: `
// [START xxx]
something here
// [END xxx]
something something
something
// [START aaa]
again something here
// [END aaa]
`,
			want: map[string]string{
				"aaa": "again something here",
				"xxx": "something here",
			},
		},
		{
			name: "interleave",
			in: `
// [START xxx]
something here
something something
// [START aaa]
something
// [END xxx]
again something here
// [END aaa]
`,
			want: map[string]string{
				"aaa": "something\nagain something here",
				"xxx": "something here\nsomething something\nsomething",
			},
		},
		{
			name: "err: interleave same",
			in: `
// [START xxx]
something here
something something
// [START xxx]
something
// [END xxx]
again something here
// [END xxx]
`,
			want:    nil,
			wantErr: ErrDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := ExtractTags(bytes.NewReader([]byte(tt.in)))
				if tt.wantErr != nil {
					if !errors.Is(err, tt.wantErr) {
						t.Fatal("should err with", tt.wantErr)
					}
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("\n- got\n+ want:\n %s", diff)
				}
			},
		)
	}
}

func TestExtractTagName(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		fix     string
		want    string
		wantErr bool
	}{
		{
			name: "start: bla",
			in:   "// [START bla]",
			fix:  START,
			want: "bla",
		},
		{
			name: "start: tab",
			in:   "\t// [START bla]",
			fix:  START,
			want: "bla",
		},
		{
			name: "start: tab and inside space",
			in:   "\t// [START bla   ]",
			fix:  START,
			want: "bla",
		},
		{
			name:    "start: end prefix",
			in:      "\t // [END bla   ]",
			fix:     START,
			want:    "",
			wantErr: true,
		},
		{
			name: "end: bla",
			in:   "\t // [END bla   ]",
			fix:  END,
			want: "bla",
		},
		{
			name:    "end: tab and inside space",
			in:      "\t// [START bla   ]",
			fix:     END,
			wantErr: true,
		},
		{
			name: "end: prefix with tabs",
			in:   "\t   \t// [END bla bla]",
			fix:  END,
			want: "bla bla",
		},
		{
			name:    "begin ignored",
			in:      "\t\t// [BEGIN sup in]",
			fix:     START,
			wantErr: true,
		},
		{
			name: "edge after",
			in:   "\t\t// [START sup!n]!!",
			fix:  START,
			want: "sup!n",
		},
		{
			name: "edge tab after name",
			in:   "\t\t// [START sup!n\t]!!",
			fix:  START,
			want: "sup!n",
		},
		{
			name: "edge tab before name",
			in:   "\t\t// [START \tsup!n]!!",
			fix:  START,
			want: "sup!n",
		},
		{
			name: "end: # beginning of line ",
			in:   "#// [END xxx]",
			fix:  END,
			want: "xxx",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, ok := extractTagName(tt.fix, tt.in)
				if ok && tt.wantErr {
					t.Fatal("should err:", "`"+tt.in+"`")
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("\n- got\n+ want:\n %s", diff)
				}
			},
		)
	}
}
