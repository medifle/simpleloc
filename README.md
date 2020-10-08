## Simpleloc
A simple tool to count lines of code, comments (single and multiline) and TODOs.

## How to use
Download ZIP and enter the directory using a terminal

Compile
```
go build main.go
```
Run
```
./main <file path>
```

## Assumption
- Code with trailing multiline comment is counted as both line of code and line of comment
- Leading multiline comment is ignored
- Line of code contains blank line
- Python does not have multiline comment (docstring is not considered as comment as it is a string)
- TODOs match is case-sensitive and only valid in comments
- Supported languages are in `languages.json` which means most of other languages can be easily supported
- File name starting with `.` is not ignored

## How to add new languages
1. Add new languages in `languages.json` and specify the corresponding `extensions`, `line_comment`, `multi_line` 
and `quotes` sections.

1. Run `go generate` at the repo's root directory before compiling.