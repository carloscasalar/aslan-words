package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/carloscasalar/aslan-words/pkg/aslanwords"
	"github.com/jessevdk/go-flags"
)

// CliOptions holds the command-line options passed to the application.
type CliOptions struct {
	NumberOfSyllables int    `short:"s" long:"syllables" default:"2" description:"Number of syllables of the word to generate. For 'male' or 'female' types, this will override their default syllable counts."`
	WordType          string `long:"type" default:"word" description:"Type of word to generate (word, male, female)."`
}

func readOptionsOrFail() CliOptions {
	var opts CliOptions
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(flags.ErrorType); ok && flagsErr == flags.ErrHelp {
			os.Exit(0)
		}
		// For other errors, the library prints a message, so we just exit.
		os.Exit(1)
	}
	return opts
}

// ToAslanGeneratorOptions converts CLI options to the aslanwords.GeneratorOption slice.
func (co *CliOptions) ToAslanGeneratorOptions() ([]aslanwords.GeneratorOption, error) {
	var genOpts []aslanwords.GeneratorOption

	// Convert and add WordType option
	var awType aslanwords.WordType
	switch strings.ToLower(co.WordType) {
	case "word":
		awType = aslanwords.TypeWord
	case "male":
		awType = aslanwords.TypeMaleName
	case "female":
		awType = aslanwords.TypeFemaleName
	default:
		return nil, fmt.Errorf("invalid word type '%s'. Valid types are: word, male, female", co.WordType)
	}
	genOpts = append(genOpts, aslanwords.WithType(awType))

	// Add NumberOfSyllables option.
	// The default:"2" in CliOptions means NumberOfSyllables will always have a value.
	// aslanwords.WithNumberOfSyllables() sets isSyllableCountUserSet=true,
	// which means this will always override the gender-specific syllable defaults if type is male/female.
	// This is acceptable for the CLI's behavior.
	// If NumberOfSyllables was 0 or negative (not possible with current default and positive user input),
	// it might indicate "not set", but the library ensures it's at least the default.
	genOpts = append(genOpts, aslanwords.WithNumberOfSyllables(co.NumberOfSyllables))

	return genOpts, nil
}
