package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"errors"

	"strconv"

	"golang.org/x/net/html"
)

func fetchHTML(URL string) (io.ReadCloser, error) {

	resp, err := http.Get(URL) // GET the URL

	if err != nil { // if there is an error, report it and exit
		return nil, err
	}

	// check response status code
	if resp.StatusCode != http.StatusOK {

		fmt.Println(resp.StatusCode)
		return nil, errors.New("response status code was " + strconv.Itoa(resp.StatusCode))

	}

	// check response content type
	ctype := resp.Header.Get("Content-Type")    // extract content type from header
	if !strings.HasPrefix(ctype, "text/html") { // check if ctype does not begin w/ "text/html"
		return nil, errors.New("response content type was " + ctype + ", not text/html")
	}

	return resp.Body, err // cannot simply return resp (which is a pointer to an http.Response object)
}

//extractTitle returns the content within the <title> element or an error
func extractTitle(body io.ReadCloser) (string, error) {

	// create a new tokenizer over the response Body
	tokenizer := html.NewTokenizer(body)

	// loop until we find the title element and its content
	// or encounter an error (including the end of the file)

	var title string

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
					title = tokenizer.Token().Data
					break
				}
			}
		}
	}

	return title, nil
}

//fetchTitle fetches the page title for a URL
func fetchTitle(URL string) (string, error) {
	//TODO: fetch the HTML, extract the title, and make sure the body gets closed
	respBody, err := fetchHTML(URL)
	if err != nil {
		return "", err
	}
	title, err := extractTitle(respBody)

	defer respBody.Close() // ensure the response body is closed, defer waits to close until main() returns
	// ******** now waits until enclosing fetchHTML function returns to close resp Body, not main

	return title, err
}

func main() {

	//if the caller didn't provide a URL to fetch...
	if len(os.Args) < 2 {
		//print the usage and exit with an error
		fmt.Printf("usage:\n  pagetitle <url>\n")
		os.Exit(1)
	}

	title, err := fetchTitle(os.Args[1])
	if err != nil {
		log.Fatalf("error fetching page title: %v\n", err)
	}

	fmt.Println(title)

}
