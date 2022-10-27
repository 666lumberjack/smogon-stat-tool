package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

// showMoveData gets moveset data for specified pokemon, generation, tier and skill weighting
func showMoveData(flags *toolFlags) error {
	pokemon := strings.ToLower(flags.pokemon)

	// Construct a URL to get the stat txt file from
	statsURL, err := constructStatURL(flags)
	if err != nil {
		return err
	}
	// Get the raw stat data in txt format
	resp, err := http.Get(statsURL)
	if err != nil {
		return err
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

	// Print move data to the terminal
	for scanner.Scan() {
		// If we see multiple consecutive dashes, we've reached the end of the moves section
		if strings.Contains(scanner.Text(), "--") {
			break
		}
		fmt.Println(scanner.Text())
	}

	return nil
}
