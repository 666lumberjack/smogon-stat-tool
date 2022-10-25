package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const statBaseURL = "http://smogon.com/stats"

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
		err := showMoveData(flags)
		if err != nil {
			fmt.Printf("Error showing moveset statistics: %e", err)
		}
	default:
		fmt.Printf("Mode string not recognised.")
	}
}

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

func modeFolderPath(flags *toolFlags) string {
	switch flags.mode {
	case "moves":
		return "moveset/"
	default:
		return ""
	}
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

// constructStatURL takes a flags struct and returns a string containing the URL
// of the stat file we want to pull from
func constructStatURL(flags *toolFlags) (string, error) {
	if flags.statURL != "" {
		return flags.statURL, nil
	}

	// Parse base stats page to establish available dates
	dateString, err := getLatestStatDate(flags)
	if err != nil {
		return "", err
	}

	statDateURL := statBaseURL + "/" + dateString

	statFolderURL := statDateURL + modeFolderPath(flags)

	return finalStatPath(statFolderURL, flags)
}

// getLatestStatDate makes a get request to the base stat page and parses
// it to find the most recent date with stats available
func getLatestStatDate(flags *toolFlags) (string, error) {
	resp, err := http.Get(statBaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	isAnchor := false
	anchorTexts := []string{}
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return anchorTexts[len(anchorTexts)-1], nil
		case tt == html.StartTagToken:
			t := z.Token()
			isAnchor = t.Data == "a"
		case tt == html.TextToken:
			t := z.Token()
			if isAnchor {
				anchorTexts = append(anchorTexts, t.Data)
			}
			isAnchor = false
		}
	}
}

// finalStatPath takes a string URL for a folder with stat files and a
// toolFlags struct pointer and returns a full path to the stats file we want
func finalStatPath(url string, flags *toolFlags) (string, error) {
	// Get page with stat files
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse page HTML for links matching format
	statFiles := []string{}
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		// Using separate if statements so break affects the loop
		if tt == html.ErrorToken {
			break
		}
		if tt == html.TextToken {
			t := z.Token()
			if strings.Contains(t.Data, flags.format()) {
				statFiles = append(statFiles, t.Data)
			}
		}
	}

	filename, err := fileForWeighting(statFiles, flags)
	if err != nil {
		return "", err
	}

	return url + filename, nil
}

// fileForWeighting takes a slice of filename strings and a pointer
// to a toolFlags struct and returns either the string at index matching
// the weighting value if it is <4 or the file most closely matching the
// provided skill rating if weighting >=4
func fileForWeighting(files []string, flags *toolFlags) (string, error) {
	if flags.weighting < 4 {
		return files[flags.weighting], nil
	}

	// construct slice of ratings each file is weighted with
	ratings := []int64{}
	for i, name := range files {
		rating, err := strconv.ParseInt(name, 0, 64)
		if err != nil {
			return "", err
		}
		ratings[i] = rating
	}

	// find index with rating most closely matching that provided
	closestIndex := 0
	for i, rating := range ratings {
		if math.Abs(float64(rating-flags.weighting)) < float64(closestIndex) {
			closestIndex = i
		}
	}

	return files[closestIndex], nil
}

// format returns generation + tier
func (f *toolFlags) format() string {
	return f.generation + f.tier
}
