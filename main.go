package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	// ErrUnknownArgs is an error specifying a wrong argument input
	ErrUnknownArgs = errors.New("unknown arguments")
	// ErrTagNotFound is an error specifying a missing tag
	ErrTagNotFound = errors.New("not found")
)

func usage() {
	u := `
"` + name + `" extracts the text between tags.

A tag is defined by "` + START + `" and "` + END + `"

Example for a file containing:

	// [START mytag]
	var Bla = "bla"
	// [END mytag]

executing: 

	> ` + name + ` mytag myfile.go 

will yield: 

	var Bla = "bla"
`
	fmt.Printf("%s\n", u)
}

const (
	// START is the tag to start a section to extract
	START = "// [START "
	// END is the tag to end a section to extract
	END = "// [END "

	name      = "inbtw"
	end_token = "]"
)

func main() {
	if err := run(os.Args[1:]...); err != nil {
		if errors.Is(err, ErrUnknownArgs) {
			usage()
			os.Exit(2)
		}
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run(args ...string) error {
	switch l := len(args); {
	case l == 1:
		tag := args[0]

		m := ExtractTags(os.Stdin)
		x, ok := m[tag]
		if !ok {
			return fmt.Errorf("tag %q: %w", tag, ErrTagNotFound)
		}
		_, _ = fmt.Fprint(os.Stdout, x)

	case l == 2:
		tag := args[0]
		fp := args[1]
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

		m := ExtractTags(f)
		x, ok := m[tag]
		if !ok {
			return fmt.Errorf("tag %q: %w", tag, ErrTagNotFound)
		}
		_, err = fmt.Fprint(os.Stdout, x)
		if err != nil {
			return fmt.Errorf("write %s: %w", fp, err)
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnknownArgs, args)
	}

	return nil
}

// ExtractTags parses the text between [START] and [END]
func ExtractTags(r io.Reader) map[string]string {
	s := bufio.NewScanner(r)
	m := make(map[string]string)
	ent := make(map[string]struct{})

	for s.Scan() {
		line := s.Text()

		if name, ok := extractTagName(START, line); ok {
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
	return m
}

func extractTagName(fix, line string) (string, bool) {
	_, after, found := strings.Cut(line, fix)
	if !found {
		return "", false
	}
	name, _, ok := strings.Cut(after, end_token)
	if !ok {
		return "", false
	}
	return strings.Trim(name, " \t"), true
}
