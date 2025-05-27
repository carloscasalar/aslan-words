package aslanwords

import (
	"fmt"
	"math/rand"
)

// WordType defines the type of word to be generated.
type WordType int

// Gender defines the gender for name generation.
type Gender int

const (
	// TypeWord generates a generic Aslan word.
	TypeWord WordType = iota
	// TypeMaleName signals generation of an Aslan male name.
	TypeMaleName
	// TypeFemaleName signals generation of an Aslan female name.
	TypeFemaleName
)

const (
	// Male specifies male gender for name generation.
	Male Gender = iota
	// Female specifies female gender for name generation.
	Female
)

// WithType sets the type of word to generate.
// If the type is TypeMaleName or TypeFemaleName, it also sets the corresponding gender.
func WithType(wordType WordType) GeneratorOption {
	return func(o *GeneratorOptions) {
		o.wordType = wordType
		switch wordType {
		case TypeMaleName:
			o.gender = Male
		case TypeFemaleName:
			o.gender = Female
		}
	}
}

// WithNumberOfSyllables sets the number of syllables to generate-word
func WithNumberOfSyllables(n int) GeneratorOption {
	return func(o *GeneratorOptions) {
		o.numberOfSyllablesOpts = fixedAmountOpt{numberOfSyllables: n}
		o.isSyllableCountUserSet = true
	}
}

// WithNumberOfSyllablesBetween Use it to generate-word a random number of syllables between the 'from' and 'to' values
func WithNumberOfSyllablesBetween(from, to int) GeneratorOption {
	return func(o *GeneratorOptions) {
		o.numberOfSyllablesOpts = randomAmountOpt{from: from, to: to}
		o.isSyllableCountUserSet = true
	}
}

// GeneratorOption Option to configure the generation of Aslan words
type GeneratorOption func(*GeneratorOptions)

// GeneratorOptions Options to configure the generation of Aslan words
type GeneratorOptions struct {
	numberOfSyllablesOpts  amountOptions
	wordType               WordType
	isSyllableCountUserSet bool // Tracks if syllable count was explicitly set by the user.
	gender                 Gender
}

func newGeneratorOptions() *GeneratorOptions {
	const defaultMinNumberOfSyllables = 2
	const defaultMaxNumberOfSyllables = 6

	opts := &GeneratorOptions{
		wordType:               TypeWord, // Default to generating a generic word
		isSyllableCountUserSet: false,    // Initially false
		// gender remains its zero value (Male, due to iota) by default,
		// which is acceptable as it's only explicitly used when wordType indicates a name.
		// Or, it's explicitly set by WithType or GenerateName.
	}
	opts.numberOfSyllablesOpts = randomAmountOpt{from: defaultMinNumberOfSyllables, to: defaultMaxNumberOfSyllables}

	return opts
}

// Validate checks if the options are valid, returning an error if not
func (o *GeneratorOptions) Validate() error {
	if err := o.numberOfSyllablesOpts.Validate(); err != nil {
		return err
	}
	return nil
}

// numberOfSyllables returns the number of syllables to generate-word
func (o *GeneratorOptions) numberOfSyllables() int {
	if o.numberOfSyllablesOpts == nil {
		return 0
	}
	return o.numberOfSyllablesOpts.NumberOfSyllables()
}

type amountOptions interface {
	Validate() error
	NumberOfSyllables() int
}
type fixedAmountOpt struct {
	numberOfSyllables int
}

func (s fixedAmountOpt) Validate() error {
	if s.numberOfSyllables < 1 {
		return fmt.Errorf("number of syllables must be one or greater")
	}
	return nil
}

func (s fixedAmountOpt) NumberOfSyllables() int {
	return s.numberOfSyllables
}

type randomAmountOpt struct {
	from int
	to   int
}

func (r randomAmountOpt) Validate() error {
	if r.from < 1 {
		return fmt.Errorf("minimum number of syllables must be one or greater")
	}
	if r.from > r.to {
		return fmt.Errorf("number of syllables 'from' cannot be greater than 'to'")
	}
	if r.to >= 15 {
		return fmt.Errorf("number of syllables 'to' cannot be greater than 15")
	}
	return nil
}

// NumberOfSyllables returns a random number of syllables between the 'from' and 'to' values
func (r randomAmountOpt) NumberOfSyllables() int {
	return r.from + rand.Intn(r.to-r.from)
}
