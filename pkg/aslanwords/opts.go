package aslanwords

import (
	"fmt"
	"math/rand"
)

// WithNumberOfSyllables sets the number of syllables to generate-word
func WithNumberOfSyllables(n int) GeneratorOption {
	return func(o *GeneratorOptions) {
		o.numberOfSyllablesOpts = fixedAmountOpt{numberOfSyllables: n}
	}
}

// WithNumberOfSyllablesBetween Use it to generate-word a random number of syllables between the 'from' and 'to' values
func WithNumberOfSyllablesBetween(from, to int) GeneratorOption {
	return func(o *GeneratorOptions) {
		o.numberOfSyllablesOpts = randomAmountOpt{from: from, to: to}
	}
}

// GeneratorOption Option to configure the generation of Aslan words
type GeneratorOption func(*GeneratorOptions)

// GeneratorOptions Options to configure the generation of Aslan words
type GeneratorOptions struct {
	numberOfSyllablesOpts amountOptions
}

func newGeneratorOptions() *GeneratorOptions {
	const defaultMinNumberOfSyllables = 2
	const defaultMaxNumberOfSyllables = 6

	opts := &GeneratorOptions{}
	WithNumberOfSyllablesBetween(defaultMinNumberOfSyllables, defaultMaxNumberOfSyllables)(opts)

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
