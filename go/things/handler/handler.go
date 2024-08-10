package handler

import (
	"context"
	"fmt"
	"html"
	"html/template"
	"net/http"

	"github.com/heyLu/lp/go/things/storage"
)

type Handler interface {
	CanHandle(input string) (string, bool)
	Kind() string
	Render(ctx context.Context, storage storage.Storage, namespace string, input string) (Renderer, error)
}

type Renderer interface {
	Render(ctx context.Context, w http.ResponseWriter) error
}

type StringRenderer string

func (sr StringRenderer) Render(ctx context.Context, w http.ResponseWriter) error {
	fmt.Fprintln(w, "<pre>", html.EscapeString(string(sr)), "</pre")
	return nil
}

type ListRenderer []Renderer

func (lr ListRenderer) Render(ctx context.Context, w http.ResponseWriter) error {
	fmt.Fprintln(w, "<ul>")
	for _, r := range lr {
		fmt.Fprintln(w, "<li>")
		err := r.Render(ctx, w)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, "</li>")
	}
	fmt.Fprintln(w, "</ul>")
	return nil
}

type TemplateRenderer struct {
	*template.Template

	Metadata *storage.Metadata
	Data     any
}

func (tr TemplateRenderer) Render(ctx context.Context, w http.ResponseWriter) error {
	return tr.Template.Execute(w, tr.Data)
}
