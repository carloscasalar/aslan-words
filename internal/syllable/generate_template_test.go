package syllable_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/carloscasalar/aslan-words/internal/syllable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateTemplate_out_of_ten_syllables(t *testing.T) {
	testCases := map[string]struct {
		expectedFrequency   int
		expectedSyllableKey string
	}{
		"three of them should be vowel syllables (V)":                     {3, "V"},
		"three of them should be consonant-vowel syllables (CV)":          {3, "CV"},
		"two of them should be vowel-consonant syllables (VC)":            {2, "VC"},
		"two of them should be consonant-vowel-consonant syllables (CVC)": {2, "CVC"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given
			syllableChancesGenerator := chanceGeneratorThatWillGenerate(t, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9)

			// When
			templates := make([]syllable.TemplateDefinition, 10)
			for i := range 10 {
				templates[i] = syllable.GenerateTemplate(1, syllable.WithSyllableChanceGenerator(syllableChancesGenerator))
			}

			// Then
			assert.Equal(t, tc.expectedFrequency, countSyllableKey(templates, tc.expectedSyllableKey))
		})
	}
}

func TestGenerateTemplate_the_syllable_ending_with_constant(t *testing.T) {
	const (
		vcChanceOver10  = 6
		cvcChanceOver10 = 8
	)
	testCases := map[string]struct {
		firstSyllableChance   int
		forbiddenSyllableKeys []string
	}{
		"VC cannot be followed by a a syllable starting with a consonant":  {vcChanceOver10, []string{"CV", "CVC"}},
		"CVC cannot be followed by a a syllable starting with a consonant": {cvcChanceOver10, []string{"CV", "CVC"}},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			for secondSyllableChanceOver10 := range 10 {
				t.Run(fmt.Sprintf("when chance for the second syllable is %d/10", secondSyllableChanceOver10+1), func(t *testing.T) {
					var chancesGenerator syllable.GenerateRandomIntegerUpToFn
					// Given
					chancesGenerator = chanceGeneratorThatWillGenerate(t, tc.firstSyllableChance, secondSyllableChanceOver10)

					// When
					template := syllable.GenerateTemplate(2, syllable.WithSyllableChanceGenerator(chancesGenerator))

					// Then
					require.NotNil(t, template)
					require.Len(t, template.SyllableKeySequence(), 2)
					assert.NotContains(t, tc.forbiddenSyllableKeys, template.SyllableKeySequence()[1], "the second syllable should not be %s but is %s", tc.forbiddenSyllableKeys, template.SyllableKeySequence()[1])
				})
			}
		})
	}
}

func TestGenerateTemplate_when_two_consecutive_vowels_are_generated(t *testing.T) {
	const (
		vowelChance = 0

		swapAChance        = 0
		swapEChance        = 1
		swapIChance        = 2
		swapOChance        = 3
		swapUChance        = 4
		swapAReverseChance = 5
		swapEReverseChance = 6
		swapIReverseChance = 7
		swapOReverseChance = 8
		swapUReverseChance = 9
	)
	testCases := map[string]struct {
		vowelTemplateChance          int
		singleVowelTemplate          string
		expectedWeightFirstSyllable  int
		expectedWeightSecondSyllable int
	}{
		"and first template doesn't contain single vowel A the second template does": {swapAChance, "(a)", 0, 10},
		"and first template doesn't contain single vowel E the second template does": {swapEChance, "(e)", 0, 6},
		"and first template doesn't contain single vowel I the second template does": {swapIChance, "(i)", 0, 4},
		"and first template doesn't contain single vowel O the second template does": {swapOChance, "(o)", 0, 2},
		"and first template doesn't contain single vowel U the second template does": {swapUChance, "(u)", 0, 1},
		"and first template contains single vowel A the second template doesn't":     {swapAReverseChance, "(a)", 10, 0},
		"and first template contains single vowel E the second template doesn't":     {swapEReverseChance, "(e)", 6, 0},
		"and first template contains single vowel I the second template doesn't":     {swapIReverseChance, "(i)", 4, 0},
		"and first template contains single vowel O the second template doesn't":     {swapOReverseChance, "(o)", 2, 0},
		"and first template contains single vowel U the second template doesn't":     {swapUReverseChance, "(u)", 1, 0},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var syllableChancesGenerator syllable.GenerateRandomIntegerUpToFn
			// Given
			syllableChancesGenerator = chanceGeneratorThatWillGenerate(t, vowelChance, vowelChance)
			vowelTemplateChanceGenerator := chanceGeneratorThatWillGenerate(t, tc.vowelTemplateChance)

			// When
			template := syllable.GenerateTemplate(2,
				syllable.WithSyllableChanceGenerator(syllableChancesGenerator),
				syllable.WithVowelTemplateChanceGenerator(vowelTemplateChanceGenerator),
			)

			// Then
			require.NotNil(t, template)
			require.Equal(t, []string{"V", "V"}, template.SyllableKeySequence(), "both syllables should be vowel syllable templates")

			singleVowelOccurrencesOnFirstSyllableTemplate := strings.Count(template.TemplateSequence()[0], tc.singleVowelTemplate)
			assert.Equal(t, tc.expectedWeightFirstSyllable, singleVowelOccurrencesOnFirstSyllableTemplate)
			singleVowelOccurrencesOnSecondSyllableTemplate := strings.Count(template.TemplateSequence()[1], tc.singleVowelTemplate)
			assert.Equal(t, tc.expectedWeightSecondSyllable, singleVowelOccurrencesOnSecondSyllableTemplate)
		})
	}

}

func countSyllableKey(templates []syllable.TemplateDefinition, key string) int {
	var count int
	for _, template := range templates {
		if len(template.SyllableKeySequence()) == 0 {
			continue
		}
		if template.SyllableKeySequence()[0] == key {
			count++
		}
	}
	return count
}

func chanceGeneratorThatWillGenerate(t *testing.T, sequence ...int) syllable.GenerateRandomIntegerUpToFn {
	var i int
	return func(n int) int {
		require.GreaterOrEqual(t, len(sequence), i, "unexpected call to chance generator for the %d time", i+1)
		chance := sequence[i]
		i++
		return chance
	}
}
