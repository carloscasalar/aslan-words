// Package syllable provides the generation of templates to generate valid aslan syllable sequences.
package syllable

import (
	"fmt"
	"math/rand/v2"
)

/*
	   The Aslan language has four types of syllables:
		 - those consisting of a single vowel (V),
		 - those beginning with a consonant (CV),
		 - those that ends with a consonant (VC),
		 - and those that begin and end with a consonant (CVC).

	   In Aslan words, the frequency of occurrence is, out of ten syllables,
		 - three will be V,
		 - three will be CV,
		 - two will be VC,
		 - and two will be CVC.

	   No syllable ending with a consonant can be followed by a syllable that starts with a consonant, so
	   A syllable of a single vowel cannot be directly followed by same syllable. For example, 'aa' cannot be
	   but 'aeae' or 'aeei' would be OK.

	   So:
		 - V can be followed by V, CV, VC or CVC.
		 - CV can be followed by V, CV, VC or CVC.
		 - VC can be followed by V or VC.
		 - CVC can be followed by V or VC.
	     - following syllable sequences are not allowed: aa, ee, ii, oo, uu
*/
type template string

const (
	firstConsonant template = `<(f)|(f)|(f)|(f)|(f)|(ft)|(ft)|(ft)|(h)|(h)|(h)|(h)|(h)|(h)|(h)|(hf)|(hf)|(hk)|(hk)|(hk)|(hk)|(hk)|(hl)|(hl)|(hl)|(hr)|(hr)|(hr)|(ht)|(ht)|(ht)|(ht)|(ht)|(hw)|(hw)|(k)|(k)|(k)|(k)|(k)|(k)|(k)|(kh)|(kh)|(kh)|(kh)|(kh)|(kh)|(kht)|(kht)|(kht)|(kht)|(kt)|(kt)|(kt)|(kt)|(l)|(l)|(r)|(r)|(r)|(r)|(s)|(s)|(s)|(s)|(st)|(st)|(st)|(t)|(t)|(t)|(t)|(t)|(t)|(t)|(t)|(tl)|(tl)|(tr)|(tr)|(w)|(w)|(w)|(w)|(w)|(w)|>`
	vowel          template = "<(a)|(a)|(a)|(a)|(a)|(a)|(a)|(a)|(a)|(a)|(ai)|(ai)|(ai)|(ao)|(ao)|(au)|(e)|(e)|(e)|(e)|(e)|(e)|(ea)|(ea)|(ea)|(ea)|(ea)|(ea)|(ei)|(ei)|(i)|(i)|(i)|(i)|(iy)|(iy)|(iy)|(o)|(o)|(o)|(oa)|(oi)|(oi)|(ou)|(u)|(ua)|(ui)|(ya)|(yu)|>"
	lastConstant   template = "<(h)|(h)|(h)|(h)|(h)|(h)|(h)|(h)|(h)|(h)|(kh)|(kh)|(kh)|(kh)|(l)|(l)|(l)|(lr)|(lr)|(lr)|(lr)|(r)|(r)|(r)|(r)|(r)|(rl)|(rl)|(rl)|(rl)|(s)|(s)|(s)|(s)|(s)|(w)|(w)|(w)|(w)|(w)|(w)|(')|(')|(')>"
)

type syllableKey string

const (
	keyV   = "v"
	keyCV  = "cv"
	keyVC  = "vc"
	keyCVC = "cvc"
)

var (
	anySyllableKeys               = []syllableKey{keyV, keyCV, keyVC, keyCVC}
	onlyVowelStartingSyllableKeys = []syllableKey{keyV, keyCV}
)

type syllableDefinition struct {
	key syllableKey
	// pattern template of the syllable
	pattern template
	// weight of the syllable over 10
	weight int
	// followedBy allowed syllables
	followedByKeys []syllableKey
}

// SyllablesThatCanFollowThis returns the syllables that can follow this syllable
func (sd syllableDefinition) SyllablesThatCanFollowThis() []syllableDefinition {
	syllableDefinitionByKey := map[syllableKey]syllableDefinition{
		keyV:   v,
		keyCV:  cv,
		keyVC:  vc,
		keyCVC: cvc,
	}
	followedBy := make([]syllableDefinition, 0, len(sd.followedByKeys))
	for _, key := range sd.followedByKeys {
		followedBy = append(followedBy, syllableDefinitionByKey[key])
	}
	return followedBy
}

var (
	v   = syllableDefinition{keyV, vowel, 3, anySyllableKeys}
	cv  = syllableDefinition{keyCV, firstConsonant + vowel, 3, anySyllableKeys}
	vc  = syllableDefinition{keyVC, vowel + lastConstant, 2, onlyVowelStartingSyllableKeys}
	cvc = syllableDefinition{keyCVC, firstConsonant + vowel + lastConstant, 2, onlyVowelStartingSyllableKeys}
)

// TemplateDefinition is a sequence of syllables that can be used to generate-word an Aslan word
type TemplateDefinition []syllableDefinition

func (td TemplateDefinition) String() string {
	s := ""
	for _, sd := range td {
		s = fmt.Sprintf("%s%s", s, sd.pattern)
	}
	return s
}

// GenerateTemplate generates a template with the given number of syllables for an Aslan word
// built with the rules of https://github.com/s0rg/fantasyname?tab=readme-ov-file#pattern-syntax
// and the Aslan language rules
func GenerateTemplate(numberOfSyllables int, opts ...TemplateOption) TemplateDefinition {
	if numberOfSyllables < 1 {
		return nil
	}
	options := applyTemplateOptions(opts...)
	sequenceGenerator := newSyllableSequenceBuilder(options.integerGenerator)
	return sequenceGenerator.randomSyllableSequence(numberOfSyllables)
}

type syllableSequenceBuilder struct {
	generateRandomIntegerUpTo GenerateRandomIntegerUpToFn
	lastSyllableGenerated     *syllableDefinition
}

func newSyllableSequenceBuilder(integerGenerator GenerateRandomIntegerUpToFn) *syllableSequenceBuilder {
	return &syllableSequenceBuilder{generateRandomIntegerUpTo: integerGenerator}
}

func (b *syllableSequenceBuilder) randomSyllableSequence(numberOfSyllables int, previousSyllables ...syllableDefinition) []syllableDefinition {
	if numberOfSyllables < 1 {
		return previousSyllables
	}
	if len(previousSyllables) == 0 {
		return b.randomSyllableSequence(numberOfSyllables-1, b.pickRandomSyllable([]syllableDefinition{v, cv, vc, cvc}))
	}
	lastSyllable := previousSyllables[len(previousSyllables)-1]
	nextSyllable := b.pickRandomSyllable(lastSyllable.SyllablesThatCanFollowThis())
	return b.randomSyllableSequence(numberOfSyllables-1, append(previousSyllables, nextSyllable)...)
}

func (b *syllableSequenceBuilder) pickRandomSyllable(definitions []syllableDefinition) syllableDefinition {
	totalWeight := 0
	for _, def := range definitions {
		totalWeight += def.weight
	}
	r := b.generateRandomIntegerUpTo(totalWeight)
	for _, def := range definitions {
		if r < def.weight {
			return def
		}
		r -= def.weight
	}
	return definitions[len(definitions)-1]
}

// GenerateRandomIntegerUpToFn is a function that is expected to generate a positive integer from zero up to the given number
type GenerateRandomIntegerUpToFn func(int) int

// TemplateOption is a function that sets options for the template generation
type TemplateOption func(*templateOptions)

type templateOptions struct {
	integerGenerator GenerateRandomIntegerUpToFn
}

// WithIntegerGenerator sets the random number generator for the template generation
func WithIntegerGenerator(fn GenerateRandomIntegerUpToFn) TemplateOption {
	return func(o *templateOptions) {
		o.integerGenerator = fn
	}
}

func applyTemplateOptions(opts ...TemplateOption) *templateOptions {
	opt := &templateOptions{
		integerGenerator: rand.IntN, // default random number generator
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}
