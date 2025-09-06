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

// generateNameSyllablesWereSetByUser checks if the user's options would alter the default syllable count.
func generateNameSyllablesWereSetByUser(passedOpts ...GeneratorOption) bool {
	if len(passedOpts) == 0 {
		return false
	}

	defaultSyllableConfig := newGeneratorOptions().numberOfSyllablesOpts
	
	optsApplied := newGeneratorOptions()
	for _, o := range passedOpts {
		o(optsApplied)
	}
	effectiveSyllableConfig := optsApplied.numberOfSyllablesOpts

	// Compare by type and value.
	// Default is randomAmountOpt{from: 2, to: 6}
	defaultRandom, okDefault := defaultSyllableConfig.(randomAmountOpt)
	effectiveRandom, okEffectiveRandom := effectiveSyllableConfig.(randomAmountOpt)
	_, okEffectiveFixed := effectiveSyllableConfig.(fixedAmountOpt)

	if okEffectiveFixed {
		return true // User explicitly set a fixed number.
	}

	if okDefault && okEffectiveRandom {
		// Both are random; if values differ, user set them.
		if defaultRandom.from == effectiveRandom.from && defaultRandom.to == effectiveRandom.to {
			return false // Same as default random range.
		}
		return true // Different random range.
	}
	
	// If types are different and not handled above (e.g. default is random, effective is something else entirely)
	// This would imply a custom amountOptions was introduced, or it's fixed (already handled).
	// If effective is still random but default wasn't (hypothetically), it's also a change.
	// For this function, if it's not fixed and not a different random, assume it wasn't explicitly set to override default.
	if !okDefault && okEffectiveRandom { // Should not happen if newGeneratorOptions is consistent
		return true
	}

	return false // If it's still the default random config, or len(opts) == 0
}


// GenerateName generates a random Aslan name based on gender and options.
func GenerateName(ctx context.Context, gender string, opts ...GeneratorOption) (string, error) {
	options := newGeneratorOptions() // Base options with library defaults (e.g., 2-6 syllables)

	// Determine if user options include specific syllable settings
	syllablesSetByPassedOpts := generateNameSyllablesWereSetByUser(opts...)

	// Apply all passed options to the main 'options' object
	for _, o := range opts {
		o(options)
	}

	if gender != "male" && gender != "female" {
		return "", fmt.Errorf("invalid gender: %s, must be 'male' or 'female'", gender)
	}

	// If user options did not specify syllable count, apply gender-based defaults.
	if !syllablesSetByPassedOpts {
		if gender == "male" {
			WithNumberOfSyllablesBetween(3, 4)(options) // Modifies 'options' in place
		} else if gender == "female" { // "female"
			WithNumberOfSyllablesBetween(2, 3)(options) // Modifies 'options' in place
		}
	}

	err := options.Validate()
	if err != nil {
		return "", fmt.Errorf("invalid options: %w", err)
	}

	wordTemplate := syllable.GenerateTemplate(options.numberOfSyllables())
	gen, err := fantasyname.Compile(wordTemplate.String(), fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		return "", fmt.Errorf("unexpected error generating the aslan name: %w", err)
	}
	return gen.String(), nil
}

// All subsequent function definitions were identified as duplicates
// from previous patching attempts and have been removed.
// The correct versions of generateNameSyllablesWereSetByUser (lines 43-78)
// and GenerateName (lines 81-125) are preserved above.
