package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "webpeek",
		Usage:  "Accepts a list of URLs and return sneak peeks of the websites",
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	urls, err := getInputURLs(c)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	peeks := peeks{
		a: make([]*peekedContent, len(urls)),
	}

	for i, u := range urls {
		wg.Add(1)

		go func(index int, url *url.URL) {
			defer wg.Done()
			peeks.a[index] = peek(url)
		}(i, u)
	}

	wg.Wait()

	return draw(peeks)
}

// getInputURLs checks whether input are from arguments or Stdin and try to
// parses it to a list of url.URL and returns the list
func getInputURLs(c *cli.Context) ([]*url.URL, error) {
	var rawURLs []string

	// Check arguments
	rawURLs = c.Args().Slice()

	// If arguments list is empty then check Stdin
	if len(rawURLs) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// Break after empty row
			if scanner.Text() == "" {
				break
			}
			rawURLs = append(rawURLs, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("could not read from Stdin: %w", err)
		}
	}

	var urls []*url.URL
	for _, s := range rawURLs {
		parsed, err := url.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("could not parse '%s' into URL: %w", s, err)
		}

		urls = append(urls, parsed)

	}

	return urls, nil
}
