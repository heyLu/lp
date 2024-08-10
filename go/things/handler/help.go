package handler

import (
	"context"
	"strings"

	"github.com/heyLu/lp/go/things/storage"
)

type HelpHandler struct{}

func (mh HelpHandler) CanHandle(input string) (string, bool) {
	return "help", strings.HasPrefix(input, "help")
}

func (mh HelpHandler) Parse(input string) (Thing, error) {
	return Help(input), nil
}

type Help string

func (m Help) Args(args []any) []any {
	return append(args, string(m))
}

func (m Help) Render(ctx context.Context, _ storage.Storage, _ string, input string) (Renderer, error) {
	return StringRenderer(`help, what is this thing?!

some examples:

- reminder 30m go stretch a bit #health
- track sleep 7.0 okay, went to bed too late
- track mood 75 #tired
- 2**10
- 30usd to eur
`), nil
}
