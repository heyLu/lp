package handler

import (
	"context"
	"fmt"
	"html"
	"html/template"
	"net/http"

	"github.com/heyLu/lp/go/things/storage"
)

var All = Handlers([]Handler{
	// TODO: bookmark <url> note
	ReminderHandler{},
	TrackHandler{},
	NoteHandler{},
	ByDateHandler{},
	MathHandler{},
	// TODO: HandleSummary
	HelpHandler{},
	OverviewHandler{},
})

type Handlers []Handler

func (hs Handlers) For(kind string) (string, Handler) {
	for _, h := range hs {
		k, _ := h.CanHandle("")
		if k == kind {
			return kind, h
		}
	}

	// try if someone can handle it, e.g. if kind is 2024-08
	for _, h := range hs {
		_, ok := h.CanHandle(kind)
		if ok {
			return kind, h
		}
	}

	return "", nil
}

type Handler interface {
	CanHandle(input string) (string, bool)
	Parse(input string) (Thing, error)
}

type Thing interface {
	Args([]any) []any
	Render(ctx context.Context, storage storage.Storage, namespace string, input string) (Renderer, error)
}

type Renderer interface {
	Render(ctx context.Context, w http.ResponseWriter) error
}

type StringRenderer string

func (sr StringRenderer) Render(ctx context.Context, w http.ResponseWriter) error {
	fmt.Fprintln(w, "<pre>"+html.EscapeString(string(sr))+"</pre")
	return nil
}

type HTMLRenderer string

func (hr HTMLRenderer) Render(ctx context.Context, w http.ResponseWriter) error {
	fmt.Fprintln(w, hr)
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

type SequenceRenderer []Renderer

func (sr SequenceRenderer) Render(ctx context.Context, w http.ResponseWriter) error {
	for _, r := range sr {
		err := r.Render(ctx, w)
		if err != nil {
			return err
		}
	}
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
