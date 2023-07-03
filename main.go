package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

//go:generate mdox fmt -l README.md

var (
	// ErrUnknownArgs is an error specifying a wrong argument input
	ErrUnknownArgs = errors.New("unknown arguments")
	// ErrNotFound is an error specifying a missing tag
	ErrNotFound = errors.New("not found")
	// ErrDuplicate is an error specifying a duplicate start tag
	ErrDuplicate = errors.New("duplicate")
)

var Usage = func() {
	u := name + `" extracts the text between tags.

	A tag is defined by "` + START + `" and "` + END + `"

	Example for a file containing:

		// [START mytag]
		var Bla = "bla"
		// [END mytag]

	executing:

		> ` + name + ` -tag mytag -f myfile.go

	will yield:

		var Bla = "bla"
	`
	fmt.Printf("%s\n", u)
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

const (
	// START is the tag to start a section to extract
	START = "// [START "
	// END is the tag to end a section to extract
	END = "// [END "

	name     = "inbtw"
	endToken = "]"
)

func main() {
	os.Exit(main1())
}

func main1() int {
	var tag, in string
	flag.StringVar(&tag, "tag", "", "tag containing the text to extract.")
	flag.StringVar(
		&in,
		"f",
		"",
		"file(s) to parse, multiple files can be separated by ',', '-' for stdin.",
	)
	flag.Usage = Usage
	flag.Parse()

	if err := run(tag, in); err != nil {
		if errors.Is(err, ErrUnknownArgs) {
			flag.Usage()
			return 2
		}
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return 1
	}
	return 0
}

func run(tag string, in string) error {
	for _, s := range strings.Split(in, ",") {
		if err := Extract(s, tag, os.Stdout); err != nil {
			return err
		}
	}

	return nil

}

func Extract(fp string, tag string, w io.Writer) error {
	var in io.Reader
	if fp == "" {
		return ErrUnknownArgs
	}
	if fp != "-" {

		f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return fmt.Errorf("opening file %s: %w", fp, err)
		}
		defer func(f *os.File) {
			ferr := f.Close()
			if ferr != nil {
				_, _ = fmt.Fprint(os.Stderr, ferr)
			}
		}(f)
		in = f
	} else {
		in = os.Stdin
	}

	m, err := ExtractTags(in)
	if err != nil {
		return err
	}
	if tag == "" {
		for k, v := range m {
			_, _ = fmt.Fprintf(w, "-- %s --\n%s\n", k, v)
		}
		return nil
	}

	x, ok := m[tag]
	if !ok {
		return fmt.Errorf("tag %q: %w", tag, ErrNotFound)
	}
	_, err = fmt.Fprint(w, x)
	if err != nil {
		return fmt.Errorf("write %s: %w", fp, err)
	}

	// return fmt.Errorf("%w: %s", ErrUnknownArgs, args)

	return nil
}

// ExtractTags parses the text between [START] and [END]
func ExtractTags(r io.Reader) (map[string]string, error) {
	s := bufio.NewScanner(r)
	m := make(map[string]string)
	ent := make(map[string]struct{})
	var i uint64 = 0

	for s.Scan() {
		line := s.Text()
		i++

		if name, ok := extractTagName(START, line); ok {
			if _, ok := ent[name]; ok {
				return nil, fmt.Errorf(
					"tag %s (line %d): %w",
					name, i, ErrDuplicate,
				)
			}
			ent[name] = struct{}{}
			continue
		}
		if name, ok := extractTagName(END, line); ok {
			if _, ok := ent[name]; ok {
				delete(ent, name)
				continue
			}
		}

		for e := range ent {
			if m[e] == "" {
				m[e] = line
				continue
			}
			m[e] = m[e] + "\n" + line
		}
	}
	return m, nil
}

func extractTagName(fix, line string) (string, bool) {
	_, after, found := strings.Cut(line, fix)
	if !found {
		return "", false
	}
	name, _, ok := strings.Cut(after, endToken)
	if !ok {
		return "", false
	}
	return strings.Trim(name, " \t"), true
}
