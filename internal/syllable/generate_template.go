package syllable

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// TemplateDefinition is a sequence of syllables that can be used to generate-word an Aslan word
type TemplateDefinition []syllableDefinition

func (td TemplateDefinition) String() string {
	s := ""
	for _, sd := range td {
		s = fmt.Sprintf("%s%s", s, sd.Template())
	}
	return s
}

// TemplateSequence returns the sequence of template syllables in the template
func (td TemplateDefinition) TemplateSequence() []string {
	sequence := make([]string, len(td))
	for i, sd := range td {
		sequence[i] = string(sd.Template())
	}
	return sequence
}

// SyllableKeySequence returns the sequence of keys of the syllables in the template where:
// - `V` is a syllable made of an aslan vowel
// - `CV` is an aslan consonant followed by an aslan vowel
// - `VC` is an aslan vowel followed by an aslan consonant
// - `CVC` is an aslan consonant followed by an aslan vowel and ending with an aslan consonant
func (td TemplateDefinition) SyllableKeySequence() []string {
	sequence := make([]string, len(td))
	for i, sd := range td {
		sequence[i] = strings.ToUpper(string(sd.Key()))
	}
	return sequence
}

// GenerateTemplate generates a template with the given number of syllables for an Aslan word
// built with the rules of https://github.com/s0rg/fantasyname?tab=readme-ov-file#pattern-syntax
// and the Aslan language rules
func GenerateTemplate(numberOfSyllables int, opts ...TemplateOption) TemplateDefinition {
	if numberOfSyllables < 1 {
		return nil
	}
	options := applyTemplateOptions(opts...)
	sequenceGenerator := newSyllableSequenceBuilder(options)
	return sequenceGenerator.randomSyllableSequence(numberOfSyllables)
}

type syllableSequenceBuilder struct {
	generateRandomIntegerUpTo    GenerateRandomIntegerUpToFn
	vowelTemplateChanceGenerator GenerateRandomIntegerUpToFn
	lastSyllableGenerated        syllableDefinition
}

func newSyllableSequenceBuilder(opt *templateOptions) *syllableSequenceBuilder {
	return &syllableSequenceBuilder{generateRandomIntegerUpTo: opt.syllableChanceGenerator, vowelTemplateChanceGenerator: opt.vowelTemplateChanceGenerator}
}

func (b *syllableSequenceBuilder) randomSyllableSequence(numberOfSyllables int, previousSyllables ...syllableDefinition) []syllableDefinition {
	if numberOfSyllables < 1 {
		return previousSyllables
	}
	if len(previousSyllables) == 0 {
		return b.randomSyllableSequence(numberOfSyllables-1, b.pickRandomSyllable(allSyllables()))
	}
	lastSyllable := previousSyllables[len(previousSyllables)-1]
	nextSyllable := b.pickRandomSyllable(lastSyllable.SyllablesThatCanFollowThis())
	lastSyllable.EnforceNoConsecutiveSingleVowels(nextSyllable, b.vowelTemplateChanceGenerator)
	return b.randomSyllableSequence(numberOfSyllables-1, append(previousSyllables, nextSyllable)...)
}

func (b *syllableSequenceBuilder) pickRandomSyllable(definitions []syllableDefinition) syllableDefinition {
	totalWeight := 0
	for _, def := range definitions {
		totalWeight += def.Weight()
	}
	chance := b.generateRandomIntegerUpTo(totalWeight)
	for _, def := range definitions {
		if chance < def.Weight() {
			return def
		}
		chance -= def.Weight()
	}
	return definitions[len(definitions)-1]
}

// GenerateRandomIntegerUpToFn is a function that is expected to generate a positive integer from zero up to the given number minus one
type GenerateRandomIntegerUpToFn func(int) int

// TemplateOption is a function that sets options for the template generation
type TemplateOption func(*templateOptions)

type templateOptions struct {
	syllableChanceGenerator      GenerateRandomIntegerUpToFn
	vowelTemplateChanceGenerator GenerateRandomIntegerUpToFn
}

// WithSyllableChanceGenerator sets the random number generator to choose the syllable using its weight over all the possible syllables
func WithSyllableChanceGenerator(fn GenerateRandomIntegerUpToFn) TemplateOption {
	return func(o *templateOptions) {
		o.syllableChanceGenerator = fn
	}
}

// WithVowelTemplateChanceGenerator sets the random number generator to choose the vowel using its weight over all the possible vowels
func WithVowelTemplateChanceGenerator(fn GenerateRandomIntegerUpToFn) TemplateOption {
	return func(o *templateOptions) {
		o.vowelTemplateChanceGenerator = fn
	}
}

func applyTemplateOptions(opts ...TemplateOption) *templateOptions {
	opt := &templateOptions{
		syllableChanceGenerator:      rand.IntN,
		vowelTemplateChanceGenerator: rand.IntN,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}
