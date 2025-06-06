// Package aslanwords provides functions to generate-word random Aslan words.
package aslanwords

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/carloscasalar/aslan-words/internal/syllable"
	"github.com/s0rg/fantasyname"
)

// Generate generates a random Aslan word with the given options.
// If no options are provided, it will generate-word a word with a random number of syllables between 2 and 6.
func Generate(ctx context.Context, opts ...GeneratorOption) (string, error) {
	options := newGeneratorOptions()
	for _, o := range opts {
		o(options)
	}
	err := options.Validate()
	if err != nil {
		return "", fmt.Errorf("invalid options: %w", err)
	}

	wordTemplate := syllable.GenerateTemplate(options.numberOfSyllables())
	gen, err := fantasyname.Compile(wordTemplate.String(), fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		return "", fmt.Errorf("unexpected error generating the aslan word: %w", err)
	}
	return gen.String(), nil
}

// MustGenerate generates a random Aslan word with the given options.
// If no options are provided, it will generate-word a word with a random number of syllables between 2 and 6.
// If an error occurs, it will panic.
func MustGenerate(ctx context.Context, opts ...GeneratorOption) string {
	word, err := Generate(ctx, opts...)
	if err != nil {
		panic(err)
	}
	return word
}
