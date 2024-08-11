package handler

import (
	"context"
	"strings"
	"time"

	"github.com/heyLu/lp/go/things/storage"
)

var _ Handler = OverviewHandler{}

type OverviewHandler struct{}

func (mh OverviewHandler) CanHandle(input string) (string, bool) {
	return "overview", input == "" || strings.HasPrefix(input, "overview")
}

func (mh OverviewHandler) Parse(input string) (Thing, error) {
	return Overview(input), nil
}

func (mh OverviewHandler) Query(ctx context.Context, db storage.Storage, namespace string, input string) (storage.Rows, error) {
	views := []string{
		time.Now().Format(time.DateOnly),
		"reminders",
		"help",
	}

	// TODO: need to do .Query for each of the views here and pass the into rows... (and then render the results in Render)

	return &overviewRows{summary: strings.Join(views, ",")}, nil
}

type overviewRows struct {
	idx     int
	summary string
}

func (o *overviewRows) Close() error { return nil }
func (o *overviewRows) Next() bool   { return o.idx < 1 }

func (o *overviewRows) Scan(row *storage.Row) error {
	o.idx += 1

	*row = storage.Row{
		Metadata: storage.Metadata{
			Kind: "overview",
		},
		Summary: o.summary,
	}

	return nil
}

func (mh OverviewHandler) Render(ctx context.Context, row *storage.Row) (Renderer, error) {
	views := strings.Split(row.Summary, ",")

	renderers := make([]Renderer, 0, len(views))
	for _, view := range views {
		_, handler := All.For(view)
		thing, err := handler.Parse(view)
		if err != nil {
			return nil, err
		}

		// TODO: query here to do a full render?

		renderer, err := handler.Render(ctx, thing.ToRow())
		if err != nil {
			return nil, err
		}

		renderers = append(renderers, renderer)
	}

	return ListRenderer(renderers), nil
}

type Overview string

func (o Overview) ToRow() *storage.Row {
	return &storage.Row{
		Metadata: storage.Metadata{
			Kind: "overview",
		},
		// Summary: "today,reminders,help",
	}
}
