package main

import (
	"github.com/medifle/simpleloc/processor"
)

//go:generate go run scripts/include.go
func main() {
	processor.Process()
}
