package handler

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/heyLu/lp/go/things/storage"
)

type ByDateHandler struct{}

var (
	byDateRe      = regexp.MustCompile(`[0-9]{4}(-[0-9]{2}(-[0-9]{2})?)?`)
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

type ByDate struct {
	input string
	from  time.Time
	to    time.Time
}

func (bd ByDate) Args(args []any) []any {
	return append(args, bd.input)
}

func (bd ByDate) Render(ctx context.Context, db storage.Storage, namespace string, input string) (Renderer, error) {
	rows, err := db.Query(ctx, namespace, "", 3,
		storage.Option{Field: "datetime(date_created, 'unixepoch')", Op: ">", Value: bd.from},
		storage.Option{Field: "datetime(date_created, 'unixepoch')", Op: "<", Value: bd.to})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	renderers := make([]Renderer, 0, 10)
	for rows.Next() {
		var val1, val2, val3 sql.NullString
		meta, err := rows.Scan(&val1, &val2, &val3)
		if err != nil {
			return nil, err
		}

		// TODO: support rendering a single Thing from Rows (also needed for /<thing>/<id> and more)
		renderers = append(renderers, StringRenderer(fmt.Sprintf("%s %s %s %s", meta.Kind, val1.String, val2.String, val3.String)))
	}

	return SequenceRenderer([]Renderer{
		StringRenderer("things from " + input),
		ListRenderer(renderers),
	}), nil
}
