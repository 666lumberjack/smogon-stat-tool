package main

import (
	"fmt"
	"os"
)

func main() {
	// Get command line arguments
	if len(os.Args) < 2 { // first arg is always filepath of executable
		fmt.Println("You must specify at least one argument.")
		return
	}
	// Log provided arguments for user

	// Run the specified mode
}
