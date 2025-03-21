// Package syllable provides the generation of templates to generate valid aslan syllable sequences.
package syllable

import (
	"strings"
)

type template string

const (
	firstConsonant template = `<(f)|(f)|(f)|(f)|(f)|(ft)|(ft)|(ft)|(h)|(h)|(h)|(h)|(h)|(h)|(h)|(hf)|(hf)|(hk)|(hk)|(hk)|(hk)|(hk)|(hl)|(hl)|(hl)|(hr)|(hr)|(hr)|(ht)|(ht)|(ht)|(ht)|(ht)|(hw)|(hw)|(k)|(k)|(k)|(k)|(k)|(k)|(k)|(kh)|(kh)|(kh)|(kh)|(kh)|(kh)|(kht)|(kht)|(kht)|(kht)|(kt)|(kt)|(kt)|(kt)|(l)|(l)|(r)|(r)|(r)|(r)|(s)|(s)|(s)|(s)|(st)|(st)|(st)|(t)|(t)|(t)|(t)|(t)|(t)|(t)|(t)|(tl)|(tl)|(tr)|(tr)|(w)|(w)|(w)|(w)|(w)|(w)|>`
	vowel          template = "<(a)|(a)|(a)|(a)|(a)|(a)|(a)|(a)|(a)|(a)|(ai)|(ai)|(ai)|(ao)|(ao)|(au)|(e)|(e)|(e)|(e)|(e)|(e)|(ea)|(ea)|(ea)|(ea)|(ea)|(ea)|(ei)|(ei)|(i)|(i)|(i)|(i)|(iy)|(iy)|(iy)|(o)|(o)|(oa)|(oi)|(oi)|(ou)|(u)|(ua)|(ui)|(ya)|(yu)|>"
	lastConstant   template = "<(h)|(h)|(h)|(h)|(h)|(h)|(h)|(h)|(h)|(h)|(kh)|(kh)|(kh)|(kh)|(l)|(l)|(l)|(lr)|(lr)|(lr)|(lr)|(r)|(r)|(r)|(r)|(r)|(rl)|(rl)|(rl)|(rl)|(s)|(s)|(s)|(s)|(s)|(w)|(w)|(w)|(w)|(w)|(w)|(')|(')|(')>"

	singleA template = "(a)|"
	singleE template = "(e)|"
	singleI template = "(i)|"
	singleO template = "(o)|"
	singleU template = "(u)|"
)

var (
	vowelWithoutSingleA  = removeFromTemplate(vowel, singleA)
	vowelWithOnlySingleA = removeFromTemplate(vowel, singleE, singleI, singleO, singleU)

	vowelWithoutSingleE  = removeFromTemplate(vowel, singleE)
	vowelWithOnlySingleE = removeFromTemplate(vowel, singleA, singleI, singleO, singleU)

	vowelWithoutSingleI  = removeFromTemplate(vowel, singleI)
	vowelWithOnlySingleI = removeFromTemplate(vowel, singleA, singleE, singleO, singleU)

	vowelWithoutSingleO  = removeFromTemplate(vowel, singleO)
	vowelWithOnlySingleO = removeFromTemplate(vowel, singleA, singleE, singleI, singleU)

	vowelWithoutSingleU  = removeFromTemplate(vowel, singleU)
	vowelWithOnlySingleU = removeFromTemplate(vowel, singleA, singleE, singleI, singleO)
)

type swapKey string

const (
	withoutSingleA  swapKey = "withoutSingleA"
	withoutSingleE  swapKey = "withoutSingleE"
	withoutSingleI  swapKey = "withoutSingleI"
	withoutSingleO  swapKey = "withoutSingleO"
	withoutSingleU  swapKey = "withoutSingleU"
	withOnlySingleA swapKey = "withOnlySingleA"
	withOnlySingleE swapKey = "withOnlySingleE"
	withOnlySingleI swapKey = "withOnlySingleI"
	withOnlySingleO swapKey = "withOnlySingleO"
	withOnlySingleU swapKey = "withOnlySingleU"
)

type templateSwap struct {
	modifiedTemplate template
	reverseSwapKey   swapKey
}

// swaps are used to enforce no consecutive single vowels are generated. They are meant to work as puzzle pieces
var swaps = map[swapKey]templateSwap{
	withoutSingleA:  {vowelWithoutSingleA, withOnlySingleA},
	withoutSingleE:  {vowelWithoutSingleE, withOnlySingleE},
	withoutSingleI:  {vowelWithoutSingleI, withOnlySingleI},
	withoutSingleO:  {vowelWithoutSingleO, withOnlySingleO},
	withoutSingleU:  {vowelWithoutSingleU, withOnlySingleU},
	withOnlySingleA: {vowelWithOnlySingleA, withoutSingleA},
	withOnlySingleE: {vowelWithOnlySingleE, withoutSingleE},
	withOnlySingleI: {vowelWithOnlySingleI, withoutSingleI},
	withOnlySingleO: {vowelWithOnlySingleO, withoutSingleO},
	withOnlySingleU: {vowelWithOnlySingleU, withoutSingleU},
}

var allSwaps = []swapKey{
	withoutSingleA, withoutSingleE, withoutSingleI, withoutSingleO, withoutSingleU,
	withOnlySingleA, withOnlySingleE, withOnlySingleI, withOnlySingleO, withOnlySingleU,
}

type syllableKey string

func (k syllableKey) StartsWithConsonant() bool {
	return k[0] == 'c'
}

func (k syllableKey) EndsWithConsonant() bool {
	return k[len(k)-1] == 'c'
}

const (
	keyV   = "v"
	keyCV  = "cv"
	keyVC  = "vc"
	keyCVC = "cvc"
)

type syllableDefinition interface {
	Key() syllableKey
	Weight() int
	Template() template
	EnforceNoConsecutiveSingleVowels(nextSyllable syllableDefinition, generateRandomSwapVowelFn GenerateRandomIntegerUpToFn)
	SwapVowelTemplate(swap templateSwap)
	SyllablesThatCanFollowThis() []syllableDefinition
	StartsWithConsonant() bool
}

type syllable struct {
	key                          syllableKey
	weight                       int
	vowelSwap                    *templateSwap
	syllablesThatCanFollowThisFn func() []syllableDefinition
}

func newV() *syllable {
	return &syllable{
		key:                          keyV,
		weight:                       3,
		syllablesThatCanFollowThisFn: allSyllables,
	}
}

func newCV() *syllable {
	return &syllable{
		key:                          keyCV,
		weight:                       3,
		syllablesThatCanFollowThisFn: allSyllables,
	}
}

func newVC() *syllable {
	return &syllable{
		key:                          keyVC,
		weight:                       2,
		syllablesThatCanFollowThisFn: onlyVowelStartingSyllables,
	}
}

func newCVC() *syllable {
	return &syllable{
		key:                          keyCVC,
		weight:                       2,
		syllablesThatCanFollowThisFn: onlyVowelStartingSyllables,
	}
}

func (d *syllable) Key() syllableKey {
	return d.key
}

func (d *syllable) Weight() int {
	return d.weight
}

func (d *syllable) vowelTemplate() template {
	if d.vowelSwap == nil {
		return vowel
	}
	return d.vowelSwap.modifiedTemplate
}

func (d *syllable) Template() template {
	templateBuilder := new(strings.Builder)
	for i, char := range d.key {
		switch char {
		case 'c':
			if i == 0 {
				templateBuilder.WriteString(string(firstConsonant))
			} else {
				templateBuilder.WriteString(string(lastConstant))
			}
		case 'v':
			templateBuilder.WriteString(string(d.vowelTemplate()))
		}
	}

	return template(templateBuilder.String())
}

func (d *syllable) EnforceNoConsecutiveSingleVowels(nextSyllable syllableDefinition, generateRandomSwapVowelFn GenerateRandomIntegerUpToFn) {
	if d.key.EndsWithConsonant() || nextSyllable.StartsWithConsonant() {
		return
	}
	if d.vowelSwap == nil {
		d.vowelSwap = pickRandomSwap(generateRandomSwapVowelFn)
	}
	nextSyllable.SwapVowelTemplate(swaps[d.vowelSwap.reverseSwapKey])
}

func (d *syllable) SwapVowelTemplate(swap templateSwap) {
	d.vowelSwap = &swap
}

func (d *syllable) SyllablesThatCanFollowThis() []syllableDefinition {
	return d.syllablesThatCanFollowThisFn()
}

func (d *syllable) StartsWithConsonant() bool {
	return d.key.StartsWithConsonant()
}

func pickRandomSwap(randomIndexPicker GenerateRandomIntegerUpToFn) *templateSwap {
	chosenSwapIndex := randomIndexPicker(len(swaps))
	vowelTemplateSwapKey := allSwaps[chosenSwapIndex]
	swap := swaps[vowelTemplateSwapKey]
	return &swap
}

func removeFromTemplate(sourceTemplate template, templatesToRemove ...template) template {
	newTemplate := string(sourceTemplate)
	for _, templateToRemove := range templatesToRemove {
		newTemplate = strings.ReplaceAll(newTemplate, string(templateToRemove), "")
	}
	return template(newTemplate)
}

func allSyllables() []syllableDefinition {
	return []syllableDefinition{newV(), newCV(), newVC(), newCVC()}
}

func onlyVowelStartingSyllables() []syllableDefinition {
	return []syllableDefinition{newV(), newVC()}
}
