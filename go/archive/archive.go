package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var flags struct {
	open bool
}

type Archive struct {
	Mappings map[string]string `json:"mappings"`
}

func init() {
	flag.BoolVar(&flags.open, "open", false, "Open the archived page")
}

func main() {
	flag.Parse()

	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		exit("url.Parse", err)
	}

	f, err := os.Open("archive.json")
	if err != nil {
		exit("os.Open", err)
	}

	var archive Archive
	dec := json.NewDecoder(f)
	err = dec.Decode(&archive)
	if err != nil {
		exit("dec.Decode", err)
	}
	f.Close()

	p, ok := archive.Mappings[u.String()]
	if ok {
		fmt.Println("==> Archived at", p)

		if flags.open {
			fmt.Println("==> Opening archive of", u)
			open(u, p)
		}

		return
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		fmt.Fprintf(os.Stderr, "Unknown url scheme %q\n", u.Scheme)
		os.Exit(1)
	}

	fmt.Println("==> Archiving", u)

	buf := make([]byte, 16)
	_, err = rand.Read(buf)
	if err != nil {
		exit("rand.Read", err)
	}

	cmd := exec.Command("prince", "--javascript", "--raster-output", fmt.Sprintf(".archive/%x-%%02d.png", buf), u.String())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		exit("prince", err)
	}

	parts, err := filepath.Glob(fmt.Sprintf(".archive/%x-*.png", buf))
	if err != nil {
		exit("filepath.Glob", err)
	}

	h := fmt.Sprintf(".archive/%x.html", buf)
	f, err = os.OpenFile(h, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0660)
	if err != nil {
		exit("os.OpenFile", err)
	}

	fmt.Fprintf(f, `<doctype html>
<html>
	<head>
		<title>%s</title>
	</head>

	<body>
`, u)

	wd, err := os.Getwd()
	if err != nil {
		exit("os.Getwd", err)
	}

	for _, p := range parts {
		fmt.Fprintf(f, "<img src=%q />\n", path.Join(wd, p))
	}

	fmt.Fprintf(f, "\n\t</body>\n</html>")
	f.Close()

	if archive.Mappings == nil {
		archive.Mappings = make(map[string]string, 1)
	}
	p = fmt.Sprintf("file://%s", path.Join(wd, h))
	archive.Mappings[u.String()] = p

	f, err = os.OpenFile("archive.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		exit("os.OpenFile", err)
	}

	enc := json.NewEncoder(f)
	err = enc.Encode(&archive)
	if err != nil {
		exit("enc.Encode", err)
	}
	f.Close()

	fmt.Println("==> Archived at", p)

	if flags.open {
		fmt.Println("==> Opening archive of", u)
		open(u, p)
	}
}

func open(u *url.URL, path string) {
	cmd := exec.Command("xdg-open", path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		exit("cmd.Start", err)
	}
}

func exit(msg string, err error) {
	fmt.Fprintf(os.Stderr, "Error: %s: %s\n", msg, err)
	os.Exit(1)
}
