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
	t.Row = &storage.Row{Metadata: storage.Metadata{Kind: "track"}}

	parts := strings.SplitN(input, " ", 4)
	if len(parts) > 1 {
		t.Summary = parts[1]
	}

	if len(parts) > 2 {
		num, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return nil, err
		}
		t.Float.Float64 = num
		t.Float.Valid = true

	}

	if len(parts) > 3 {
		t.Content.String = parts[3]
		t.Content.Valid = true
	}

	return &t, nil
}

func (th TrackHandler) Query(ctx context.Context, db storage.Storage, namespace string, input string) (storage.Rows, error) {
	thing, err := th.Parse(input)
	if err != nil {
		return nil, err
	}

	track := thing.(*Track)

	if track.Summary == "" {
		return db.Query(ctx, namespace, storage.Kind(track.Kind))
	}

	return db.Query(ctx, namespace, storage.Kind(track.Kind), storage.Summary(track.Summary))
}

func (th TrackHandler) Render(ctx context.Context, row *storage.Row) (Renderer, error) {
	return TemplateRenderer{Template: trackTemplate, Data: &Track{Row: row}}, nil
}

var _ Thing = &Track{}

type Track struct{ *storage.Row }

func (t *Track) Category() string { return t.Summary }
func (t *Track) Num() *float64 {
	if t.Float.Valid {
		return &t.Float.Float64
	} else {
		return nil
	}
}
func (t *Track) Notes() string { return t.Content.String }

func (t *Track) ToRow() *storage.Row { return t.Row }

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
