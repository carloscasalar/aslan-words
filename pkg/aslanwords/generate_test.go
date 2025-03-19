package aslanwords_test

import (
	"context"
	"testing"

	"github.com/carloscasalar/aslan-words/v1/pkg/aslanwords"
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

func TestGenerate_when_invalid_range_of_syllables_generate_should_return_error(t *testing.T) {
	ctx := context.Background()
	_, err := aslanwords.Generate(ctx, aslanwords.WithNumberOfSyllablesBetween(5, 3))
	assert.Error(t, err)
}
