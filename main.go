package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	//if the caller didn't provide a URL to fetch...
	if len(os.Args) < 2 {
		//print the usage and exit with an error
		fmt.Printf("usage:\n  pagetitle <url>\n")
		os.Exit(1)
	}

	URL := os.Args[1]

	// fmt.Printf(URL + "\n") // TODO delete ******

	// GET the URL
	resp, err := http.Get(URL)

	// if there is an error, report it and exit
	if err != nil {
		log.Fatalf("error fetching URL: %v\n", err)
	}

	fmt.Println(resp) // TODO delete *********

	defer resp.Body.Close() // ensure the response body is closed, defer waits to close until main() returns

	// check response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("response status code was %d\n", resp.StatusCode)
	}

	// check response content type
	ctype := resp.Header.Get("Content-Type")    // extract content type from header
	if !strings.HasPrefix(ctype, "text/html") { // check if ctype does not begin w/ "text/html"
		log.Fatalf("response content type was %s, not text/html\n", ctype)
	}

	fmt.Println("Response Body: ")
	fmt.Println(resp.Body)

	// create a new tokenizer over the response Body
	tokenizer := html.NewTokenizer(resp.Body)

	// loop until we find the title element and its content
	// or encounter an error (including the end of the file)

	for { // same as while loop
		// get the next token type
		tokenType := tokenizer.Next()

		// if it's an error token, we either reached the end of the file or the HTML was malformed
		if tokenType == html.ErrorToken {
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
			// Err returns the error associated with the most recent ErrorToken token.
			// This is typically io.EOF, meaning the end of tokenization.
		}

		// if this is a start tag token
		if tokenType == html.StartTagToken {
			// get the token
			token := tokenizer.Token()
			// if the name of the element is "title"
			if "title" == token.Data {
				// the next token should be the page title
				tokenType = tokenizer.Next()
				// ensure it is actually a text token
				if tokenType == html.TextToken {
					// report the page title and break out of the loop
					fmt.Println(tokenizer.Token().Data)
					break
				}
			}
		}

	}
}
