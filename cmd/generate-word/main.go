// generate-word is a simple command line tool that generates an Aslan word.
package main

import (
	"context"
	"fmt"
	"log" // Added log for fatal errors
	"os"  // Added os import for os.Exit

	"github.com/carloscasalar/aslan-words/pkg/aslanwords"
)

var osExit = os.Exit // Allow os.Exit to be mocked for testing

func main() {
	cliOpts := readOptionsOrFail() // Renamed to cliOpts to match opts.go

	ctx := context.Background()
	generatorOpts, err := cliOpts.ToAslanGeneratorOptions()
	if err != nil {
		log.Printf("Error processing command-line options: %v", err)
		osExit(1) // Use mocked exit
	}

	word, err := aslanwords.Generate(ctx, generatorOpts...)
	if err != nil {
		log.Printf("Unable to generate an Aslan word: %v", err)
		osExit(1) // Use mocked exit
	}

	fmt.Println(word)
}
