package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// showMoveData gets moveset data for specified pokemon, generation, tier and skill weighting
func showMoveData(args []string, flags []int64) {
	// Check and log provided arguments for user
	if len(args) < 3 {
		fmt.Printf("Insufficient number of arguments provided: expected pokemon, tier, weighting but got only %d arguments", len(args))
	}

	pokemon, tier, weighting, generation, urlOverride := parseMovesetArguments(args)

	// weighting can either be an interger value corresponding to a specific skill level,
	// or a numerical rating value to find the closest match for - we log each differently
	if weighting > 4 {
		fmt.Printf("Common moves for %s in %s%s near specified weighting %d:\n", pokemon, generation, tier, weighting)
	} else {
		fmt.Printf("Common moves for %s in %s%s at give skill weighting %s:\n", pokemon, generation, tier, weightingTierName(weighting))
	}

	movesPathSpec := &pathSpec{
		pokemon:     pokemon,
		tier:        tier,
		weighting:   weighting,
		generation:  generation,
		overrideURL: urlOverride,
	}

	// Construct a URL to get the stat txt file from
	statsURL, err := constructStatURL("moves", movesPathSpec)
	if err != nil {
		fmt.Printf("Error constructing stat URL: %e", err)
		return
	}
	// Get the raw stat data in txt format
	resp, err := http.Get(statsURL)
	if err != nil {
		fmt.Printf("Error getting stats file: %e", err)
		return
	}
	defer resp.Body.Close()

	// Parse out the move data we're interested in
	scanner := bufio.NewScanner(resp.Body)
	foundPokemon := false

	// Scan until we find the first line of the moves section for the specified pokemon
	for scanner.Scan() {
		text := strings.ToLower(scanner.Text())
		if foundPokemon {
			if strings.Contains(text, "moves") {
				break
			}
			continue
		}
		// Pokemon names may be listed in teammates section, so we check for no percent sign to ensure
		// we've found the data entry for the pokemon itself
		if strings.Contains(text, pokemon) && !strings.Contains(scanner.Text(), "%") {
			foundPokemon = true
		}
	}

	if !foundPokemon {
		fmt.Printf("Could not find stats for those parameters. Was %s used at least one time in that tier and skill bracket?\n", pokemon)
	}

	// Print move data to the terminal
	for scanner.Scan() {
		// If we see multiple consecutive dashes, we've reached the end of the moves section
		if strings.Contains(scanner.Text(), "--") {
			break
		}
		fmt.Println(scanner.Text())
	}
}

// parseMovesetArguments takes a slice of string arguments and extracts the data we need to specify
// which moveset stats to get. It expects the args slice is already validated as at least 3 entries long
func parseMovesetArguments(args []string) (string, string, int64, string, string) {
	// In moves mode, we expect the first three arguments to be pokemon, tier, and weighting in order
	// an optional fourth argument specifies generation; if not present, we default to currentGen
	// a fourth or fifth argument can specify an override URL to use instead of generating one
	pokemon := strings.ToLower(args[0])
	tier := strings.ToLower(args[1])
	weighting, _ := strconv.ParseInt(args[2], 10, 64)
	generation := currentGen
	urlOverride := ""
	// Fourth argument could be generation or override url
	if len(args) > 3 {
		switch {
		case len(args[3]) < 6:
			generation = strings.ToLower(args[3])
		default:
			urlOverride = args[3]
		}
	}
	// If there is a fifth arg, it's always override url
	if len(args) > 4 {
		urlOverride = args[4]
	}
	return pokemon, tier, weighting, generation, urlOverride
}
