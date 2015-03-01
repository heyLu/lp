package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	feed "github.com/SlyMarbo/rss"
)

// store posts somewhere (reverse chronical ordering, text file)
// - title, url, timestamp, via
// - periodically (or triggered if fetches get new posts, but with a delay?)
// periodic fetching
// - once per hour
// - at most n connections at a time
// web frontend
// - simple list with links
// - filterable with query params
// - support "live" search if i want to be fancy (not right now)
// - support json and edn output (and transit?)
// test (see feeds_test.go)

var commands = []string{"fetch-all", "fetch-one", "fetch-background", "help"}

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	flag.CommandLine.Parse(os.Args[2:])

	switch cmd {
	case "fetch-all":
		FetchAll()
	case "fetch-one":
		if len(flag.Args()) != 1 {
			fmt.Printf("Usage: %s [<options>] fetch-one <url>\n", os.Args[0])
			os.Exit(1)
		}

		FetchOne(flag.Args()[0])
	case "fetch-background":
		feeds, err := ReadConfig("config.txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		FetchBackground(feeds)
	case "help":
		printUsage()
		flag.PrintDefaults()
		fmt.Println("\nAvaillable commands:")
		for _, command := range commands {
			fmt.Println("\t", command)
		}
		os.Exit(0)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("Usage: %s [<options>] <cmd> [<args>]\n", os.Args[0])
}

func FetchOne(u string) {
	f, err := Fetch(u)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%s: %s", u, f.Title)
	if strings.TrimSpace(f.Description) != "" {
		fmt.Printf(" - %s", f.Description)
	}
	fmt.Printf(" (%d entries)\n", len(f.Items))
	for _, item := range f.Items {
		fmt.Printf("\t%s\n", item.Title)
	}
}

func FetchAll() {
	feeds, err := ReadConfig("config.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("fetching %d feeds...\n", len(*feeds))
	for _, fn := range *feeds {
		fmt.Printf("fetching %s\n", fn)
		f, err := Fetch(fn)
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

type FeedResult struct {
	*feed.Feed
	URL string
}

func FetchBackground(us *[]string) {
	canFetchCh := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		canFetchCh <- true
	}

	resultCh := make(chan FeedResult, 10)
	for _, u := range *us {
		go Fetcher(u, canFetchCh, resultCh)
	}

	feeds := make(map[string](*feed.Feed))
	go func() {
		for {
			fmt.Printf("storing %d feeds\n", len(feeds))
			StoreFeeds("feeds.json", feeds)
			time.Sleep(1 * time.Minute)
		}
	}()

	for {
		f := <-resultCh
		canFetchCh <- true

		if f.Feed == nil {
			fmt.Printf("%s: empty feed\n", f.URL)
		} else {
			fmt.Printf("%s - %d entries\n", f.UpdateURL, len(f.Items))
			feeds[f.URL] = f.Feed
		}
	}
}

func StoreFeeds(n string, feeds map[string](*feed.Feed)) {
	f, err := os.Create(n)
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	enc := json.NewEncoder(f)
	if err = enc.Encode(feeds); err != nil {
		fmt.Println(err)
	}
}

func Fetcher(u string, canFetchCh chan bool, resultCh chan FeedResult) {
	for {
		<- canFetchCh

		fmt.Printf("fetching %s\n", u)
		f, _ := Fetch(u)
		resultCh <- FeedResult{f, u}

		time.Sleep(1 * time.Minute)
	}
}

func Fetch(fn string) (*feed.Feed, error) {
	fu, err := GetFeedUrl(fn)
	if err != nil {
		return nil, err
	}
	fn = fu

	f, err :=  feed.Fetch(fn)
	if err != nil {
		return nil, err
	}

	return f, nil
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
