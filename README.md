## Simpleloc
A simple tool to count lines of code, comments (single and multiline), TODOs.

## How to use
Download zip and enter the directory using terminal

compile
```
go build main.go
```
run
```
./main <file path>
```

## Assumption
- Line of code with a trailing multiline comment is considered
- Line of code with a leading multiline comment is ignored
- Line of code contains blank line
- Python does not have multiline comment (docstring is not considered as comment as it is a string)
- Language support can be added in languages.json and run `go generate` in repo's root directory before compile
- TODOs match is case-sensitive and only valid in comments
- Supported language is in language.json which means most of other languages can be easily supported
- File name starting with '.' is not ignored