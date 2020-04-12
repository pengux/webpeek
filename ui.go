package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func draw(peeks peeks) error {
	var out []string

	uiTextView := tview.NewTextView()
	uiTextView.
		SetDynamicColors(true).
		SetWrap(false).
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

		if p.err != nil {
			content = append(content, "[red]"+p.err.Error()+"[white]")
		}
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
