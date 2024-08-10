package handler

import (
	"bytes"
	"context"
	"html/template"
	"regexp"
	"strings"

	"github.com/heyLu/lp/go/things/storage"
	"github.com/yuin/goldmark"
)

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
	return &Note{
		Content: content,
		About:   urlRe.FindString(content),
	}, nil
}

type Scanner interface {
	ScanRow(storage.Rows) (Renderer, error)
	NumParams() int
}

func (nh NoteHandler) ScanRow(rows storage.Rows) (Renderer, error) {
	var note Note
	meta, err := rows.Scan(&note.About, &note.Content)
	if err != nil {
		return nil, err
	}
	note.Metadata = meta

	return TemplateRenderer{Template: noteTemplate, Data: note}, nil
}

func (nh NoteHandler) NumParams() int {
	return 2
}

type Note struct {
	*storage.Metadata

	Content string `json:"content"`
	About   string `json:"about"`
}

func (n *Note) Args(args []any) []any {
	return append(args, n.About, n.Content)
}

func (n *Note) Render(ctx context.Context, storage storage.Storage, namespace string, input string) (Renderer, error) {
	seq := []Renderer{}
	if n.Content != "" {
		seq = append(seq,
			TemplateRenderer{Template: noteTemplate, Data: n}, // in-progress thing
			HTMLRenderer("<hr />"),
		)
	}

	// args := n.QueryArgs(make([]any, 0, 1)) // TODO: filter by url

	rows, err := storage.Query(ctx, namespace, "note", 2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []Renderer{}
	for rows.Next() {
		noteRenderer, err := NoteHandler{}.ScanRow(rows)
		if err != nil {
			return nil, err
		}

		res = append(res, noteRenderer)
	}

	seq = append(seq, ListRenderer(res))

	return SequenceRenderer(seq), nil
}

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

	<div>{{ markdown .Content }}</div>
</section>
`))
