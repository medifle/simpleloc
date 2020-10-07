package processor

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type stateType uint32

const (
	SBlank stateType = iota
	SCode
	SSingleLineComment
	SMultiLineComment
)

const (
	Space byte = 32
	Lf         = 10 // new line "\n"
)

// stateType machine for counting lines
type machine struct {
	state stateType
}

// method for machine to advance to the next stateType
func (p *machine) step(input byte) {
	switch p.state {
	case SBlank:
		// TODO:

	}
}

func countStats(file *File) {
	// If file is empty, no line is counted
	if len(file.Content) == 0 {
		return
	}
	//machine := &machine{state: SBlank}
	//for i := 0; i < len(file.Content); i++ {
	//	machine.step(file.Content[i])
	//}
}

func printStats(file *File) {
	fmt.Printf("Total # of lines: %d\n", file.Lines)
	fmt.Printf("Total # of comment lines: %d\n", file.SingleLineComment+file.MultiLineComment)
	fmt.Printf("Total # of single line comments: %d\n", file.SingleLineComment)
	fmt.Printf("Total # of comment lines within block comments: %d\n", file.MultiLineComment)
	fmt.Printf("Total # of block line comments: %d\n", file.BlockComment)
	fmt.Printf("Total # of TODO’s: %d\n", file.Todo)
	fmt.Println()
}

func Process() {
	// Parse CLI argument
	if len(os.Args) <= 1 {
		fmt.Println("usage: simpleloc <filepath>")
		return
	}
	filePath := os.Args[1]

	// TODO: get language specific config
	extension := filepath.Ext(filePath)

	// Read the entire file content into memory
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	file := &File{
		Extension: extension,
		Content:   content,
	}

	// Count lines
	countStats(file)

	// Output the result
	printStats(file)
}
