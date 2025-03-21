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
	reverseSwap      swapKey
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

const (
	keyV   = "v"
	keyCV  = "cv"
	keyVC  = "vc"
	keyCVC = "cvc"
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
	EnforceNoConsecutiveSingleVowels(nextSyllable syllableDefinition, generateRandomSwapVowelFn GenerateRandomIntegerUpToFn)
	SwapVowelTemplate(swap templateSwap)
	SyllablesThatCanFollowThis() []syllableDefinition
	StartsWithConsonant() bool
}

type vowelSyllableDefinition struct {
	vowelSwap *templateSwap
}

func newV() *vowelSyllableDefinition {
	return &vowelSyllableDefinition{
		vowelSwap: nil,
	}
}

func (d *vowelSyllableDefinition) Key() syllableKey {
	return keyV
}

func (d *vowelSyllableDefinition) Weight() int {
	const chanceOfVowelSyllableOutOf10 = 3
	return chanceOfVowelSyllableOutOf10
}

func (d *vowelSyllableDefinition) Template() template {
	if d.vowelSwap == nil {
		return vowel
	}
	return d.vowelSwap.modifiedTemplate
}

func (d *vowelSyllableDefinition) SwapVowelTemplate(swap templateSwap) {
	d.vowelSwap = &swap
}

func (d *vowelSyllableDefinition) EnforceNoConsecutiveSingleVowels(nextSyllable syllableDefinition, generateRandomSwapVowelFn GenerateRandomIntegerUpToFn) {
	if nextSyllable.StartsWithConsonant() {
		return
	}
	if d.vowelSwap == nil {
		d.vowelSwap = pickRandomSwap(generateRandomSwapVowelFn)
	}
	nextSyllable.SwapVowelTemplate(swaps[d.vowelSwap.reverseSwap])
}

func (d *vowelSyllableDefinition) SyllablesThatCanFollowThis() []syllableDefinition {
	return allSyllables()
}

func (d *vowelSyllableDefinition) StartsWithConsonant() bool {
	return false
}

type consonantVowelSyllableDefinition struct {
	vowelSwap *templateSwap
}

func newCV() *consonantVowelSyllableDefinition {
	return &consonantVowelSyllableDefinition{}
}

func (d *consonantVowelSyllableDefinition) Key() syllableKey {
	return keyCV
}

func (d *consonantVowelSyllableDefinition) Weight() int {
	const chanceOfCVSyllableOutOf10 = 3
	return chanceOfCVSyllableOutOf10
}

func (d *consonantVowelSyllableDefinition) Template() template {
	if d.vowelSwap == nil {
		return firstConsonant + vowel
	}
	return firstConsonant + d.vowelSwap.modifiedTemplate
}

func (d *consonantVowelSyllableDefinition) SwapVowelTemplate(swap templateSwap) {
	d.vowelSwap = &swap
}

func (d *consonantVowelSyllableDefinition) EnforceNoConsecutiveSingleVowels(nextSyllable syllableDefinition, generateRandomVowelFn GenerateRandomIntegerUpToFn) {
	if nextSyllable.StartsWithConsonant() {
		return
	}
	if d.vowelSwap == nil {
		d.vowelSwap = pickRandomSwap(generateRandomVowelFn)
	}
	nextSyllable.SwapVowelTemplate(swaps[d.vowelSwap.reverseSwap])
}

func (d *consonantVowelSyllableDefinition) SyllablesThatCanFollowThis() []syllableDefinition {
	return allSyllables()
}

func (d *consonantVowelSyllableDefinition) StartsWithConsonant() bool {
	return true
}

type vowelConsonantSyllableDefinition struct {
	vowelSwap *templateSwap
}

func newVC() *vowelConsonantSyllableDefinition {
	return &vowelConsonantSyllableDefinition{}
}

func (d *vowelConsonantSyllableDefinition) Key() syllableKey {
	return keyVC
}

func (d *vowelConsonantSyllableDefinition) Weight() int {
	const chanceOfVCSyllableOutOf10 = 2
	return chanceOfVCSyllableOutOf10
}

func (d *vowelConsonantSyllableDefinition) Template() template {
	if d.vowelSwap == nil {
		return vowel + lastConstant
	}
	return d.vowelSwap.modifiedTemplate + lastConstant
}

func (d *vowelConsonantSyllableDefinition) SwapVowelTemplate(swap templateSwap) {
	d.vowelSwap = &swap
}

func (d *vowelConsonantSyllableDefinition) EnforceNoConsecutiveSingleVowels(_ syllableDefinition, _ GenerateRandomIntegerUpToFn) {
	// VC ends with a consonant, so no restrictions with following syllables
}

func (d *vowelConsonantSyllableDefinition) SyllablesThatCanFollowThis() []syllableDefinition {
	return onlyVowelStartingSyllables()
}

func (d *vowelConsonantSyllableDefinition) StartsWithConsonant() bool {
	return false
}

type consonantVowelConsonantSyllableDefinition struct {
	vowelSwap *templateSwap
}

func (d *consonantVowelConsonantSyllableDefinition) SwapVowelTemplate(swap templateSwap) {
	d.vowelSwap = &swap
}

func newCVC() *consonantVowelConsonantSyllableDefinition {
	return &consonantVowelConsonantSyllableDefinition{}
}

func (d *consonantVowelConsonantSyllableDefinition) Key() syllableKey {
	return keyCVC
}

func (d *consonantVowelConsonantSyllableDefinition) Weight() int {
	const chanceOfCVCSyllableOutOf10 = 2
	return chanceOfCVCSyllableOutOf10
}

func (d *consonantVowelConsonantSyllableDefinition) Template() template {
	if d.vowelSwap == nil {
		return firstConsonant + vowel + lastConstant
	}

	return firstConsonant + d.vowelSwap.modifiedTemplate + lastConstant
}

func (d *consonantVowelConsonantSyllableDefinition) EnforceNoConsecutiveSingleVowels(_ syllableDefinition, _ GenerateRandomIntegerUpToFn) {
	// CVC ends with a consonant, so no restrictions with following syllables
}

func (d *consonantVowelConsonantSyllableDefinition) SyllablesThatCanFollowThis() []syllableDefinition {
	return onlyVowelStartingSyllables()
}

func (d *consonantVowelConsonantSyllableDefinition) StartsWithConsonant() bool {
	return true
}

func pickRandomSwap(randomIndexPicker GenerateRandomIntegerUpToFn) *templateSwap {
	chosenSwapIndex := randomIndexPicker(len(swaps))
	swapKey := allSwaps[chosenSwapIndex]
	swap := swaps[swapKey]
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
