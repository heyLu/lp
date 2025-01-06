// Command blog renders a YAML file to a little HTML blog.
//
// Feel free to adapt it and use it in your setup, if you want your own
// personal digital writing implement.
package main

import (
	"bytes"
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
	"path"
	"sort"
	"strings"
	"unicode"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
)

// Post is a single post, title, tags and all.
//
// It contains all information to get rendered by something, as determined by
// the `Type`.
type Post struct {
	ID      string   `yaml:"id"`
	Title   string   `yaml:"title"`
	URL     string   `yaml:"url"`
	Content string   `yaml:"content"`
	Tags    []string `yaml:"tags"`
	Date    string   `yaml:"date"`
	Type    string   `yaml:"type"`
}

// Options are options that can be specified in the YAML file itself or on the
// commandline, to control how the output is done.
type Options struct {
	WriteBack      bool   `yaml:"write_back"`
	HashIDs        bool   `yaml:"hash_ids"`
	Reverse        bool   `yaml:"reverse"`
	CSS            string `yaml:"css"`
	NoDefaultStyle bool   `yaml:"no_default_style"`
	Title          string `yaml:"title"`
	After          string `yaml:"after"`
	PostsDir       string `yaml:"posts_dir"`
	Favicon        string `yaml:"favicon"`

	// options only in the YAML prefix:

	Output string `yaml:"output"`

	// options only for file output
	title    string
	showTags bool
}

var flags struct {
	Options
	PrintDefaultStyle bool
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
	margin-left: 0.5em;
	text-decoration: none;
	color: #555;

	visibility: hidden;
}

article header:hover .permalink {
	visibility: visible;
}

article time {
	margin-left: 0.5em;
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
	flag.BoolVar(&flags.WriteBack, "write-back", false, "Rewrite the YAML file with the generated ids")
	flag.BoolVar(&flags.HashIDs, "hash-ids", false, "Use hash-based post ids")
	flag.BoolVar(&flags.Reverse, "reverse", false, "Reverse the order of the articles in the file")
	flag.StringVar(&flags.Favicon, "favicon", "", "Link to favicon (can be a data:... url)")
	flag.StringVar(&flags.CSS, "css", "", "Use custom `css` styles")
	flag.BoolVar(&flags.NoDefaultStyle, "no-default-style", false, "Don't use the default styles")
	flag.BoolVar(&flags.PrintDefaultStyle, "print-default-style", false, "Print the default styles")
	flag.StringVar(&flags.Title, "title", "A blog", "Custom `title` to use")
	flag.StringVar(&flags.After, "after", "", "Insert additional `html` at the end of the generated page")
	flag.StringVar(&flags.PostsDir, "posts-dir", "", "Directory to write per-post html files to")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [<blog.yaml> [<blog.html>]]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flags.PrintDefaultStyle {
		fmt.Print(defaultStyle)
		os.Exit(0)
	}

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

	i := bytes.Index(data, []byte{'\n', '-', '-', '-', '\n'})
	if i != -1 && i != 0 {
		var opts Options
		err = yaml.Unmarshal(data[0:i+1], &opts)
		if err != nil {
			exit(err)
		}
		data = data[i+5:]

		isSet := map[string]bool{}
		flag.Visit(func(f *flag.Flag) {
			isSet[f.Name] = true
		})

		flag.VisitAll(func(f *flag.Flag) {
			_, ok := isSet[f.Name]
			if ok {
				return
			}

			switch f.Name {
			case "write-back":
				flags.WriteBack = opts.WriteBack
			case "hash-ids":
				flags.HashIDs = opts.HashIDs
			case "reverse":
				flags.Reverse = opts.Reverse
			case "favicon":
				flags.Favicon = opts.Favicon
			case "css":
				flags.CSS = opts.CSS
			case "no-default-style":
				flags.NoDefaultStyle = opts.NoDefaultStyle
			case "title":
				if opts.Title != "" {
					flags.Title = opts.Title
				}
			case "after":
				flags.After = opts.After
			case "posts-dir":
				if opts.PostsDir != "" {
					flags.PostsDir = opts.PostsDir
				}
			}
		})

		// write to specified `output` file if none was given as an argument
		if opts.Output != "" && out == os.Stdout {
			// write to same directory as the input file if opts.Output is
			// just the file name
			if path.Base(opts.Output) == opts.Output && path.Dir(dataPath) != "." {
				opts.Output = path.Join(path.Dir(dataPath), opts.Output)
			}

			out, err = os.Create(opts.Output)
			if err != nil {
				exit(err)
			}
		}
	}

	var posts []Post
	err = yaml.Unmarshal(data, &posts)
	if err != nil {
		exit(err)
	}

	for i, post := range posts {
		if post.ID == "" {
			posts[i].ID = generateID(post)
		}
	}

	flags.showTags = true
	writePosts(posts, out, flags.Options)

	if flags.PostsDir != "" {
		for _, post := range posts {
			postPath := path.Join(path.Dir(dataPath), flags.PostsDir, post.ID+".html")

			f, err := os.Create(postPath)
			if err != nil {
				exit(err)
			}

			flags.Title = post.Title
			flags.showTags = false
			writePosts([]Post{post}, f, flags.Options)

			fmt.Println("wrote", postPath)
		}
	}
}

func writePosts(posts []Post, out io.WriteCloser, opts Options) {
	if flags.NoDefaultStyle {
		defaultStyle = ""
	}
	fmt.Fprintf(out, `<!doctype html>
<html>
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<link rel="icon" href="%s" />
	<title>%s</title>
	<style>%s</style>
	<style>%s</style>
</head>

<body>
`, flags.Favicon, template.HTMLEscapeString(flags.Title), defaultStyle, flags.CSS)

	if flags.Reverse {
		l := len(posts)
		reversePosts := make([]Post, l)
		for i := 0; i < l; i++ {
			reversePosts[i] = posts[l-i-1]
		}
		posts = reversePosts
	}

	tags := map[string]bool{}
	for _, post := range posts {
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

			var provider string
			switch {
			case strings.HasSuffix(u.Path, ".mp4") || strings.HasSuffix(u.Path, ".ogv"):
				provider = "native"
			case (u.Host == "youtube.com" || u.Host == "www.youtube.com") && u.Query().Get("v") != "":
				provider = "youtube"
				post.URL = fmt.Sprintf("https://www.youtube.com/embed/%s", u.Query().Get("v"))
			case u.Host == "vimeo.com" && getVimeoID(u.Path) != "":
				provider = "vimeo"
				post.URL = fmt.Sprintf("https://player.vimeo.com/video/%s", getVimeoID(u.Path))
			default:
				exit(fmt.Errorf("unsupported video url '%s'", post.URL))
			}
			p := struct {
				Post
				Provider string
			}{post, provider}
			err = videoTmpl.Execute(out, p)
		default:
			fmt.Fprintf(os.Stderr, "Error: no output for type '%s'\n", post.Type)
			os.Exit(1)
		}
		if err != nil {
			exit(err)
		}
	}

	if opts.showTags {
		sortedTags := []string{}
		for tag := range tags {
			sortedTags = append(sortedTags, tag)
		}
		sort.Strings(sortedTags)
		err := tagsTmpl.Execute(out, sortedTags)
		if err != nil {
			exit(err)
		}
	}

	fmt.Fprintf(out, `

	<script>
	var baseTitle = document.title;

	var currentFilter = null;

	window.addEventListener("DOMContentLoaded", function(ev) {
		runFilterFromURL(document.location);
	});

	window.addEventListener("hashchange", function(ev) {
		runFilterFromURL(new URL(ev.newURL));
	});

	window.addEventListener("click", function(ev) {
		if (!ev.target.classList.contains("tag-link")) {
			return;
		}

		if (ev.target.href == "") {
			return;
		}

		var filter = filterFromURL(new URL(ev.target.href));
		if (isSameFilter(currentFilter, filter)) {
			clearFilter();
			location.hash = "";
			ev.preventDefault();
		} else {
			filterTag(tag);
		}
	});

	function isSameFilter(f1, f2) {
		return f1 && f2 && f1.type == f2.type && f1.argument == f2.argument;
	}

	function runFilterFromURL(u) {
		var filter = filterFromURL(u);
		if (filter == null) {
			clearFilter();
		} else {
			filter.run(filter.argument);
			currentFilter = filter;
		}
	}

	function filterFromURL(u) {
		if (u.hash.startsWith("#tag:")) {
			return { type: "tag", run: filterTag, argument: u.hash.substr(5) };
		} else if (u.hash.startsWith("#title:")) {
			return { type: "title", run: filterTitle, argument: u.hash.substr(7) };
		} else if (u.hash.startsWith("#id:")) {
			return { type: "id", run: filterId, argument: u.hash.substr(4) };
		} else if (u.hash.startsWith("#type:")) {
			return { type: "type", run: filterType, argument: u.hash.substr(6) };
		} else {
			return null;
		}
	}

	function filterTag(tag) {
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
			if (matches) {
				article.classList.remove("does-not-match");
			} else {
				article.classList.add("does-not-match");
			}
		}

		document.title = baseTitle + " (Posts tagged '" + tag + "')";
	}

	function filterTitle(match) {
		var match = match.toLowerCase();
		var articles = document.querySelectorAll("article");
		for (var i = 0; i < articles.length; i++) {
			var article = articles[i];
			var title = article.querySelector("header h1");
			if (title && title.textContent.toLowerCase().match(match)) {
				article.classList.remove("does-not-match");
			} else {
				article.classList.add("does-not-match");
			}
		}

		document.title = baseTitle + " (Posts matching '" + match + "')";
	}

	function filterId(id) {
		var articles = document.querySelectorAll("article");
		for (var i = 0; i < articles.length; i++) {
			var article = articles[i];
			if (article.id == id) {
				article.classList.remove("does-not-match");
			} else {
				article.classList.add("does-not-match");
			}
		}

		document.title = baseTitle + " (Only post id '" + id + "')";
	}

	function filterType(type) {
		var articles = document.querySelectorAll("article");
		for (var i = 0; i < articles.length; i++) {
			var article = articles[i];
			if (article.classList.contains(type)) {
				article.classList.remove("does-not-match");
			} else {
				article.classList.add("does-not-match");
			}
		}

		document.title = baseTitle + " (Only " + type + " posts)";
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

	if flags.After != "" {
		fmt.Fprintf(out, "\n%s\n", flags.After)
	}

	fmt.Fprintf(out, "\n</body>\n</html>\n")
	out.Close()

	if flags.WriteBack {
		dataOut, err := yaml.Marshal(posts)
		if err != nil {
			exit(err)
		}
		ioutil.WriteFile(dataPath, dataOut, 0664)
	}
}

var funcs = template.FuncMap{
	"permalink": func(post Post) string {
		if flags.PostsDir == "" {
			return "#" + post.ID
		}

		return path.Join(flags.PostsDir, post.ID+".html")
	},
	"markdown": func(markdown string) template.HTML {
		return template.HTML(
			blackfriday.Run([]byte(markdown),
				blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.Footnotes),
				blackfriday.WithRenderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{Flags: blackfriday.CommonHTMLFlags | blackfriday.FootnoteReturnLinks})),
			))
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
{{ $options := .Options }}
{{ define "title" }}
	<header>
		{{- if .Title }}
		<h1>{{ .Title }}</h1>{{ end }}
		<a class="permalink" href="{{ permalink . }}">∞</a>
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
<article id="{{ .ID }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
	<header>
		<h1><code class="language-shell">{{ .Title }}</code></h1>
		<a class="permalink" href="{{ permalink . }}">∞</a>
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
<article id="{{ .ID }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
	<header>
		<h1><a href="{{ .URL }}">{{ .Title }}</a></h1>
		<a class="permalink" href="{{ permalink . }}">∞</a>
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
<article id="{{ .ID }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
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
<article id="{{ .ID }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
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
<article id="{{ .ID }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
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
<article id="{{ .ID }}" class="{{ .Type }}" {{- if .Tags }} data-tags="{{ json .Tags }}"{{ end }}>
	{{- template "title" . }}
	{{ if (eq .Provider "youtube" "vimeo") -}}
	<iframe width="560" height="315" src="{{ safe_url .URL }}" frameborder="0" allowfullscreen loading="lazy"></iframe>
	{{- else if (eq .Provider "native") -}}
	<video src="{{ safe_url .URL }}" controls></video>
	{{- end }}
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

func generateID(p Post) string {
	if flags.HashIDs {
		return hashID(p)
	}
	return slugID(p)
}

func hashID(p Post) string {
	h := md5.New()
	io.WriteString(h, p.Title)
	io.WriteString(h, p.Content)
	io.WriteString(h, p.Type)
	return hex.EncodeToString(h.Sum(nil))
}

func randomID() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(buf)
}

var usedSlugs = map[string]int{}

func slugID(p Post) string {
	slug := toSlug(p.Title)
	n, ok := usedSlugs[slug]
	if ok {
		n++
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

func getVimeoID(p string) string {
	i := strings.LastIndex(p, "/")
	if i == -1 || len(p) == i+1 {
		return ""
	}

	return p[i+1:]
}
