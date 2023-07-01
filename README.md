# inbtw

`inbtw` (a.k.a. "in between"), is a small utility to extract the text in between two tags.

```
"inbtw" extracts the text between tags.

A tag is defined by "// [START " and "// [END "

Example for a file containing:

        // [START mytag]
        var Bla = "bla"
        // [END mytag]

executing: 

        > inbtw mytag myfile.go 

will yield: 

        var Bla = "bla"
```

## Purpose

To extract code snippets in order to render them in a document. 
It requires less maintenance than something like `sed` because it needs line numbers,
and less complicated than regex or `awk`. 

## Installing

```shell
go install github.com/veggiemonk/inbtw@latest 
```
