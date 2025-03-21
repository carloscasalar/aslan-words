package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type commandOptions struct {
	NumberOfSyllables int `short:"s" default:"2" long:"number-of-syllables" description:"Number of syllables of the aslan word to generate-word"`
}

func readOptionsOrFail() commandOptions {
	var opts commandOptions
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(flags.ErrorType); ok && flagsErr == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}
	return opts
}
