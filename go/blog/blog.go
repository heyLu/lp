package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

type Post struct {
	Id      string   `yaml:"id"`
	Title   string   `yaml:"title"`
	URL     string   `yaml:"url"`
	Content string   `yaml:"content"`
	Tags    []string `yaml:"tags"`
	Date    string   `yaml:"date"`
	Type    string   `yaml:"type"`
}

var flags struct {
	writeBack bool
	hashIds   bool
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
}

article header .permalink {
	margin-left: 0.1em;
	text-decoration: none;
	color: #555;

	visibility: hidden;
}

article header:hover .permalink {
	visibility: visible;
}

article time {
	margin-left: 1em;
	color: #666;
}

article img {
	max-width: 80vw;
	max-height: 50vh;
}

.tags {
	display: flex;
	flex-wrap: wrap;
	list-style-type: none;
	padding: 0;
}

.tags .tag-link {
	color: black;
	margin-right: 0.5em;
}

article .tags .tag-link:visited {
	color: #555;
}

article.does-not-match {
	display: none;
}

#tags {
	color: #555;
	font-size: smaller;
	margin-bottom: 1em;
}

@media(max-width: 25em) {
	article h1 {
		margin-right: 0;
		text-align: center;
	}

	article header {
		flex-direction: column;
	}

	article header .permalink {
		visibility: visible;
	}

	article pre {
		white-space: pre-wrap;
	}

	iframe {
		max-width: 95vw;
	}
}

`

func init() {
	flag.BoolVar(&flags.writeBack, "write-back", false, "Rewrite the YAML file with the generated ids")
	flag.BoolVar(&flags.hashIds, "hash-ids", false, "Use hash-based post ids")
	flag.BoolVar(&flags.reverse, "reverse", false, "Reverse the order of the articles in the file")
	flag.StringVar(&flags.css, "css", defaultStyle, "Custom `css` styles to use")
	flag.StringVar(&flags.title, "title", "A blog", "Custom `title` to use")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [<blog.yaml> [<blog.html>]]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
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

	out := os.Stdout
	if flag.NArg() > 1 {
		out, err = os.Create(flag.Arg(1))
		if err != nil {
			exit(err)
		}
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		exit(err)
	}

	var posts []Post
	err = yaml.Unmarshal(data, &posts)
	if err != nil {
		exit(err)
	}

	fmt.Fprintf(out, `<!doctype html>
<html>
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
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

	tags := map[string]bool{}
	for i, post := range posts {
		if post.Id == "" {
			posts[i].Id = generateId(post)
			post = posts[i]
		}

		if post.Tags != nil {
			for _, tag := range post.Tags {
				tags[tag] = true
			}
		}

		var err error
		switch post.Type {
		case "shell":
			err = shellTmpl.Execute(out, post)
		case "link":
			err = linkTmpl.Execute(out, post)
		case "image":
			err = imageTmpl.Execute(out, post)
		case "song":
			err = songTmpl.Execute(out, post)
		case "text":
			err = textTmpl.Execute(out, post)
		case "video":
			u, err := url.Parse(post.URL)
			if err != nil {
				exit(fmt.Errorf("invalid video url '%s'", post.URL))
			}
			id := u.Query().Get("v")
			if (u.Host != "youtube.com" && u.Host != "www.youtube.com") || id == "" {
				exit(fmt.Errorf("unsupported video url '%s'", post.URL))
			}
			post.URL = id
			err = videoTmpl.Execute(out, post)
		default:
			fmt.Fprintf(os.Stderr, "Error: no output for type '%s'\n", post.Type)
			os.Exit(1)
		}
		if err != nil {
			exit(err)
		}
	}

	sortedTags := []string{}
	for tag, _ := range tags {
		sortedTags = append(sortedTags, tag)
	}
	sort.Strings(sortedTags)
	err = tagsTmpl.Execute(out, sortedTags)
	if err != nil {
		exit(err)
	}

	fmt.Fprintf(out, `

	<script>
	var baseTitle = document.title;

	var currentFilter = null;

	window.addEventListener("DOMContentLoaded", function(ev) {
		filterFromURL(document.location);
	});

	window.addEventListener("hashchange", function(ev) {
		filterFromURL(new URL(ev.newURL));
	});

	window.addEventListener("click", function(ev) {
		if (!ev.target.classList.contains("tag-link")) {
			return;
		}

		if (ev.target.href == "") {
			return;
		}

		var tag = tagFromURL(new URL(ev.target.href));
		if (currentFilter == tag) {
			clearFilter();
			location.hash = "";
			ev.preventDefault();
		} else {
			filterTag(tag);
		}
	});

	function filterFromURL(u) {
		var tag = tagFromURL(u);
		if (tag == null) {
			clearFilter();
		} else {
			filterTag(tag);
		}
	}

	function tagFromURL(u) {
		if (!u.hash.startsWith("#tag:")) {
			return null;
		}
		return u.hash.substr(5);
	}

	function filterTag(tag) {
		currentFilter = tag;

		var articles = document.querySelectorAll("article");
		for (var i = 0; i < articles.length; i++) {
			var article = articles[i];
			var matches = false;
			if (article && 'tags' in article.dataset) {
				var tags = JSON.parse(article.dataset.tags);
				for (var j = 0; j < tags.length; j++) {
					if (tags[j] == tag) {
						matches = true;
						break;
					}
				}
			}
			if (!matches) {
				article.classList.add("does-not-match");
			}
		}

		document.title = baseTitle + " (Posts tagged '" + tag + "')";
	}

	function filterId(id) {
		var articles = document.querySelectorAll("article");
		for (var i = 0; i < articles.length; i++) {
			var article = articles[i];
			if (article.id != id) {
				article.classList.add("does-not-match");
			}
		}

		document.title = baseTitle + " (Only post id '" + id + "')";
	}

	function clearFilter() {
		var articles = document.querySelectorAll("article.does-not-match");
		for (var i = 0; i < articles.length; i++) {
			articles[i].classList.remove("does-not-match");
		}

		currentFilter = null;
		document.title = baseTitle;
	}
	</script>`)
	fmt.Fprintf(out, "\n</body>\n</html>\n")
	out.Close()

	if flags.writeBack {
		dataOut, err := yaml.Marshal(posts)
		if err != nil {
			exit(err)
		}
		ioutil.WriteFile(dataPath, dataOut, 0664)
	}
}

var funcs = template.FuncMap{
	"markdown": func(markdown string) template.HTML {
		return template.HTML(blackfriday.MarkdownCommon([]byte(markdown)))
	},
	"safe_url": func(s string) template.URL {
		return template.URL(s)
	},
	"json": func(v interface{}) (string, error) {
		buf, err := json.Marshal(v)
		return string(buf), err
	},
}

var baseTmpl = template.Must(template.New("base").
	Funcs(funcs).Parse(`
{{ define "title" }}
	<header>
		{{- if .Title }}
		<h1>{{ .Title }}</h1>{{ end }}
		<a class="permalink" href="#{{ .Id }}">∞</a>
		{{- if .Date }}
		<time>{{ .Date }}</time>{{ end }}
	</header>
{{ end }}

{{ define "tags" }}
	<ul class="tags">
	{{- range . }}
		<li><a class="tag-link" href="#tag:{{ safe_url . }}">{{ . }}</a></li>
	{{ end -}}
	</ul>
{{ end }}
`))

var shellTmpl = template.Must(baseTmpl.New("shell").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
	<header>
		<h1><code class="language-shell">{{ .Title }}</code></h1>
		<a class="permalink" href="#{{ .Id }}">∞</a>
		{{- if .Date }}
		<time>{{ .Date }}</time>{{ end }}
	</header>
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}

	<footer>
		{{ template "tags" .Tags }}
	</footer>
</article>
`))

var linkTmpl = template.Must(baseTmpl.New("link").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
	<header>
		<h1><a href="{{ .URL }}">{{ .Title }}</a></h1>
		<a class="permalink" href="#{{ .Id }}">∞</a>
		{{- if .Date }}
		<time>{{ .Date }}</time>{{ end }}
	</header>
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}

	<footer>
		{{ template "tags" .Tags }}
	</footer>
</article>
`))

var imageTmpl = template.Must(baseTmpl.New("image").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
	{{- template "title" . }}
	<img src="{{ safe_url .URL }}" />
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}

	<footer>
		{{ template "tags" .Tags }}
	</footer>
</article>
`))

var songTmpl = template.Must(baseTmpl.New("song").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
	{{- template "title" . }}
	<audio src="{{ safe_url .URL }}" controls>
		Your browser can't play {{ .URL }}.
	</audio>
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}

	<footer>
		{{ template "tags" .Tags }}
	</footer>
</article>
`))

var textTmpl = template.Must(baseTmpl.New("text").
	Funcs(funcs).Parse(`
<article id="{{ .Id }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
	{{- template "title" . }}
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}

	<footer>
		{{ template "tags" .Tags }}
	</footer>
</article>
`))

var videoTmpl = template.Must(baseTmpl.New("video").Parse(`
<article id="{{ .Id }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
	{{- template "title" . }}
	<iframe width="560" height="315" src="https://www.youtube.com/embed/{{ .URL }}" frameborder="0" allowfullscreen></iframe>
	{{- if .Content }}

	{{ markdown .Content }}
	{{- end -}}

	<footer>
		{{ template "tags" .Tags }}
	</footer>
</article>
`))

var tagsTmpl = template.Must(baseTmpl.New("tags-list").Parse(`
	<section id="tags">
		<h1>All tags:</h1>
		{{ template "tags" . }}
	</section>
`))

func exit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func generateId(p Post) string {
	if flags.hashIds {
		return hashId(p)
	}
	return slugId(p)
}

func hashId(p Post) string {
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

var usedSlugs = map[string]int{}

func slugId(p Post) string {
	slug := toSlug(p.Title)
	n, ok := usedSlugs[slug]
	if ok {
		n += 1
	} else {
		n = 1
	}
	usedSlugs[slug] = n

	if slug != "" && n == 1 {
		return slug
	} else if slug == "" {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%s-%d", slug, n)
}

func toSlug(s string) string {
	lastChar := ' '
	s = strings.Map(func(ch rune) rune {
		var newChar rune
		switch {
		case unicode.IsLetter(ch) || unicode.IsDigit(ch):
			newChar = unicode.ToLower(ch)
		default:
			if lastChar == '-' {
				return -1
			}
			newChar = '-'
		}
		lastChar = newChar
		return newChar
	}, s)
	return strings.Trim(s, "-")
}
