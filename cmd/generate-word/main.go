// generate-word is a simple command line tool that generates an Aslan word.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/carloscasalar/aslan-words/pkg/aslanwords"
)

func main() {
	opts := readOptionsOrFail()

	ctx := context.Background()
	word, err := aslanwords.Generate(ctx, aslanwords.WithNumberOfSyllables(opts.NumberOfSyllables))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(word)
}
