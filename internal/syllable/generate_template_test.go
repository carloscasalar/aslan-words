package syllable_test

import (
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
			template := syllable.GenerateTemplate(1, syllable.WithIntegerGenerator(chancesGenerator))

			// Then
			require.NotNil(t, template)
			require.Len(t, template.SyllableKeySequence(), 1, "only one syllable is expected to be generated")
			assert.Equal(t, tc.expectedSyllableKey, template.SyllableKeySequence()[0], "out of ten times, when chance is %d, the expected syllable key is %s", tc.chance, tc.expectedSyllableKey)
		})
	}
}
