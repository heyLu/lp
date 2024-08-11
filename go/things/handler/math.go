package handler

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/heyLu/lp/go/things/storage"
)

var _ Handler = MathHandler{}

var mathRe = regexp.MustCompile(`([0-9]|eur|usd)`)

type MathHandler struct{}

func (mh MathHandler) CanHandle(input string) (string, bool) {
	return "math", mathRe.MatchString(input)
}

func (mh MathHandler) Parse(input string) (Thing, error) {
	return Math(input), nil
}

func (mh MathHandler) Query(ctx context.Context, db storage.Storage, namespace string, input string) (storage.Rows, error) {
	return db.Query(ctx, namespace, storage.Kind("help"))
}

func (mh MathHandler) Render(ctx context.Context, row *storage.Row) (Renderer, error) {
	cmd := exec.CommandContext(ctx, "qalc", "--terse", "--color=0", row.Summary)

	buf := new(bytes.Buffer)
	cmd.Stderr = buf
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf(buf.String())
	}

	return StringRenderer(buf.String()), nil
}

type Math string

func (m Math) ToRow() *storage.Row {
	return &storage.Row{
		Metadata: storage.Metadata{
			Kind: "math",
		},
		Summary: string(m),
	}
}
