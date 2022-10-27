package main

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

func weightingTierName(weighting int64) string {
	switch weighting {
	case unskilled:
		return "unskilled"
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
const (
	ubers      = "ubers"
	overused   = "ou"
	underused  = "uu"
	rarelyused = "ru"
	neverused  = "nu"
	pu         = "pu" // PU and ZU don't stand for anything
	zu         = "zu"
)
