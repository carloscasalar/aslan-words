// generate-word is a simple command line tool that generates an Aslan word.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/carloscasalar/aslan-words/pkg/aslanwords"
)

func main() {
	parsedCliOpts := readOptionsOrFail() // Assume this function will be updated in the next step to return an object with a Gender field

	ctx := context.Background()
	var word string
	var err error

	genOpts := []aslanwords.GeneratorOption{}
	// Assuming NumberOfSyllables will be 0 if not set by the user, or a positive integer if set.
	// The readOptionsOrFail function would be responsible for this parsing logic.
	if parsedCliOpts.NumberOfSyllables > 0 {
		genOpts = append(genOpts, aslanwords.WithNumberOfSyllables(parsedCliOpts.NumberOfSyllables))
	}

	if parsedCliOpts.Gender != "" {
		word, err = aslanwords.GenerateName(ctx, parsedCliOpts.Gender, genOpts...)
	} else {
		word, err = aslanwords.Generate(ctx, genOpts...)
	}

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(word)
}
