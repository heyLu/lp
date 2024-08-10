package handler

import (
	"context"
	"html/template"
	"strconv"
	"strings"

	"github.com/heyLu/lp/go/things/storage"
)

var _ Handler = TrackHandler{}

type TrackHandler struct{}

func (th TrackHandler) CanHandle(input string) (string, bool) {
	return "track", strings.HasPrefix(input, "track")
}
func (th TrackHandler) Parse(input string) (Thing, error) {
	var t Track

	parts := strings.SplitN(input, " ", 4)
	if len(parts) > 1 {
		t.Category = parts[1]
	}

	if len(parts) > 2 {
		num, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return nil, err
		}
		t.Num = &num

	}

	if len(parts) > 3 {
		t.Notes = &parts[3]
	}

	return &t, nil
}

var _ Thing = &Track{}

type Track struct {
	*storage.Metadata

	Category string   `json:"category"`
	Num      *float64 `json:"num"`
	Notes    *string  `json:"notes"`
}

func (t *Track) Args(args []any) []any {
	return append(args, t.Category, t.Num, t.Notes)
}

func (t *Track) QueryArgs(args []any) []any {
	if t.Category != "" {
		return append(args, t.Category)
	}
	return args
}

func (t *Track) Render(ctx context.Context, storage storage.Storage, namespace string, input string) (Renderer, error) {
	seq := []Renderer{}
	if t.Category != "" && strings.Index(input, t.Category)+len(t.Category) != len(input) {
		seq = append(seq,
			TemplateRenderer{Template: trackTemplate, Metadata: nil, Data: t}, // in-progress thing
			HTMLRenderer("<hr />"),
		)
	}

	args := t.QueryArgs(make([]any, 0, 1))

	rows, err := storage.Query(ctx, namespace, "track", 3, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []Renderer{}
	for rows.Next() {
		var track Track
		meta, err := rows.Scan(&track.Category, &track.Num, &track.Notes)
		if err != nil {
			return nil, err
		}
		track.Metadata = meta

		res = append(res, TemplateRenderer{Template: trackTemplate, Metadata: meta, Data: track})
	}

	seq = append(seq, ListRenderer(res))

	return SequenceRenderer(seq), nil
}

var trackTemplate = template.Must(template.New("").Parse(`
<section class="thing track">
{{ if .Metadata }}
	<time class="meta date-created" time="{{ .DateCreated }}" title="{{ .DateCreated }}">{{ .DateCreated.Format "2006-01-02 15:04:05" }}</time>
{{ end }}

	<span>
{{ .Category }}{{ if .Num }} <span{{ if (eq .Category "mood") }} style="opacity: calc({{ .Num }}/100)"{{ end }}>{{ .Num }}{{ if (eq .Category "stretch")}}min{{ end }}</span>{{ end }}{{ if .Notes }}<p>{{ .Notes }}</p>{{ end }}
	</span>
</section>
`))
