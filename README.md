# Aslan Words [![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/carloscasalar/aslan-words)

This is a library for generating Aslan words according to the guidelines of the awesome [Pirates of Drinax](https://www.mongoosepublishing.com/products/the-pirates-of-drinax?srsltid=AfmBOoq1tbdk_O_QGAK5xoYQc13tXsNSOn8tu5wHto5TOTVtvqMkq-pH) campaign, in the Traveller RPG universe.

## Installation

To install the library into your go project run:

```sh
go get github.com/carloscasalar/aslan-words
```

## Usage

### Library

Here is an example of how to use the library to generate Aslan names:

```go
package main

import (
	"context"
	"fmt"

	"github.com/carloscasalar/aslan-words/pkg/aslanwords"
)

func main() {
	ctx := context.Background()
	word, err := aslanwords.Generate(ctx, aslanwords.WithNumberOfSyllables(3))
	if err != nil {
		fmt.Printf("Error generating Aslan word (3 syllables): %v\n", err)
		return
	}
	fmt.Printf("Generic Aslan word with 3 syllables: %s\n", word)

	word = aslanwords.MustGenerate(ctx, aslanwords.WithNumberOfSyllablesBetween(3, 6))
	fmt.Printf("Generic Aslan word with between 3 and 6 syllables: %s\n", word)

	fmt.Println("\n--- Name Generation Examples (using new GenerateName function) ---")

	// Using the new GenerateName function (recommended)
	maleName, err := aslanwords.GenerateName(ctx, aslanwords.Male)
	if err != nil {
		fmt.Printf("Error generating male name with GenerateName: %v\n", err)
	} else {
		fmt.Printf("Default male name (using GenerateName): %s (typically 3-5 syllables)\n", maleName)
	}

	femaleName, err := aslanwords.GenerateName(ctx, aslanwords.Female)
	if err != nil {
		fmt.Printf("Error generating female name with GenerateName: %v\n", err)
	} else {
		fmt.Printf("Default female name (using GenerateName): %s (typically 2-4 syllables)\n", femaleName)
	}

	maleName2Syllables, err := aslanwords.GenerateName(ctx, aslanwords.Male, aslanwords.WithNumberOfSyllables(2))
	if err != nil {
		fmt.Printf("Error generating male name (2 syllables) with GenerateName: %v\n", err)
	} else {
		fmt.Printf("Male name with 2 syllables (using GenerateName): %s\n", maleName2Syllables)
	}

	femaleName4Syllables, err := aslanwords.GenerateName(ctx, aslanwords.Female, aslanwords.WithNumberOfSyllables(4))
	if err != nil {
		fmt.Printf("Error generating female name (4 syllables) with GenerateName: %v\n", err)
	} else {
		fmt.Printf("Female name with 4 syllables (using GenerateName): %s\n", femaleName4Syllables)
	}
	
	fmt.Println("\n--- Name Generation using Generate with WithType (still valid) ---")
	// Using Generate with WithType for more control (still valid)
	maleNameWithType, err := aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeMaleName))
	if err != nil {
		fmt.Printf("Error generating male name using Generate/WithType: %v\n", err)
	} else {
		fmt.Printf("Male name using Generate/WithType (default syllables): %s\n", maleNameWithType)
	}

	femaleNameWithType3Syllables, err := aslanwords.Generate(ctx, 
		aslanwords.WithType(aslanwords.TypeFemaleName),
		aslanwords.WithNumberOfSyllables(3),
	)
	if err != nil {
		fmt.Printf("Error generating female name (3 syllables) using Generate/WithType: %v\n", err)
	} else {
		fmt.Printf("Female name with 3 syllables using Generate/WithType: %s\n", femaleNameWithType3Syllables)
	}

	fmt.Println("\n--- Examples of deprecated functions (for awareness, prefer GenerateName) ---")
	// Note: GenerateMaleName and GenerateFemaleName are deprecated.
	// These examples show they still work but will internally call GenerateName.
	deprecatedMaleName, err := aslanwords.GenerateMaleName(ctx) // Deprecated
	if err == nil {
		fmt.Printf("Deprecated GenerateMaleName example: %s\n", deprecatedMaleName)
	}
	deprecatedFemaleName, err := aslanwords.GenerateFemaleName(ctx, aslanwords.WithNumberOfSyllables(3)) // Deprecated
	if err == nil {
		fmt.Printf("Deprecated GenerateFemaleName (3 syllables) example: %s\n", deprecatedFemaleName)
	}
}
```

### Generating Male and Female Names

The primary way to generate gendered Aslan names is using the `GenerateName` function:

-   `GenerateName(ctx context.Context, gender aslanwords.Gender, opts ...aslanwords.GeneratorOption) (string, error)`:
    -   Takes a `gender` argument (`aslanwords.Male` or `aslanwords.Female`).
    -   By default, male names are 3-5 syllables long, and female names are 2-4 syllables long and tend to end in a vowel sound.
    -   Accepts `GeneratorOption` arguments like `aslanwords.WithNumberOfSyllables` to override default syllable counts.

**Deprecated Functions:**
-   `GenerateMaleName(ctx, opts...)` and `GenerateFemaleName(ctx, opts...)` are now **deprecated**. They internally call the new `GenerateName` function. Users are encouraged to migrate to `GenerateName` for clarity and future enhancements.

**Alternative using `Generate`:**
-   You can also use the main `Generate` function with the `aslanwords.WithType` option:
    -   `aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeMaleName), ...)`
    -   `aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeFemaleName), ...)`
    -   This approach now routes to the same underlying logic as `GenerateName`.

All methods for generating names accept the same `GeneratorOption` arguments (e.g., `aslanwords.WithNumberOfSyllables`, `aslanwords.WithNumberOfSyllablesBetween`). If a syllable option is provided, it will override the default syllable range for that name type.

## Testing

To run the tests, use the following command:

```sh
make test
```

## CLI

The project includes a basic CLI to generate Aslan names. You can build and run the CLI using the following commands:

```sh
# Build the CLI
make build

# Run the CLI
./out/generate-word --help
```

The CLI supports the following options:

| Flag                      | Description                                                                                                | Default |
| ------------------------- | ---------------------------------------------------------------------------------------------------------- | ------- |
| `-s, --syllables <n>`     | Number of syllables of the word to generate. For 'male' or 'female' types, this overrides their defaults. | 2       |
| `--type <type>`           | Type of word to generate ("word", "male", "female").                                                       | "word"  |

**CLI Examples:**

Generate a generic word (default behavior):
```sh
./out/generate-word -s 3
# Output: (A random 3-syllable Aslan word)
```

Generate a male name with default syllables (3-5):
```sh
./out/generate-word --type male
# Output: (A random 3-5 syllable male Aslan name)
```

Generate a female name with default syllables (2-4):
```sh
./out/generate-word --type female
# Output: (A random 2-4 syllable female Aslan name)
```

Generate a male name with a specific number of syllables:
```sh
./out/generate-word --type male -s 4
# Output: (A random 4-syllable male Aslan name)
```

Generate a female name with a specific number of syllables:
```sh
./out/generate-word --type female -s 2
# Output: (A random 2-syllable female Aslan name)
```

**Note on Syllable Counts for Names:** When using `--type male` or `--type female`, the default syllable ranges (3-5 for male, 2-4 for female) are used. If the `-s` (or `--syllables`) flag is also provided, it will override these defaults.

![cli demo](demo/demo.gif)
