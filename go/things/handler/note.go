package handler

import (
	"bytes"
	"context"
	"html/template"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"

	"github.com/heyLu/lp/go/things/storage"
)

var _ Handler = NoteHandler{}

type NoteHandler struct{}

func (nh NoteHandler) CanHandle(input string) (string, bool) {
	return "note", strings.HasPrefix(input, "note")
}

var urlRe = regexp.MustCompile(`(\w+)://[^ ]+`)

func (nh NoteHandler) Parse(input string) (Thing, error) {
	idx := strings.Index(input, " ")
	if idx == -1 {
		idx = len(input)
	}
	content := input[idx:]

	note := Note{
		Row: &storage.Row{
			Metadata: storage.Metadata{
				Kind: "note",
			},
			Summary: content,
		},
	}

	about := urlRe.FindString(content)
	if about != "" {
		note.Ref.String = about
	}

	return note, nil
}

func (nh NoteHandler) Query(ctx context.Context, db storage.Storage, namespace string, input string) (storage.Rows, error) {
	parts := strings.SplitN(input, " ", 2)
	if len(parts) == 1 {
		return db.Query(ctx, namespace, storage.Kind("note"))
	}
	return db.Query(ctx, namespace, storage.Kind("note"), storage.Match("summary", parts[1]))
}

func (nh NoteHandler) Render(ctx context.Context, row *storage.Row) (Renderer, error) {
	return TemplateRenderer{
		Template: noteTemplate,
		Data:     Note{Row: row},
	}, nil
}

type Note struct {
	*storage.Row
}

func (n Note) ToRow() *storage.Row { return n.Row }

var noteFuncs = template.FuncMap{
	"markdown": func(md string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		err := goldmark.Convert([]byte(md), buf)
		if err != nil {
			return "", err
		}
		return template.HTML(buf.String()), nil
	},
}

var noteTemplate = template.Must(template.New("").Funcs(noteFuncs).Parse(`
<section class="thing note">
{{ if .Metadata }}
	<time class="meta date-created" time="{{ .DateCreated }}" title="{{ .DateCreated }}">{{ .DateCreated.Format "2006-01-02 15:04:05" }}</time>
{{ end }}

	<div>{{ markdown .Summary }}</div>
</section>
`))
