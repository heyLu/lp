package handler

import (
	"context"
	"strings"
	"time"

	"github.com/heyLu/lp/go/things/storage"
)

type OverviewHandler struct{}

func (mh OverviewHandler) CanHandle(input string) (string, bool) {
	return "overview", input == "" || strings.HasPrefix(input, "overview")
}

func (mh OverviewHandler) Parse(input string) (Thing, error) {
	return Overview(input), nil
}

type Overview string

func (m Overview) Args(args []any) []any {
	return append(args, string(m))
}

func (m Overview) Render(ctx context.Context, storage storage.Storage, namespace string, input string) (Renderer, error) {
	rows, err := storage.Query(ctx, namespace, "settings", 2, "overview.value")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	views := make([]string, 0, 5)
	for rows.Next() {
		var key, view string
		_, err := rows.Scan(&key, &view)
		if err != nil {
			return nil, err
		}

		views = append(views, view)
	}

	if len(views) == 0 {
		views = []string{
			time.Now().Format(time.DateOnly),
			"reminders",
			"help",
		}
	}

	renderers := make([]Renderer, 0, len(views))
	for _, view := range views {
		for _, handler := range All {
			_, ok := handler.CanHandle(view)
			if !ok {
				continue
			}

			thing, err := handler.Parse(view)
			if err != nil {
				return nil, err
			}

			r, err := thing.Render(ctx, storage, namespace, view)
			if err != nil {
				return nil, err
			}

			renderers = append(renderers, r)

			break
		}
	}

	return ListRenderer(renderers), nil
}
