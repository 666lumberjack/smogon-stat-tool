package main

import (
	"fmt"
)

func main() {
	// Get command line arguments and flags
	mode, args, flags, err := getArguments()
	if err != nil {
		fmt.Printf("Error parsing arguments and flags: %e", err)
		return
	}

	// Run the specified mode
	switch mode {
	case "moves":
		showMoveData(args, flags)
	default:
		fmt.Printf("Mode string did not map to valid mode")
	}
}
