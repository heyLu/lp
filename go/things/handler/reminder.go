package handler

import (
	"context"
	"html/template"
	"strings"
	"time"

	"github.com/heyLu/lp/go/things/storage"
)

var _ Handler = ReminderHandler{}
var _ Thing = &Reminder{}

type ReminderHandler struct{}

func (rh ReminderHandler) CanHandle(input string) (string, bool) {
	return "reminder", strings.HasPrefix(input, "remind")
}

func (rh ReminderHandler) Parse(input string) (Thing, error) {
	reminder := Reminder{
		Row: &storage.Row{
			Metadata: storage.Metadata{
				Kind: "reminder",
			},
		},
	}

	reminder.Bool.Bool = false

	parts := strings.SplitN(input, " ", 3)

	var dur time.Duration
	var err error
	if len(parts) > 1 {
		dur, err = time.ParseDuration(parts[1])
		if err != nil {
			return nil, err
		}

		reminder.Time.Time = time.Now().UTC().Add(dur).Truncate(time.Minute)
	}

	if len(parts) > 2 {
		reminder.Summary = parts[2]
	}

	return &reminder, nil
}

func (rh ReminderHandler) Query(ctx context.Context, db storage.Storage, namespace string, input string) (storage.Rows, error) {
	if !strings.Contains(input, " ") {
		return db.Query(ctx, namespace, storage.Kind("reminder"))
	}
	return db.Query(ctx, namespace, storage.Kind("reminder"), storage.Match("summary", input))
}

func (rh ReminderHandler) Render(ctx context.Context, row *storage.Row) (Renderer, error) {
	return TemplateRenderer{Template: reminderTemplate, Data: Reminder{Row: row}}, nil
}

type Reminder struct{ *storage.Row }

func (r *Reminder) ToRow() *storage.Row {
	return r.Row
}

func (r Reminder) Until() time.Duration {
	return r.Time.Time.Sub(time.Now()).Truncate(time.Minute)
}

var reminderTemplate = template.Must(template.New("").Parse(`
<section class="thing reminder">
{{ if .Metadata }}
	<time class="meta date-created" time="{{ .DateCreated }}" title="{{ .DateCreated }}">{{ .DateCreated.Format "2006-01-02 15:04:05" }}</time>
{{ end }}

	<span>
		<time time="{{ .Time }}" title="{{ .Time }}">in {{ .Until }}</time>
		{{ .Summary }}
	</span>
</section>
`))
