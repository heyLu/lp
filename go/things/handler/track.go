package handler

import (
	"context"
	"html/template"
	"strconv"
	"strings"

	"github.com/heyLu/lp/go/things/storage"
)

var _ Handler = &Track{}

type Track struct {
	Category string   `json:"category"`
	Num      *float64 `json:"num"`
	Notes    *string  `json:"notes"`
}

func (t *Track) CanHandle(input string) (string, bool) {
	return "track", strings.HasPrefix(input, "track")
}
func (t *Track) Kind() string { return "track" }

func (t *Track) Parse(input string) error {
	parts := strings.SplitN(input, " ", 4)
	if len(parts) > 1 {
		t.Category = parts[1]
	}

	if len(parts) > 2 {
		num, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return err
		}
		t.Num = &num

	}

	if len(parts) > 3 {
		t.Notes = &parts[3]
	}

	return nil
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
	err := t.Parse(input)
	if err != nil {
		return nil, err
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
		_, err := rows.Scan(&track.Category, &track.Num, &track.Notes)
		if err != nil {
			return nil, err
		}

		res = append(res, TemplateRenderer{Template: trackTemplate, Data: track})
	}

	return ListRenderer(res), nil
}

var trackTemplate = template.Must(template.New("").Parse(`
<section class="track">
{{ .Category }}{{ if .Num }} <span{{ if (eq .Category "mood") }} style="opacity: calc({{ .Num }}/100)"{{ end }}>{{ .Num }}{{ if (eq .Category "stretch")}}min{{ end }}</span>{{ end }}{{ if .Notes }}<p>{{ .Notes }}</p>{{ end }}
</section>
`))
