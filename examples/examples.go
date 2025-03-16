// Description: Example of how to use the aslanwords package to generate Aslan words.
package main

import (
	"context"
	"fmt"

	"github.com/carloscasalar/aslan-words/pkg/aslanwords"
)

func main() {
	ctx := context.Background()
	word, err := aslanwords.Generate(ctx, aslanwords.WithNumberOfSyllables(3))
	if err != nil {
		fmt.Println("Error generating Aslan word:", err)
		return
	}
	fmt.Printf("Aslan word with 3 syllables: %s\n", word)

	word = aslanwords.MustGenerate(ctx, aslanwords.WithNumberOfSyllablesBetween(3, 6))
	fmt.Printf("Aslan word with between 3 and 6 syllables: %s\n", word)
}
