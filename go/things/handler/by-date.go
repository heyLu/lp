package handler

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/heyLu/lp/go/things/storage"
)

var _ Handler = ByDateHandler{}

type ByDateHandler struct{}

var (
	byDateRe      = regexp.MustCompile(`^[0-9]{4}(-[0-9]{2}(-[0-9]{2})?)?$`)
	byDateFormats = []string{
		"2006-01-02",
		"2006-01",
		"2006",
	}
)

func (_ ByDateHandler) CanHandle(input string) (string, bool) {
	return "by-date", byDateRe.MatchString(input)
}

func (_ ByDateHandler) Parse(input string) (Thing, error) {
	// TODO: support '<from> to <to>' syntax for custom ranges
	for i, format := range byDateFormats {
		t, err := time.Parse(format, input)
		if err != nil {
			continue
		}

		byDate := &ByDate{
			input: input,
			from:  t,
		}

		switch i {
		case 0:
			byDate.to = t.AddDate(0, 0, 1)
			return byDate, nil
		case 1:
			byDate.to = t.AddDate(0, 1, 0)
			return byDate, nil
		case 2:
			byDate.to = t.AddDate(1, 0, 0)
			return byDate, nil
		default:
			break
		}
	}

	return nil, fmt.Errorf("can't parse %q", input)
}

func (bdh ByDateHandler) Query(ctx context.Context, db storage.Storage, namespace string, input string) (storage.Rows, error) {
	thing, err := bdh.Parse(input)
	if err != nil {
		return nil, err
	}

	byDate := thing.(*ByDate)

	return db.Query(ctx, namespace,
		storage.Gt("date_created", byDate.from.UTC().Unix()),
		storage.Lt("date_created", byDate.to.UTC().Unix()),
	)
}

func (_ ByDateHandler) Render(ctx context.Context, row *storage.Row) (Renderer, error) {
	return StringRenderer(row.Summary), nil
}

type ByDate struct {
	input string
	from  time.Time
	to    time.Time
}

func (bd ByDate) ToRow() *storage.Row {
	return &storage.Row{
		Metadata: storage.Metadata{
			Kind: "by-date",
		},
		Summary: bd.input,
		Fields: map[string]any{
			"from": bd.from,
			"to":   bd.to,
		},
	}
}
