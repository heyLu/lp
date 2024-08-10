package handler

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/heyLu/lp/go/things/storage"
)

var mathRe = regexp.MustCompile(`([0-9]|eur|usd)`)

type MathHandler struct{}

func (mh MathHandler) CanHandle(input string) (string, bool) {
	return "math", mathRe.MatchString(input)
}

func (mh MathHandler) Parse(input string) (Thing, error) {
	return Math(input), nil
}

type Math string

func (m Math) Args(args []any) []any {
	return append(args, string(m))
}

func (m Math) Render(ctx context.Context, _ storage.Storage, _ string, input string) (Renderer, error) {
	cmd := exec.CommandContext(ctx, "qalc", "--terse", "--color=0", input)

	buf := new(bytes.Buffer)
	cmd.Stderr = buf
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf(buf.String())
	}

	return StringRenderer(buf.String()), nil
}
