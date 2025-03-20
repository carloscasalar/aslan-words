package syllable_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/carloscasalar/aslan-words/v1/internal/syllable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGenerateTemplate_out_of_ten_syllables(t *testing.T) {
	testCases := map[string]struct {
		chance              int
		expectedSyllableKey string
	}{
		"when random chance is 0, syllable should be V":   {0, "V"},
		"when random chance is 1, syllable should be V":   {1, "V"},
		"when random chance is 2, syllable should be V":   {2, "V"},
		"when random chance is 3, syllable should be CV":  {3, "CV"},
		"when random chance is 4, syllable should be CV":  {4, "CV"},
		"when random chance is 5, syllable should be CV":  {5, "CV"},
		"when random chance is 6, syllable should be VC":  {6, "VC"},
		"when random chance is 7, syllable should be VC":  {7, "VC"},
		"when random chance is 8, syllable should be CVC": {8, "CVC"},
		"when random chance is 9, syllable should be CVC": {9, "CVC"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var chancesGenerator syllable.GenerateRandomIntegerUpToFn
			// Given
			chancesGenerator = func(n int) int {
				return tc.chance
			}

			// When
			template := syllable.GenerateTemplate(1, syllable.WithSyllableChanceGenerator(chancesGenerator))

			// Then
			require.NotNil(t, template)
			require.Len(t, template.SyllableKeySequence(), 1, "only one syllable is expected to be generated")
			assert.Equal(t, tc.expectedSyllableKey, template.SyllableKeySequence()[0], "out of ten times, when chance is %d, the expected syllable key is %s", tc.chance, tc.expectedSyllableKey)
		})
	}
}

func TestGenerateTemplate_the_syllable_ending_with_constant(t *testing.T) {
	const (
		vcChance  = 6
		cvcChance = 8
	)
	testCases := map[string]struct {
		firstSyllableChance   int
		forbiddenSyllableKeys []string
	}{
		"VC cannot be followed by a a syllable starting with a consonant":  {vcChance, []string{"CV", "CVC"}},
		"CVC cannot be followed by a a syllable starting with a consonant": {cvcChance, []string{"CV", "CVC"}},
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

func chanceGeneratorThatWillGenerate(t *testing.T, sequence ...int) syllable.GenerateRandomIntegerUpToFn {
	var i int
	return func(n int) int {
		require.GreaterOrEqual(t, len(sequence), i, "unexpected call to chance generator for the %d time", i+1)
		chance := sequence[i]
		i++
		return chance
	}
}
