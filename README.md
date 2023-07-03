# inbtw

`inbtw` (a.k.a. "in between"), is a small utility to extract the text in between two tags.

```bash mdox-exec="inbtw" mdox-expect-exit-code=2
inbtw" extracts the text between tags.

	A tag is defined by "// [START " and "// [END "

	Example for a file containing:

		// [START mytag]
		var Bla = "bla"
		// [END mytag]

	executing:

		> inbtw-tag mytag -f myfile.go

	will yield:

		var Bla = "bla"
	
Usage of inbtw:
  -f string
    	file(s) to parse, multiple files can be separated by ',', '-' for stdin.
  -tag string
    	tag containing the text to extract.
  -trim int
    	trim left number of spaces.
```

## Purpose

To extract code snippets in order to render them in a document. It requires less maintenance than something like `sed` because it needs line numbers, and less complicated than regex or `awk`.

## Installing

```shell
go install github.com/veggiemonk/inbtw@latest 
```
