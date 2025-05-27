// Description: Example of how to use the aslanwords package to generate Aslan words.
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
	maleNameWithType, err := aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeMaleName))
	if err != nil {
		fmt.Printf("Error generating male name using Generate/WithType: %v\n", err)
	} else {
		fmt.Printf("Male name using Generate/WithType (default syllables): %s\n", maleNameWithType)
	}

	femaleNameWithType3Syllables, err := aslanwords.Generate(ctx, aslanwords.WithType(aslanwords.TypeFemaleName), aslanwords.WithNumberOfSyllables(3))
	if err != nil {
		fmt.Printf("Error generating female name (3 syllables) using Generate/WithType: %v\n", err)
	} else {
		fmt.Printf("Female name with 3 syllables using Generate/WithType: %s\n", femaleNameWithType3Syllables)
	}

	fmt.Println("\n--- Examples of deprecated functions (for awareness, prefer GenerateName) ---")
	// Note: GenerateMaleName and GenerateFemaleName are deprecated.
	// These examples show they still work but will internally call GenerateName.
	deprecatedMaleName, err := aslanwords.GenerateMaleName(ctx)
	if err == nil {
		fmt.Printf("Deprecated GenerateMaleName example: %s\n", deprecatedMaleName)
	}
	deprecatedFemaleName, err := aslanwords.GenerateFemaleName(ctx, aslanwords.WithNumberOfSyllables(3))
	if err == nil {
		fmt.Printf("Deprecated GenerateFemaleName (3 syllables) example: %s\n", deprecatedFemaleName)
	}
}
