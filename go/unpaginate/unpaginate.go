// `unpaginate` unpaginates json resources.
//
// The requested resource is assumed to return an array of JSON
// objects.  `unpaginate` prints a new array containing the JSON
// objects from all pages on stdout.
//
// Pagination is assumed to be in the format that the GitHub v3
// API uses:
//
//  HTTP/1.1 200 OK
//  ...
//  Link: <https://api.github.com/user/527119/repos?per_page=42&page=2>; rel="next", <https://api.github.com/user/527119/repos?per_page=42&page=2>; rel="last"
package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"os"
	"strings"
)

var config struct {
	userInfo string
}

func init() {
	flag.StringVar(&config.userInfo, "user", "", "The basic auth user information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <url>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, `Fetches JSON documents from a paginated resource
and returns a single JSON document.

`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "  -h, --help\n\tDisplay this message\n")
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	url := flag.Arg(0)

	authorization := ""
	if config.userInfo != "" {
		if !strings.Contains(config.userInfo, ":") {
			fmt.Fprint(os.Stderr, "Type in your password/token: ")
			password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			config.userInfo += ":" + string(password)
		}

		authorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(config.userInfo))
	}

	os.Stdout.WriteString("[\n")

	first := true
	for url != "" {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if authorization != "" {
			req.Header.Set("Authorization", authorization)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		decoder := json.NewDecoder(res.Body)
		_, err = decoder.Token()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for decoder.More() {
			if !first {
				os.Stdout.WriteString(", ")
			} else {
				first = false
			}

			var data interface{}
			err := decoder.Decode(&data)
			if err != nil {
				fmt.Println("decode:", err)
				os.Exit(1)
			}

			out, err := json.MarshalIndent(data, "  ", "  ")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			os.Stdout.Write(out)
		}

		_, err = decoder.Token()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		url = ""
		links := res.Header.Get("Link")
		if links != "" {
			for _, link := range strings.Split(links, ",") {
				link := strings.TrimSpace(link)
				start := strings.Index(link, "<")
				end := strings.Index(link, ">")
				if start != -1 && end != -1 && start < end &&
					strings.HasSuffix(link, "rel=\"next\"") {
					url = link[start+1 : end]
					break
				}
			}
		}
	}

	os.Stdout.WriteString("\n]\n")
}
