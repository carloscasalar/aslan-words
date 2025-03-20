// Package syllable provides the generation of templates to generate valid aslan syllable sequences.
package syllable

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
	onlyVowelStartingSyllableKeys = []syllableKey{keyV, keyVC}
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

// Pattern returns the pattern of the syllable
func (sd syllableDefinition) Pattern() string {
	return string(sd.pattern)
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
