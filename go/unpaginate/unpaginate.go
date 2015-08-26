// `unpaginate` unpaginates json resources.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("Usage: %s [flags] <url>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	url := flag.Arg(0)
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

	os.Stdout.WriteString("[\n")

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

	os.Stdout.WriteString("\n]")
}
