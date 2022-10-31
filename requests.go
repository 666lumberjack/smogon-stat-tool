package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// pathSpec contains everything necessary to build a path to the file with stats we want
type pathSpec struct {
	pokemon     string
	tier        string
	generation  string
	overrideURL string
	weighting   int64
}

// format returns generation + tier
func (f *pathSpec) format() string {
	return f.generation + f.tier
}

// modeFolderPath returns the path to the subfolder with the stats we want for a given mode
func modeFolderPath(mode string) string {
	switch mode {
	case "moves":
		return "moveset/"
	default:
		return ""
	}
}

// guessStatURL tries to guess the location of a stat file based on preknowledge of the filepaths
func guessStatURL(mode string, spec *pathSpec) string {
	if spec.overrideURL != "" {
		fmt.Printf("Using override URL: %s\n", spec.overrideURL)
		return spec.overrideURL
	}

	dateGuess := guessStatDate(time.Now().UTC())

	modeFolder := modeFolderPath(mode)

	filenameGuess := guessFilename(spec)

	return statBaseURL + "/" + dateGuess + modeFolder + filenameGuess
}

// guessStatDate takes a time and returns a string representation of the estimated month that the latest stats are available for
// new stats are uploaded on the 1st, 2nd or 3rd of the month, but usually late on the 1st
func guessStatDate(statTime time.Time) string {
	pathTimeFormat := "2006-01/"

	day := statTime.Day()

	if day < 2 { // if first of the month, jump back 48 hours to last month
		statTime = statTime.Add(-time.Hour * 48)
	}

	// We always jump back one month and a day to account for stats for one months being released the next
	statTime = statTime.AddDate(0, -1, -1)

	return statTime.Format(pathTimeFormat)
}

// guessFilename takes a pathspec struct pointer and returns a string containing the estimated name of the file with stats we want
func guessFilename(spec *pathSpec) string {
	ratingsForWeightings := []int64{0, 1500, 1630, 1760}

	if spec.tier == overused {
		// The overused tier has many more players that all others, and the skill threshold for the upper two weightings is thus higher
		ratingsForWeightings[2], ratingsForWeightings[3] = 1695, 1825
	}

	ratingString := ""
	if spec.weighting < 4 {
		ratingString = fmt.Sprintf("%d", ratingsForWeightings[spec.weighting])
	} else {
		closestWeighting := closestIndex(ratingsForWeightings, spec.weighting)
		ratingString = fmt.Sprintf("%d", ratingsForWeightings[closestWeighting])
	}

	return spec.format() + "-" + ratingString + ".txt"
}

// constructStatURL takes a mode string and a flags struct and returns a string containing the URL
// of the stat file we want to pull from
func constructStatURL(mode string, spec *pathSpec) (string, error) {

	// Parse base stats page to establish available dates
	dateString, err := getLatestStatDate(spec)
	if err != nil {
		return "", err
	}

	statDateURL := statBaseURL + "/" + dateString

	statFolderURL := statDateURL + modeFolderPath(mode)

	return finalStatPath(statFolderURL, spec)
}

// getLatestStatDate makes a get request to the base stat page and parses
// it to find the most recent date with stats available
func getLatestStatDate(spec *pathSpec) (string, error) {
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
func finalStatPath(url string, spec *pathSpec) (string, error) {
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
			if strings.Contains(t.Data, spec.format()) {
				statFiles = append(statFiles, t.Data)
			}
		}
	}

	filename, err := fileForWeighting(statFiles, spec)
	if err != nil {
		return "", err
	}

	return url + filename, nil
}

// fileForWeighting takes a slice of filename strings and a pointer
// to a toolFlags struct and returns either the string at index matching
// the weighting value if it is <4 or the file most closely matching the
// provided skill rating if weighting >=4
func fileForWeighting(files []string, spec *pathSpec) (string, error) {
	if spec.weighting < 4 {
		return files[spec.weighting], nil
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
	fileIndex := closestIndex(ratings, spec.weighting)

	return files[fileIndex], nil
}
