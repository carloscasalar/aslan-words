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

func TestGenerateName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		gender        string
		opts          []aslanwords.GeneratorOption
		wantErr       bool
		errorContains string
	}{
		{
			name:    "male name generation default syllables",
			gender:  "male",
			opts:    nil,
			wantErr: false,
		},
		{
			name:   "male name generation with 3 syllables",
			gender: "male",
			opts: []aslanwords.GeneratorOption{
				aslanwords.WithNumberOfSyllables(3),
			},
			wantErr: false,
		},
		{
			name:   "male name generation with 1 syllable",
			gender: "male",
			opts: []aslanwords.GeneratorOption{
				aslanwords.WithNumberOfSyllables(1),
			},
			wantErr: false,
		},
		{
			name:    "female name generation default syllables",
			gender:  "female",
			opts:    nil,
			wantErr: false,
		},
		{
			name:   "female name generation with 2 syllables",
			gender: "female",
			opts: []aslanwords.GeneratorOption{
				aslanwords.WithNumberOfSyllables(2),
			},
			wantErr: false,
		},
		{
			name:   "female name generation with 1 syllable",
			gender: "female",
			opts: []aslanwords.GeneratorOption{
				aslanwords.WithNumberOfSyllables(1),
			},
			wantErr: false,
		},
		{
			name:          "invalid gender unknown",
			gender:        "unknown",
			opts:          nil,
			wantErr:       true,
			errorContains: "invalid gender: unknown",
		},
		{
			name:          "invalid gender empty",
			gender:        "",
			opts:          nil,
			wantErr:       true,
			errorContains: "invalid gender: ",
		},
		{
			name:   "male name generation with 0 syllables (error case)",
			gender: "male",
			opts: []aslanwords.GeneratorOption{
				aslanwords.WithNumberOfSyllables(0),
			},
			wantErr:       true,
			errorContains: "number of syllables must be one or greater",
		},
		{
			name:   "female name generation with 0 syllables (error case)",
			gender: "female",
			opts: []aslanwords.GeneratorOption{
				aslanwords.WithNumberOfSyllables(0),
			},
			wantErr:       true,
			errorContains: "number of syllables must be one or greater",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := aslanwords.GenerateName(ctx, tt.gender, tt.opts...)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, name)
				// Verifying exact syllable counts is complex and depends on the fantasyname library.
				// For now, we ensure a name is generated.
			}
		})
	}
}

func TestGenerate_when_invalid_range_of_syllables_generate_should_return_error(t *testing.T) {
	ctx := context.Background()
	_, err := aslanwords.Generate(ctx, aslanwords.WithNumberOfSyllablesBetween(5, 3))
	assert.Error(t, err)
}
