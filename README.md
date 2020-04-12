# webpeek

A CLI tool which accepts a list of URLs and peek at their content. The content can then be displayed in a table format or as slides.


## Rationale

I subscribe to a couple of RSS feeds (e.g. Hacker News) which contains links to websites and need a quick way to determine of the link has something that should interest me. Previously, I would just open all links in a web browser but this seems to be too wasteful, both in time and computing resources so that's why this tool exists.


## Usage

```
echo "https://duckduckgo.com\nhttps://archlinux.org" | webpeek
```

or a bunch of URLs and open selected (see [Keybindings](#keybindings)) in Firefox

```
cat urls.txt | webpeek | xargs firefox -new-tab
```

## Keybindings

 Keys    | Feature                      |
---------|------------------------------|
 `l`     | Keep the URL                 |
 `q`     | Quit                         |
 `r`     | Reload the peek              |
 `Space` | Next without keeping the URL |
