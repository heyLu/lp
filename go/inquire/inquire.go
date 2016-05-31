package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
)

type PageInfo struct {
	RawURL string `json:"url"`
	// Title is the value of the `title` element, or the value of the
	// meta tag `title` or `twitter:title`.
	Title string `json:"title"`
	// Description is the value of the meta tag `description` or
	// `og:description`.
	Description string `json:"description,omitempty"`
	// Image is the value of the meta tag `og:description`.
	//
	// Note that this is expected to be an image, not the icon of the
	// webpage.
	Image string `json:"image,omitempty"`
}

var config = struct {
	output string
}{}

func init() {
	flag.StringVar(&config.output, "output", "text", "what format to output")
}

func main() {
	flag.Parse()

	u := flag.Args()[0]
	url, err := url.Parse(u)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if url.Scheme == "" {
		u = "http://" + u
	}

	info, err := GetPageInfo(u)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch config.output {
	case "text":
		fmt.Printf("url: %s\ntitle: %s\ndescription: %s\nimage: %s\n",
			info.RawURL, info.Title, info.Description, info.Image)
	case "html":
		fmt.Printf("<h1><a href=\"%s\">%s</a></h1>\n", info.RawURL, info.Title)
		if info.Image != "" {
			fmt.Printf("<img src=\"%s\" />\n", info.Image)
		}
		if info.Description != "" {
			fmt.Printf("<p>%s</p>\n", info.Description)
		}
	case "json":
		out, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Stdout.Write(out)
	case "yaml":
		fmt.Printf("- url: %s\n", info.RawURL)
		fmt.Printf("  title: %s\n", info.Title)
		if info.Description != "" {
			fmt.Printf("  description: %s\n", info.Description)
		}
		if info.Image != "" {
			fmt.Printf("  image: %s\n", info.Image)
		}
	default:
		fmt.Fprintln(os.Stderr, "unknown output format:", config.output)
		os.Exit(1)
	}
}

func GetPageInfo(u string) (*PageInfo, error) {
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	tree, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	sel := cascadia.MustCompile("meta")
	meta := sel.MatchAll(tree)

	found, title := findTitle(tree)
	if !found {
		_, title = findProperty(meta, "title", "twitter:title")
	}

	_, description := findProperty(meta, "description", "og:description")
	_, image := findProperty(meta, "og:image")

	return &PageInfo{
		Title:       title,
		Description: description,
		Image:       image,
		RawURL:      u,
	}, nil
}

func findTitle(tree *html.Node) (found bool, title string) {
	sel := cascadia.MustCompile("title")
	node := sel.MatchFirst(tree)
	if node == nil {
		return false, ""
	}

	if node.Type == html.ElementNode {
		node = node.FirstChild
	}

	buf := new(bytes.Buffer)
	for node != nil {
		if node.Type == html.TextNode {
			buf.WriteString(node.Data)
		}

		node = node.NextSibling
	}

	return true, string(buf.Bytes())
}

func findProperty(nodes []*html.Node, properties ...string) (found bool, value string) {
	props := make(map[string]struct{}, len(properties))
	for _, prop := range properties {
		props[prop] = struct{}{}
	}

	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key != "property" && attr.Key != "name" {
				continue
			}

			if _, ok := props[attr.Val]; ok {
				found, value := findAttr("content", node)
				if found {
					return true, value
				}
			}
		}
	}

	return false, ""
}

func findAttr(name string, node *html.Node) (bool, string) {
	if node == nil {
		return false, ""
	}

	for _, attr := range node.Attr {
		if attr.Key == name {
			return true, attr.Val
		}
	}

	return false, ""
}
