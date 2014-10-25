package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	feed "github.com/SlyMarbo/rss"
)

func main() {
	feeds, err := ReadConfig("config.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("fetching %d feeds...\n", len(*feeds))
	for _, fn := range *feeds {
		fmt.Printf("fetching %s\n", fn)
		fu, err := GetFeedUrl(fn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fn = fu

		f, err :=  feed.Fetch(fn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%s: %s - %s (%d entries)\n", fn, f.Title, f.Description, len(f.Items))
		for i, item := range f.Items {
			if i < 10 {
				fmt.Printf("%s: %s\n", fn, item.Title)
			}
		}
	}
}

func GetFeedUrl(u string) (string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}

	if strings.Contains(resp.Header.Get("Content-Type"), "xml") {
		return u, nil
	}

	tree, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	sel := cascadia.MustCompile("link[rel=alternate][type*=xml]")
	alt := sel.MatchFirst(tree)
	if alt == nil {
		return "", errors.New("no feed link found")
	}

	altUrl, found := FindAttr("href", alt.Attr)
	if !found {
		return "", errors.New("missing link in alternate")
	}

	return ToAbsolute(resp.Request.URL, altUrl.Val), nil
}

func FindAttr(name string, attributes []html.Attribute) (*html.Attribute, bool) {
	for _, attr := range attributes {
		if attr.Key == name {
			return &attr, true
		}
	}
	return nil, false
}

func ToAbsolute(base *url.URL, rawUrl string) string {
	url, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}
	return base.ResolveReference(url).String()
}

func ReadConfig(fileName string) (*[]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0)

	r := bufio.NewReader(f)

	line, err := r.ReadString('\n')
	for err == nil {
		line = strings.TrimSpace(line)
		if line != "" && line[0] != '#' && line[0] != ';' {
			lines = append(lines, line)
		}

		line, err = r.ReadString('\n')
	}

	return &lines, nil
}
