package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
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

	var out []string

	uiTextView := tview.NewTextView()
	uiTextView.
		SetWordWrap(true).
		SetBorder(false)

	// Returns a new primitive which puts the provided primitive in the center and
	// sets its size to the given width and height.
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, false).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	err = screen.Init()
	if err != nil {
		return err
	}

	screenWidth, screenHeight := screen.Size()
	width := screenWidth / 3 * 2
	if width < 40 {
		width = 40
	}

	height := screenHeight / 3 * 2
	if height < 20 {
		height = 20
	}

	pages := tview.NewPages().
		AddPage("modal", modal(uiTextView, width, height), true, true)
	uiApp := tview.NewApplication().SetRoot(pages, true).
		SetScreen(screen).
		SetFocus(uiTextView)

	showPeek := func(p *peekedContent) {
		content := []string{p.url.String()}
		// if p.metaDesc != "" {
		// 	content = append(content, p.metaDesc)
		// }
		content = append(content, p.h1s...)
		content = append(content, p.markdown)

		// uiTextView.SetTitle(p.title)
		uiTextView.SetText(strings.Join(content, "\n\n")).
			ScrollToBeginning()
	}

	uiTextView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'l':
			out = append(out, peeks.Value().url.String())

			if !peeks.Next() {
				uiApp.Stop()
				return event
			}
		case 'q':
			uiApp.Stop()
			return event
		case 'r':
			peeks.Value().Reload()
		case 'h':
			fallthrough
		case ' ':
			if !peeks.Next() {
				uiApp.Stop()
				return event
			}
		}

		showPeek(peeks.Value())
		return event
	})

	showPeek(peeks.Value())
	err = uiApp.Run()
	if err != nil {
		return err
	}

	fmt.Println(strings.Join(out, "\n"))
	return nil
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
