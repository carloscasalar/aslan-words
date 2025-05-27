// Package aslanwords provides functions to generate-word random Aslan words.
package aslanwords

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/carloscasalar/aslan-words/internal/syllable"
	"github.com/s0rg/fantasyname"
)

// generateNameFromSyllableTemplate is a helper function to reduce duplication.
// It takes a syllable template definition and generates a name.
func generateNameFromSyllableTemplate(template syllable.TemplateDefinition) (string, error) {
	if template == nil {
		return "", fmt.Errorf("cannot generate name from nil template")
	}
	gen, err := fantasyname.Compile(template.String(), fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		return "", fmt.Errorf("unexpected error compiling template: %w", err)
	}
	return gen.String(), nil
}

// GenerateName generates a random Aslan name for the specified gender.
// Default syllable counts (3-5 for male, 2-4 for female) are applied if no specific syllable options are provided.
func GenerateName(ctx context.Context, gender Gender, userOpts ...GeneratorOption) (string, error) {
	options := newGeneratorOptions()
	options.gender = gender
	// Set wordType for consistency if any downstream logic within template generation might use it,
	// though syllable.GenerateMaleNameTemplate/GenerateFemaleNameTemplate don't directly depend on it.
	if gender == Male {
		options.wordType = TypeMaleName
	} else if gender == Female {
		options.wordType = TypeFemaleName
	}

	for _, o := range userOpts {
		o(options) // This applies user-defined syllable counts and sets isSyllableCountUserSet
	}

	// Apply default syllable counts if not set by user
	if !options.isSyllableCountUserSet {
		if gender == Male {
			WithNumberOfSyllablesBetween(3, 5)(options)
		} else if gender == Female {
			WithNumberOfSyllablesBetween(2, 4)(options)
		}
	}

	if err := options.Validate(); err != nil {
		return "", fmt.Errorf("invalid options for %v name: %w", gender, err)
	}

	var wordTemplate syllable.TemplateDefinition
	numSyllables := options.numberOfSyllables()

	switch gender {
	case Male:
		wordTemplate = syllable.GenerateMaleNameTemplate(numSyllables)
	case Female:
		wordTemplate = syllable.GenerateFemaleNameTemplate(numSyllables)
	default:
		// This case should ideally not be reached if Gender enum is used correctly.
		return "", fmt.Errorf("unknown gender specified for name generation: %v", gender)
	}

	return generateNameFromSyllableTemplate(wordTemplate)
}

// MustGenerateName generates a random Aslan name for the specified gender, panicking on error.
func MustGenerateName(ctx context.Context, gender Gender, userOpts ...GeneratorOption) string {
	name, err := GenerateName(ctx, gender, userOpts...)
	if err != nil {
		panic(fmt.Sprintf("error generating %v name: %v", gender, err))
	}
	return name
}

// Generate generates a random Aslan word, male name, or female name based on the provided options.
// If no specific type option is given, it defaults to generating a generic Aslan word.
func Generate(ctx context.Context, userOpts ...GeneratorOption) (string, error) {
	options := newGeneratorOptions()
	// Apply user options first to determine wordType, gender, and if syllables were user-set
	for _, o := range userOpts {
		o(options)
	}

	switch options.wordType {
	case TypeMaleName:
		// Pass userOpts directly; GenerateName will handle defaults based on isSyllableCountUserSet
		return GenerateName(ctx, Male, userOpts...)
	case TypeFemaleName:
		// Pass userOpts directly; GenerateName will handle defaults based on isSyllableCountUserSet
		return GenerateName(ctx, Female, userOpts...)
	case TypeWord:
		fallthrough
	default: // Catches TypeWord and any undefined wordType
		// For generic words, ensure syllable options are validated and used.
		// Default syllables (2-6) are set by newGeneratorOptions if not in userOpts.
		// If userOpts contain syllable settings, they will be used.
		if err := options.Validate(); err != nil {
			return "", fmt.Errorf("invalid options for word: %w", err)
		}
		numSyllables := options.numberOfSyllables()
		wordTemplate := syllable.GenerateTemplate(numSyllables)
		return generateNameFromSyllableTemplate(wordTemplate)
	}
}

// GenerateMaleName generates a random male Aslan name.
// By default, male names are 3-5 syllables long. This can be overridden by providing GeneratorOption arguments.
// Deprecated: use GenerateName(ctx, Male, opts...) instead.
func GenerateMaleName(ctx context.Context, opts ...GeneratorOption) (string, error) {
	return GenerateName(ctx, Male, opts...)
}

// GenerateFemaleName generates a random female Aslan name.
// By default, female names are 2-4 syllables long. This can be overridden by providing GeneratorOption arguments.
// Deprecated: use GenerateName(ctx, Female, opts...) instead.
func GenerateFemaleName(ctx context.Context, opts ...GeneratorOption) (string, error) {
	return GenerateName(ctx, Female, opts...)
}

// MustGenerate generates a random Aslan word with the given options.
// If no options are provided, it will generate-word a word with a random number of syllables between 2 and 6.
// If an error occurs, it will panic.
func MustGenerate(ctx context.Context, opts ...GeneratorOption) string {
	word, err := Generate(ctx, opts...)
	if err != nil {
		panic(fmt.Sprintf("error generating word: %v", err))
	}
	return word
}

// MustGenerateMaleName generates a random male Aslan name, panicking on error.
// By default, male names are 3-5 syllables long.
// Deprecated: use MustGenerateName(ctx, Male, opts...) instead.
func MustGenerateMaleName(ctx context.Context, opts ...GeneratorOption) string {
	// Calls the deprecated GenerateMaleName to maintain behavior if anyone was relying on its exact path,
	// though it now routes to GenerateName.
	name, err := GenerateMaleName(ctx, opts...)
	if err != nil {
		panic(fmt.Sprintf("error generating male name: %v", err))
	}
	return name
}

// MustGenerateFemaleName generates a random female Aslan name, panicking on error.
// By default, female names are 2-4 syllables long.
// Deprecated: use MustGenerateName(ctx, Female, opts...) instead.
func MustGenerateFemaleName(ctx context.Context, opts ...GeneratorOption) string {
	// Calls the deprecated GenerateFemaleName.
	name, err := GenerateFemaleName(ctx, opts...)
	if err != nil {
		panic(fmt.Sprintf("error generating female name: %v", err))
	}
	return name
}
