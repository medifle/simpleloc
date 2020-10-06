package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// parse CLI argument
	if len(os.Args) <= 1 {
		fmt.Println("usage: simpleloc <filepath>")
		return
	}
	filePath := os.Args[1]

	// TODO: get language specific config
	extension := filepath.Ext(filePath)
	fmt.Println(extension)

	// read the entire file content into memory
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File contents:\n%s\n", content)
}
