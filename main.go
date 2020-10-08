package main

import (
	"github.com/medifle/simpleloc/processor"
)

//pre-generate language config in Go format
//go:generate go run scripts/include.go
func main() {
	processor.Process()
}
