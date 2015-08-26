// `unpaginate` unpaginates json resources.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("Usage: %s [flags] <url>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	url := flag.Arg(0)

	os.Stdout.WriteString("[\n")

	for url != "" {
		res, err := http.Get(url)
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

		first := true
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
