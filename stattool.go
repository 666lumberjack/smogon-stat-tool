package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

const statBaseURL = "smogon.com/stats"

// Stats are typically available weighted by four different player skill ratings for each tier
// unskilled represents no weighting by skill at all
// average represents stats weighted to approximate what a player of average skill sees
// skilled represents stats weighted to approximate what a player striving to seriously compete sees
// elite represents stats weighted to approximate what what an excellent player sees
const (
	unskilled = iota
	average
	skilled
	elite
)

// Tier names and their shorthand abbreviations as used in filenames
const (
	ubers      = "ubers"
	overused   = "ou"
	underused  = "uu"
	rarelyused = "ru"
	neverused  = "nu"
	pu         = "pu" // PU and ZU don't stand for anything
	zu         = "zu"
)

type toolFlags struct {
	mode       string
	pokemon    string
	tier       string
	generation string
	statURL    string
	weighting  int64
}

func main() {
	// Get command line arguments

	flags, err := getFlags()
	if err != nil {
		fmt.Printf("Error parsing arguments: %e", err)
		return
	}

	// Log provided arguments for user

	fmt.Printf("Running Smogon Stat Tool in %s mode\n", flags.mode)
	fmt.Printf("Using generation and tier %s%s\n", flags.generation, flags.tier)
	fmt.Printf("Getting stats for Pokemon: %s\n", flags.pokemon)

	// weighting can either be an interger value corresponding to a specific skill level,
	// or a numerical rating value to find the closest match for
	if flags.weighting > 4 {
		fmt.Printf("Weighting stats with closest available rating to %d\n", flags.weighting)
	} else {
		fmt.Printf("Weighting stats using predefined value %s\n", weightingTierName(flags.weighting))
	}

	if flags.statURL != "" {
		fmt.Printf("Using override URL: %s\n", flags.statURL)
	}

	// Run the specified mode
	switch flags.mode {
	case "moves":
		showMoveData(flags)
	default:
		fmt.Printf("Mode string not recognised.")
	}
}

// showMoveData gets moveset data for specified pokemon, generation, tier and skill weighting
func showMoveData(flags *toolFlags) {

}

func weightingTierName(weighting int64) string {
	switch weighting {
	case 0:
		return "unskilled"
	case 1:
		return "average"
	case 2:
		return "skilled"
	case 3:
		return "elite"
	default:
		return "unrecognised"
	}
}

// getFlags retrieves flags and arguments to use
func getFlags() (*toolFlags, error) {
	flags := &toolFlags{}

	if len(os.Args) < 2 { // first arg is always filepath of executable
		return nil, errors.New("No arguments specified.")
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
