// Package syllable provides the generation of templates to generate valid aslan syllable sequences.
package syllable

import "strings"

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

const (
	singleA template = "(a)|"
	singleE template = "(e)|"
	singleI template = "(i)|"
	singleO template = "(o)|"
	singleU template = "(u)|"
)

type syllableDefinition interface {
	Key() syllableKey
	Weight() int
	Template() template
	RemoveVowelsFromTemplate(vowels []template) syllableDefinition
	SyllablesThatCanFollowThis() []syllableDefinition
}

var (
	v   = newV()
	cv  = newCV()
	vc  = newVC()
	cvc = newCVC()
)

type vowelSyllableDefinition struct {
	template template
}

func newV() vowelSyllableDefinition {
	return vowelSyllableDefinition{template: vowel}
}

func (d vowelSyllableDefinition) Key() syllableKey {
	return keyV
}

func (d vowelSyllableDefinition) Weight() int {
	const chanceOfVowelSyllableOutOf10 = 3
	return chanceOfVowelSyllableOutOf10
}

func (d vowelSyllableDefinition) Template() template {
	return d.template
}

func (d vowelSyllableDefinition) RemoveVowelsFromTemplate(vowels []template) syllableDefinition {
	return vowelSyllableDefinition{vowelTemplateWithoutSingleVowels(vowels)}
}

func (d vowelSyllableDefinition) SyllablesThatCanFollowThis() []syllableDefinition {
	return []syllableDefinition{v, cv, vc, cvc}
}

type consonantVowelSyllableDefinition struct {
	template template
}

func newCV() consonantVowelSyllableDefinition {
	return consonantVowelSyllableDefinition{template: firstConsonant + vowel}
}

func (d consonantVowelSyllableDefinition) Key() syllableKey {
	return keyCV
}

func (d consonantVowelSyllableDefinition) Weight() int {
	const chanceOfCVSyllableOutOf10 = 3
	return chanceOfCVSyllableOutOf10
}

func (d consonantVowelSyllableDefinition) Template() template {
	return d.template
}

func (d consonantVowelSyllableDefinition) RemoveVowelsFromTemplate(vowels []template) syllableDefinition {
	return vowelSyllableDefinition{firstConsonant + vowelTemplateWithoutSingleVowels(vowels)}
}

func (d consonantVowelSyllableDefinition) SyllablesThatCanFollowThis() []syllableDefinition {
	return []syllableDefinition{v, cv, vc, cvc}
}

type vowelConsonantSyllableDefinition struct {
	template template
}

func newVC() vowelConsonantSyllableDefinition {
	return vowelConsonantSyllableDefinition{template: vowel + lastConstant}
}

func (d vowelConsonantSyllableDefinition) Key() syllableKey {
	return keyVC
}

func (d vowelConsonantSyllableDefinition) Weight() int {
	const chanceOfVCSyllableOutOf10 = 2
	return chanceOfVCSyllableOutOf10
}

func (d vowelConsonantSyllableDefinition) Template() template {
	return d.template
}

func (d vowelConsonantSyllableDefinition) RemoveVowelsFromTemplate(vowels []template) syllableDefinition {
	return vowelSyllableDefinition{vowelTemplateWithoutSingleVowels(vowels) + lastConstant}
}

func (d vowelConsonantSyllableDefinition) SyllablesThatCanFollowThis() []syllableDefinition {
	return []syllableDefinition{v, vc}
}

type vowelConsonantVowelSyllableDefinition struct {
	template template
}

func newCVC() vowelConsonantVowelSyllableDefinition {
	return vowelConsonantVowelSyllableDefinition{firstConsonant + vowel + lastConstant}
}

func (d vowelConsonantVowelSyllableDefinition) Key() syllableKey {
	return keyCVC
}

func (d vowelConsonantVowelSyllableDefinition) Weight() int {
	const chanceOfCVCSyllableOutOf10 = 2
	return chanceOfCVCSyllableOutOf10
}

func (d vowelConsonantVowelSyllableDefinition) Template() template {
	return d.template
}

func (d vowelConsonantVowelSyllableDefinition) RemoveVowelsFromTemplate(vowels []template) syllableDefinition {
	return vowelSyllableDefinition{firstConsonant + vowelTemplateWithoutSingleVowels(vowels) + lastConstant}
}

func (d vowelConsonantVowelSyllableDefinition) SyllablesThatCanFollowThis() []syllableDefinition {
	return []syllableDefinition{v, vc}
}

func vowelTemplateWithoutSingleVowels(singleVowels []template) template {
	if len(singleVowels) == 0 {
		return vowel
	}
	newVowelPattern := string(vowel)
	for _, singleVowel := range singleVowels {
		newVowelPattern = strings.ReplaceAll(newVowelPattern, string(singleVowel), "")
	}
	return template(newVowelPattern)
}
