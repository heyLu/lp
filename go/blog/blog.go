package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"

	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

type Post struct {
	Id      string `yaml:"id"`
	Title   string `yaml:"title"`
	URL     string `yaml:"url"`
	Content string `yaml:"content"`
	Type    string `yaml:"type"`
}

func main() {
	f, err := os.Open("blog.yaml")
	if err != nil {
		exit(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		exit(err)
	}

	var posts []Post
	err = yaml.Unmarshal(data, &posts)
	if err != nil {
		exit(err)
	}

	fmt.Printf(`<doctype html>
<html>
<head>
	<meta charset="utf-8" />
	<title>A blog</title>
</head>

<body>
`)

	for i, post := range posts {
		if post.Id == "" {
			posts[i].Id = generateId(post)
			post = posts[i]
		}

		var err error
		switch post.Type {
		case "shell":
			err = shellTmpl.Execute(os.Stdout, post)
		case "link":
			err = linkTmpl.Execute(os.Stdout, post)
		default:
			fmt.Fprintf(os.Stderr, "Error: no output for type '%s'\n", post.Type)
			os.Exit(1)
		}
		if err != nil {
			exit(err)
		}
	}

	fmt.Printf("\n</body>\n</html>\n")
}

var funcs = template.FuncMap{
	"markdown": func(markdown string) template.HTML {
		return template.HTML(blackfriday.MarkdownCommon([]byte(markdown)))
	},
}

var shellTmpl = template.Must(template.New("shell").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="shell">
	<h1><code class="language-shell">{{ .Title }}</code></h1>
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}
</article>
`))

var linkTmpl = template.Must(template.New("link").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="link">
	<h1><a href="{{ .URL }}">{{ .Title }}</a></h1>
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}
</article>
`))

func exit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func generateId(p Post) string {
	h := md5.New()
	io.WriteString(h, p.Title)
	io.WriteString(h, p.Content)
	io.WriteString(h, p.Type)
	return hex.EncodeToString(h.Sum(nil))
}

func randomId() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(buf)
}
