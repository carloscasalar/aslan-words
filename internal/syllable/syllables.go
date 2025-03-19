package syllable

import (
	"fmt"
	"math/rand"
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

type TemplateDefinition []syllableDefinition

func (td TemplateDefinition) String() string {
	s := ""
	for _, sd := range td {
		s = fmt.Sprintf("%s%s", s, sd.pattern)
	}
	return s
}

func GenerateTemplate(numberOfSyllables int) TemplateDefinition {
	if numberOfSyllables < 1 {
		return nil
	}
	return generateSyllables(numberOfSyllables, pickRandomSyllable([]syllableDefinition{v, cv, vc, cvc}))
}

func generateSyllables(remaining int, lastSyllable syllableDefinition) TemplateDefinition {
	if remaining == 1 {
		return TemplateDefinition{lastSyllable}
	}
	followedBy := lastSyllable.SyllablesThatCanFollowThis()
	nextSyllable := pickRandomSyllable(followedBy)
	return append(TemplateDefinition{lastSyllable}, generateSyllables(remaining-1, nextSyllable)...)
}

func pickRandomSyllable(definitions []syllableDefinition) syllableDefinition {
	totalWeight := 0
	for _, def := range definitions {
		totalWeight += def.weight
	}
	r := rand.Intn(totalWeight)
	for _, def := range definitions {
		if r < def.weight {
			return def
		}
		r -= def.weight
	}
	return definitions[len(definitions)-1]
}
