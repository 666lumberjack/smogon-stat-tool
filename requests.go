package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

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

// constructStatURL takes a mode string and a flags struct and returns a string containing the URL
// of the stat file we want to pull from
func constructStatURL(mode string, spec *pathSpec) (string, error) {
	if spec.overrideURL != "" {
		fmt.Printf("Using override URL: %s\n", spec.overrideURL)
		return spec.overrideURL, nil
	}

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
	closestIndex := 0
	for i, rating := range ratings {
		if math.Abs(float64(rating-spec.weighting)) < float64(closestIndex) {
			closestIndex = i
		}
	}

	return files[closestIndex], nil
}
