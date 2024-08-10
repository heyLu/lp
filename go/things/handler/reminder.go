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
	var reminder Reminder
	parts := strings.SplitN(input, " ", 3)

	var dur time.Duration
	var err error
	if len(parts) > 1 {
		dur, err = time.ParseDuration(parts[1])
		if err != nil {
			return nil, err
		}

		reminder.Date = time.Now().UTC().Add(dur).Truncate(time.Minute)
	}

	if len(parts) > 2 {
		reminder.Description = parts[2]
	}

	return &reminder, nil
}

type Reminder struct {
	*storage.Metadata

	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func (r Reminder) Until() time.Duration {
	return r.Date.Sub(time.Now()).Truncate(time.Minute)
}

func (r *Reminder) Args(args []any) []any {
	return append(args, r.Description, r.Date)

}

func (r *Reminder) Render(ctx context.Context, storage storage.Storage, namespace string, input string) (Renderer, error) {
	seq := []Renderer{}
	if r.Description != "" {
		seq = append(seq,
			TemplateRenderer{Template: reminderTemplate, Data: r},
			HTMLRenderer("<hr />"))
	}

	rows, err := storage.Query(ctx, namespace, "reminder", 2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []Renderer{}
	for rows.Next() {
		var reminder Reminder
		var date string
		meta, err := rows.Scan(&reminder.Description, &date)
		if err != nil {
			return nil, err
		}

		reminder.Date, err = time.Parse("2006-01-02 15:04:05Z07:00", date)
		if err != nil {
			return nil, err
		}

		reminder.Metadata = meta

		res = append(res, TemplateRenderer{Template: reminderTemplate, Data: reminder})
	}

	seq = append(seq, ListRenderer(res))

	return SequenceRenderer(seq), nil
}

var templateHelpers = template.FuncMap{
	"now":      func() time.Time { return time.Now() },
	"truncate": func(t time.Time) time.Time { return t.Truncate(time.Minute) },
}

var reminderTemplate = template.Must(template.New("").Funcs(templateHelpers).Parse(`
<section class="thing reminder">
{{ if .Metadata }}
	<time class="meta date-created" time="{{ .DateCreated }}" title="{{ .DateCreated }}">{{ .DateCreated.Format "2006-01-02 15:04:05" }}</time>
{{ end }}

	<span>
		<time time="{{ .Date }}" title="{{ .Date }}">in {{ .Until }}</time>
		{{ .Description }}
	</span>
</section>
`))
