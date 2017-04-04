package main

import (
	"fmt"
	"os"
)

func main() {

	//if the caller didn't provide a URL to fetch...
	if len(os.Args) < 2 {
		//print the usage and exit with an error
		fmt.Printf("usage:\n  pagetitle <url>\n")
		os.Exit(1)
	}

	URL := os.Args[1]

	fmt.Printf(URL + "\n") // TODO delete ******
}
