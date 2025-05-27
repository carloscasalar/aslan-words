package aslanwords_test

import (
	"context"
	"testing"

	"github.com/carloscasalar/aslan-words/pkg/aslanwords"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_when_no_options_it_should_generate_a_word_with_default_number_of_syllables(t *testing.T) {
	ctx := context.Background()
	word, err := aslanwords.Generate(ctx, nil...)
	require.NoError(t, err)
	assert.NotEmpty(t, word)
}

func TestGenerate_when_called_with_fixed_number_of_syllables_generate_should_return_a_word(t *testing.T) {
	ctx := context.Background()
	word, err := aslanwords.Generate(ctx, aslanwords.WithNumberOfSyllables(3))
	require.NoError(t, err)
	assert.NotEmpty(t, word)
}

func TestGenerate_when_random_number_of_syllables_between_2_and_4_generate_should_not_return_error(t *testing.T) {
	ctx := context.Background()
	word, err := aslanwords.Generate(ctx, aslanwords.WithNumberOfSyllablesBetween(2, 4))
	require.NoError(t, err)
	assert.NotEmpty(t, word)
}

func TestGenerate_when_asked_to_generate_word_with_less_than_one_syllables_should_return_error(t *testing.T) {
	ctx := context.Background()
	_, err := aslanwords.Generate(ctx, aslanwords.WithNumberOfSyllables(0))
	assert.Error(t, err)
}

// Test for deprecated GenerateMaleName, which now calls GenerateName.
func TestGenerateMaleName_generates_non_empty_string(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.GenerateMaleName(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, name, "GenerateMaleName should produce a non-empty string")
}

// Test for deprecated GenerateMaleName, which now calls GenerateName.
func TestGenerateMaleName_with_explicit_syllables(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.GenerateMaleName(ctx, aslanwords.WithNumberOfSyllables(2))
	require.NoError(t, err)
	assert.NotEmpty(t, name, "GenerateMaleName with 2 syllables should produce a non-empty string")
}

// Test for deprecated GenerateFemaleName, which now calls GenerateName.
func TestGenerateFemaleName_generates_non_empty_string(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.GenerateFemaleName(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, name, "GenerateFemaleName should produce a non-empty string")
}

// Test for deprecated GenerateFemaleName, which now calls GenerateName.
func TestGenerateFemaleName_with_explicit_syllables(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.GenerateFemaleName(ctx, aslanwords.WithNumberOfSyllables(4))
	require.NoError(t, err)
	assert.NotEmpty(t, name, "GenerateFemaleName with 4 syllables should produce a non-empty string")
}

func TestGenerateName_Male_generates_non_empty_string(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.GenerateName(ctx, aslanwords.Male)
	require.NoError(t, err)
	assert.NotEmpty(t, name, "GenerateName with Male should produce a non-empty string")
}

func TestGenerateName_Female_generates_non_empty_string(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.GenerateName(ctx, aslanwords.Female)
	require.NoError(t, err)
	assert.NotEmpty(t, name, "GenerateName with Female should produce a non-empty string")
}

func TestGenerateName_Male_with_explicit_syllables(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.GenerateName(ctx, aslanwords.Male, aslanwords.WithNumberOfSyllables(2))
	require.NoError(t, err)
	assert.NotEmpty(t, name, "GenerateName with Male and 2 syllables should produce a non-empty string")
}

func TestGenerateName_Female_with_explicit_syllables(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.GenerateName(ctx, aslanwords.Female, aslanwords.WithNumberOfSyllables(5))
	require.NoError(t, err)
	assert.NotEmpty(t, name, "GenerateName with Female and 5 syllables should produce a non-empty string")
}

func TestGenerateName_invalid_options_combination(t *testing.T) {
	ctx := context.Background()
	_, err := aslanwords.GenerateName(ctx, aslanwords.Male, aslanwords.WithNumberOfSyllables(0))
	assert.Error(t, err, "GenerateName with Male and 0 syllables should return an error")

	_, err = aslanwords.GenerateName(ctx, aslanwords.Female, aslanwords.WithNumberOfSyllables(0))
	assert.Error(t, err, "GenerateName with Female and 0 syllables should return an error")
}

func TestGenerate_with_TypeMaleName_option(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeMaleName))
	require.NoError(t, err)
	assert.NotEmpty(t, name, "Generate with TypeMaleName should produce a non-empty string")
}

func TestGenerate_with_TypeFemaleName_option(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeFemaleName))
	require.NoError(t, err)
	assert.NotEmpty(t, name, "Generate with TypeFemaleName should produce a non-empty string")
}

func TestGenerate_with_TypeWord_option(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeWord))
	require.NoError(t, err)
	assert.NotEmpty(t, name, "Generate with TypeWord should produce a non-empty string")
}

func TestGenerate_with_TypeMaleName_and_syllables_option(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeMaleName), aslanwords.WithNumberOfSyllables(2))
	require.NoError(t, err)
	assert.NotEmpty(t, name, "Generate with TypeMaleName and 2 syllables should produce a non-empty string")
}

func TestGenerate_with_TypeFemaleName_and_syllables_option(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeFemaleName), aslanwords.WithNumberOfSyllables(5))
	require.NoError(t, err)
	assert.NotEmpty(t, name, "Generate with TypeFemaleName and 5 syllables should produce a non-empty string")
}

func TestGenerate_with_TypeWord_and_syllables_option(t *testing.T) {
	ctx := context.Background()
	name, err := aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeWord), aslanwords.WithNumberOfSyllables(3))
	require.NoError(t, err)
	assert.NotEmpty(t, name, "Generate with TypeWord and 3 syllables should produce a non-empty string")
}

// Test to ensure that providing an invalid combination (though hard with current options) or future validation logic is caught.
// For now, an invalid WordType is not possible at compile time if using the enum.
// This test serves as a placeholder if more complex validation rules are added to option combinations.
func TestGenerate_invalid_options_combination(t *testing.T) {
	ctx := context.Background()
		// Test that providing a syllable count that itself is invalid, along with a type, still fails.
		// These tests now route through Generate -> GenerateName for TypeMaleName/TypeFemaleName
	_, err := aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeMaleName), aslanwords.WithNumberOfSyllables(0))
	assert.Error(t, err, "Generate with TypeMaleName and 0 syllables should return an error")

	_, err = aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeFemaleName), aslanwords.WithNumberOfSyllables(0))
	assert.Error(t, err, "Generate with TypeFemaleName and 0 syllables should return an error")

	_, err = aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeWord), aslanwords.WithNumberOfSyllables(0))
	assert.Error(t, err, "Generate with TypeWord and 0 syllables should return an error")
}

func TestGenerate_when_invalid_range_of_syllables_generate_should_return_error(t *testing.T) {
	ctx := context.Background()
	_, err := aslanwords.Generate(ctx, aslanwords.WithNumberOfSyllablesBetween(5, 3))
	assert.Error(t, err)
}
