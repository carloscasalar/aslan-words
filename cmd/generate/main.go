// CLI util to generate aslan words
package main

import (
	"context"
	"fmt"
	"github.com/carloscasalar/aslan-words"
	"log"
	"math/rand"
)

func main() {
	numberOfSyllables := rand.Intn(5) + 2

	ctx := context.Background()

	log.Println(fmt.Sprintf("%d syllables:", numberOfSyllables))
	log.Println(aslanwords.MustGenerate(ctx, aslanwords.WithNumberOfSyllables(numberOfSyllables)))
	log.Println(aslanwords.MustGenerate(ctx))
	log.Println(aslanwords.MustGenerate(ctx, aslanwords.WithNumberOfSyllablesBetween(1, numberOfSyllables)))
}
