package main

import (
	"fmt"
)

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
		err := showMoveData(flags)
		if err != nil {
			fmt.Printf("Error showing moveset statistics: %e", err)
		}
	default:
		fmt.Printf("Mode string not recognised.")
	}
}
