package main

import (
	"errors"
	"flag"
	"os"
)

// getArguments retrieves command line arguments and flags
// it returns first the mode as a string, then a slice of string arguments,
// then a slice of int64 flags, and finally an error if one occurred
func getArguments() (string, []string, []int64, error) {

	if len(os.Args) < 2 { // first arg is always filepath of executable
		return "", nil, nil, errors.New("no arguments specified")
	}

	// Define and parse flags first to leave us only the other arguments
	flags := getFlags()

	// Get remaining arguments and validate and separate mode
	args := flag.Args()

	mode, otherArgs, err := parseMode(args)
	if err != nil {
		return "", nil, nil, err
	}

	return mode, otherArgs, flags, nil
}

// defineFlags defines all accepted flags against the default CommandLine flagSet
func defineFlags() map[int64]bool {
	flags := make(map[int64]bool)

	flags[forcePathGen] = *flag.Bool("forcePathGen", false, "Force SST to traverse the smogon stats folder with GET requests instead of guessing the path to the file first")

	return flags
}

// getFlags retrieves flags in the form of a slice of int64s
// It expects flags to have already been defined on the default flagSet
func getFlags() []int64 {

	flagsMap := defineFlags()

	flag.Parse()

	var flags []int64
	for flag, present := range flagsMap {
		if present {
			flags = append(flags, flag)
		}
	}

	return flags
}

// parseMode takes a slice of string arguments, extracts and verifies the first item as a valid mode,
// and then returns the mode and remaining arguments
func parseMode(arguments []string) (string, []string, error) {

	if len(arguments) < 1 {
		return "", nil, errors.New("no arguments provided, only flags")
	}

	mode, otherArgs := arguments[0], arguments[1:]

	switch {
	case contains(movesStrings, mode):
		return "moves", otherArgs, nil
	default:
		return "", nil, errors.New("mode not recognised")
	}
}
