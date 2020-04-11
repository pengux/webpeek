package main

import (
	"net/http"
	"net/url"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

const (
	userAgent         = "webpeek/1.0"
	defaultReqTimeout = 10 * time.Second
)

type (
	peekedContent struct {
		url   *url.URL
		title string
		// metaDesc string
		h1s      []string
		markdown string
		// bodyPreview string
		err error
	}

	peeks struct {
		a      []*peekedContent
		curVal int
	}
)

var (
	titleSelector = cascadia.MustCompile("head > title")
	// metaDescSelector = cascadia.MustCompile(`head > meta[name="description"]`)
	h1Selector = cascadia.MustCompile("h1")

	c = http.Client{
		Timeout: defaultReqTimeout,
	}
)

func peek(url *url.URL) *peekedContent {
	p := &peekedContent{
		url: url,
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		p.err = err
		return p
	}
	req.Header.Add("User-Agent", userAgent)

	resp, err := c.Do(req)
	if err != nil {
		p.err = err
		return p
	}

	htmlDoc, err := html.Parse(resp.Body)
	cErr := resp.Body.Close()
	if err != nil {
		p.err = err
		return p
	}
	if cErr != nil {
		p.err = cErr
		return p
	}

	// Parse title
	if v := titleSelector.MatchFirst(htmlDoc); v != nil {
		p.title = extractTextContent(v)
	}

	// // Parse meta description
	// if v := metaDescSelector.MatchFirst(htmlDoc); v != nil {
	// 	for _, attr := range v.Attr {
	// 		if attr.Key == "content" {
	// 			p.metaDesc = attr.Val
	// 		}
	// 	}
	// }
	//
	// Parse h1 tags
	h1s := h1Selector.MatchAll(htmlDoc)
	for _, h1 := range h1s {
		p.h1s = append(p.h1s, extractTextContent(h1))
	}

	// Parse whole document to markdown
	converter := md.NewConverter("", true, nil)
	sel := goquery.NewDocumentFromNode(htmlDoc).Find("body")
	p.markdown = converter.Convert(sel)

	return p
}

func extractTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.FirstChild == nil {
		return ""
	}

	return extractTextContent(n.FirstChild)
}

func (p *peeks) Next() bool {
	p.curVal++

	return p.curVal < len(p.a)
}

func (p *peeks) Value() *peekedContent {
	return p.a[p.curVal]
}

func (p *peekedContent) Reload() {
	newP := peek(p.url)
	*p = *newP
}
