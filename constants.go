package main

const statBaseURL = "http://smogon.com/stats"

const currentGen = "gen8"

// Stats are typically available weighted by four different player skill ratings for each tier
// none represents no weighting by skill at all
// average represents stats weighted to approximate what a player of average skill sees
// skilled represents stats weighted to approximate what a player striving to seriously compete sees
// elite represents stats weighted to approximate what what an excellent player sees
const (
	none = iota
	average
	skilled
	elite
)

// weightingTierName returns the string name corresponding to a 'skill tier' stats could be weighted by
func weightingTierName(weighting int64) string {
	switch weighting {
	case none:
		return "none"
	case average:
		return "average"
	case skilled:
		return "skilled"
	case elite:
		return "elite"
	default:
		return "unrecognised"
	}
}

// Tier names and their shorthand abbreviations as used in filenames
// TODO: Add more formats and tiers
const (
	ubers      = "ubers"
	overused   = "ou"
	underused  = "uu"
	rarelyused = "ru"
	neverused  = "nu"
	pu         = "pu" // PU and ZU don't stand for anything
	zu         = "zu"
)

// We represent runtime flags with a slice of int64s
// If the flags slice contains the integer representing a given flag, it was provided
const (
	forcePathGen = iota
)

// We define a slice of multiple valid strings that map to each mode
var movesStrings = []string{"m", "mv", "move", "moves"}
