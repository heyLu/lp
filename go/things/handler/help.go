package handler

import (
	"context"
	"strings"

	"github.com/heyLu/lp/go/things/storage"
)

var _ Handler = HelpHandler{}

type HelpHandler struct{}

func (h HelpHandler) CanHandle(input string) (string, bool) {
	return "help", strings.HasPrefix(input, "help")
}

func (h HelpHandler) Parse(input string) (Thing, error) {
	return Help(input), nil
}

func (h HelpHandler) Query(ctx context.Context, db storage.Storage, namespace string, input string) (storage.Rows, error) {
	return db.Query(ctx, namespace, storage.Kind("help"))
}

func (h HelpHandler) Render(ctx context.Context, row *storage.Row) (Renderer, error) {
	return StringRenderer(`help, what is this thing?!

some examples:

- reminder 30m go stretch a bit #health
- track sleep 7.0 okay, went to bed too late
- track mood 75 #tired
- 2**10
- 30usd to eur
`), nil
}

type Help string

func (h Help) ToRow() *storage.Row {
	return &storage.Row{
		Metadata: storage.Metadata{
			Kind: "help",
		},
		Summary: string(h),
	}
}
