package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var flags struct {
	open bool

	archiveDir string
}

type Archive struct {
	Mappings map[string]string `json:"mappings"`
}

func init() {
	flag.BoolVar(&flags.open, "open", false, "Open the archived page")

	wd, err := os.Getwd()
	if err != nil {
		exit("os.Getwd", err)
	}
	flags.archiveDir = path.Join(wd, ".archive")
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

	var archiver string
	var archiveFunc func(string, *url.URL) (string, error)

	switch {
	default:
		archiver, archiveFunc = "prince", archiveWithPrince
	}

	resultPath, err := archiveFunc(flags.archiveDir, u)
	if err != nil {
		exit(archiver, err)
	}

	if archive.Mappings == nil {
		archive.Mappings = make(map[string]string, 1)
	}
	p = fmt.Sprintf("file://%s", resultPath)
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

func archiveWithPrince(dir string, u *url.URL) (resultPath string, err error) {
	buf := make([]byte, 16)
	_, err = rand.Read(buf)
	if err != nil {
		exit("rand.Read", err)
	}

	cmd := exec.Command("prince", "--javascript", "--raster-output", fmt.Sprintf("%x-%%02d.png", buf), u.String())
	cmd.Dir = dir
	cmd.Stderr = prefixWriter("    | ", os.Stderr)
	cmd.Stdout = prefixWriter("    | ", os.Stdout)
	fmt.Print("    ", cmd.Args[0])
	for _, arg := range cmd.Args[1:] {
		fmt.Print(" ", arg)
	}
	fmt.Println()
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	parts, err := filepath.Glob(path.Join(dir, fmt.Sprintf("%x-*.png", buf)))
	if err != nil {
		return "", fmt.Errorf("filepath.Glob: %s", err)
	}

	if len(parts) == 0 {
		return "", fmt.Errorf("filepath.Glob: no matches")
	}

	h := path.Join(dir, fmt.Sprintf("%x.html", buf))
	f, err := os.OpenFile(h, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0660)
	if err != nil {
		return "", fmt.Errorf("os.OpenFile: %s", err)
	}

	fmt.Fprintf(f, `<doctype html>
<html>
	<head>
		<title>%s</title>
	</head>

	<body>
`, u)

	for _, p := range parts {
		fmt.Fprintf(f, "<img src=%q />\n", p)
	}

	fmt.Fprintf(f, "\n\t</body>\n</html>")
	f.Close()

	return h, nil
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

type prefixLineWriter struct {
	prefix      string
	needsPrefix bool

	w io.Writer
}

func prefixWriter(prefix string, w io.Writer) io.Writer {
	return &prefixLineWriter{
		prefix:      prefix,
		needsPrefix: true,

		w: w,
	}
}

func (p *prefixLineWriter) Write(b []byte) (n int, err error) {
	if p.needsPrefix {
		p.w.Write([]byte(p.prefix))
		p.needsPrefix = false
	}

	n = 0
	for {
		i := bytes.IndexByte(b, '\n')
		if i == -1 {
			nn, err := p.w.Write(b)
			return n + nn, err
		}

		if i+1 == len(b) {
			p.needsPrefix = true
			nn, err := p.w.Write(b)
			return n + nn, err
		}

		nn, err := p.w.Write(b[:i+1])
		if err != nil {
			return n + nn, err
		}
		n += nn
		p.w.Write([]byte(p.prefix))
		b = b[i+1:]
	}
}
