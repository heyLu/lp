package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"flag"
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
	Date    string `yaml:"date"`
	Type    string `yaml:"type"`
}

var flags struct {
	writeBack bool
	reverse   bool
	css       string
	title     string
}
var dataPath string = "blog.yaml"

var defaultStyle = `
article {
	margin-bottom: 1em;
}

article header {
	display: flex;
	align-items: center;
}

article h1 {
	margin: 0;
	margin-right: 1em;
}

article time {
	color: #666;
}

article img {
	max-width: 80vw;
	max-height: 50vh;
}
`

func init() {
	flag.BoolVar(&flags.writeBack, "write-back", false, "Rewrite the YAML file with the generated ids")
	flag.BoolVar(&flags.reverse, "reverse", false, "Reverse the order of the articles in the file")
	flag.StringVar(&flags.css, "css", defaultStyle, "Custom styles to use")
	flag.StringVar(&flags.title, "title", "A blog", "Custom title to use")
}

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		dataPath = flag.Arg(0)
	}
	f, err := os.Open(dataPath)
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

	fmt.Printf(`<!doctype html>
<html>
<head>
	<meta charset="utf-8" />
	<title>%s</title>
	<style>%s</style>
</head>

<body>
`, flags.title, flags.css)

	if flags.reverse {
		l := len(posts)
		reversePosts := make([]Post, l)
		for i := 0; i < l; i++ {
			reversePosts[i] = posts[l-i-1]
		}
		posts = reversePosts
	}

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
		case "image":
			err = imageTmpl.Execute(os.Stdout, post)
		case "song":
			err = songTmpl.Execute(os.Stdout, post)
		case "text":
			err = textTmpl.Execute(os.Stdout, post)
		default:
			fmt.Fprintf(os.Stderr, "Error: no output for type '%s'\n", post.Type)
			os.Exit(1)
		}
		if err != nil {
			exit(err)
		}
	}

	fmt.Printf("\n</body>\n</html>\n")

	if flags.writeBack {
		out, err := yaml.Marshal(posts)
		if err != nil {
			exit(err)
		}
		ioutil.WriteFile(dataPath, out, 0664)
	}
}

var funcs = template.FuncMap{
	"markdown": func(markdown string) template.HTML {
		return template.HTML(blackfriday.MarkdownCommon([]byte(markdown)))
	},
	"safe_url": func(s string) template.URL {
		return template.URL(s)
	},
}

var shellTmpl = template.Must(template.New("shell").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="shell">
	<header>
		<h1><code class="language-shell">{{ .Title }}</code></h1>
		{{- if .Date }}<time>{{ .Date }}</time>{{ end -}}
	</header>
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}
</article>
`))

var linkTmpl = template.Must(template.New("link").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="link">
	<header>
		<h1><a href="{{ .URL }}">{{ .Title }}</a></h1>
		{{- if .Date }}<time>{{ .Date }}</time>{{ end -}}
	</header>
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}
</article>
`))

var imageTmpl = template.Must(template.New("image").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="image">
	{{- if .Title }}
	<header>
		<h1>{{ .Title }}</h1>
		{{- if .Date }}<time>{{ .Date }}</time>{{ end -}}
	</header>
	{{- end }}
	<img src="{{ safe_url .URL }}" />
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}
</article>
`))

var songTmpl = template.Must(template.New("song").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="song">
	{{- if .Title }}
	<header>
		<h1>{{ .Title }}</h1>
		{{- if .Date }}<time>{{ .Date }}</time>{{ end -}}
	</header>
	{{- end }}
	<audio src="{{ safe_url .URL }}" controls>
		Your browser can't play {{ .URL }}.
	</audio>
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}
</article>
`))

var textTmpl = template.Must(template.New("text").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="text">
	<header>
		<h1>{{ .Title }}</h1>
		{{- if .Date }}<time>{{ .Date }}</time>{{ end -}}
	</header>
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
