package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

type toolFlags struct {
	mode       string
	pokemon    string
	tier       string
	generation string
	statURL    string
	weighting  int64
}

// getFlags retrieves flags and arguments to use
func getFlags() (*toolFlags, error) {
	flags := &toolFlags{}

	if len(os.Args) < 2 { // first arg is always filepath of executable
		return nil, errors.New("no arguments specified")
	}

	// Define flags
	flags.mode = *flag.String("mode", "moves", "The mode to run Smogon Stat Tool in")
	flags.pokemon = *flag.String("pokemon", "pikachu", "The Pokemon to get stats for")
	flags.tier = *flag.String("tier", overused, "The tier to get stats from")
	flags.generation = *flag.String("generation", "gen8", "The generation to get stats from")
	flags.statURL = *flag.String("url", "", "An override URL to pull stats from")
	flags.weighting = int64(*flag.Int("weighting", skilled, "The player skill level to draw stats from"))

	flag.Parse()

	// If no flags are defined, default to quick move lookup mode using arguments
	if flag.NFlag() == 0 {
		flags.pokemon = flag.Arg(0)
		flags.tier = flag.Arg(1)
		weighting, err := strconv.ParseInt(flag.Arg(2), 10, 0)
		if err != nil {
			return &toolFlags{}, err
		}
		flags.weighting = weighting
	}

	return flags, nil
}

// format returns generation + tier
func (f *toolFlags) format() string {
	return f.generation + f.tier
}
