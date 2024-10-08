package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/heyLu/lp/go/things/handler"
	"github.com/heyLu/lp/go/things/storage"
)

var settings struct {
	Addr   string
	DBPath string
}

func main() {
	flag.StringVar(&settings.Addr, "addr", "localhost:5000", "Address to listen on")
	flag.StringVar(&settings.DBPath, "db-path", "things.db", "Path to db file")
	flag.Parse()

	dbStorage, err := storage.NewDBStorage(context.Background(), "file:"+settings.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbStorage.Close()

	things := &Things{
		handlers: handler.All,
		storage:  dbStorage,
	}

	router := chi.NewRouter()

	router.Get("/*", things.HandleIndex)
	router.Get("/thing", things.HandleThing)
	router.Post("/thing", things.HandleThing)

	router.Get("/{kind}", things.HandleList)
	router.Get("/{kind}/{category}", things.HandleList)
	router.Get("/{kind}/{category}/{id}", things.HandleFind)

	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Listening on http://%s", settings.Addr)
	log.Fatal(http.ListenAndServe(settings.Addr, router))
}

type Things struct {
	handlers handler.Handlers

	storage storage.Storage
}

type Handler func(ctx context.Context, storage storage.Storage, namespace string, w http.ResponseWriter, input string, save bool) error

var ErrNotHandled = errors.New("not handled")

func (t *Things) HandleIndex(w http.ResponseWriter, req *http.Request) {
	pageWithContent(w, req, "", nil)
}

func pageWithContent(w http.ResponseWriter, req *http.Request, input string, content handler.Renderer) {

	fmt.Fprintf(w, `<!doctype html>
<html>
<head>
	<meta charset="utf-8" />
	<title>things</title>

	<link rel="stylesheet" href="/static/things.css" />
</head>

<body>
	<main>
		<form hx-post="/thing" hx-target="#answer" hx-indicator="#waiting">
			<input id="tell-me" name="tell-me" type="text" autofocus autocomplete="off" placeholder="tell me things"
				value=%q
				hx-get="/thing"
				hx-trigger="input changed delay:250ms"
				hx-target="#answer"
				hx-indicator="#waiting" />
			<input name="save" value="yes" hidden />
			<input type="submit" value="💾" />
		    <img id="waiting" class="htmx-indicator" src="/static/three-dots.svg" />
	    </form>

		<section id="answer">`, input)

	if content != nil {
		content.Render(req.Context(), w)
	}

	fmt.Fprintln(w, `
		</section>
	</main>

	<script src="/static/htmx.min.js"></script>
	<script src="/static/things.js"></script>
</body>
</html>`)
}

func (t *Things) HandleThing(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}

	tellMe := req.Form.Get("tell-me")

	namespace := "test" // FIXME: get from path/cookie/stuff (like trackl does)

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	save := req.Method == http.MethodPost

	handled := false
	for _, handler := range t.handlers {
		err := t.handle(ctx, handler, t.storage, namespace, w, tellMe, save)
		if err == ErrNotHandled {
			continue
		}

		handled = true

		if err != nil {
			fmt.Fprintln(w, err)
		}

		break
	}

	if !handled {
		fmt.Fprintln(w, "don't know what to do with that (yet)")
	}
}

func (t *Things) handle(ctx context.Context, hndl handler.Handler, storage storage.Storage, namespace string, w http.ResponseWriter, input string, save bool) error {
	kind, ok := hndl.CanHandle(input)
	if !ok {
		return ErrNotHandled
	}

	fmt.Fprintln(w, kind)

	thing, err := hndl.Parse(input)
	if err != nil {
		return err
	}

	row := thing.(handler.Thing).ToRow()
	row.Namespace = namespace // TODO: from context in Insert probably?  or get it from context and pass it in 🤔

	if save {
		err := storage.Insert(ctx, row)
		if err != nil {
			return err
		}

		fmt.Fprintln(w, "saved!")
	}

	seq := []handler.Renderer{}

	if row.Summary != "" {
		// TODO: thing.CanSave and only then preview?
		previewRenderer, err := hndl.(handler.Handler).Render(ctx, row)
		if err != nil {
			return err
		}

		seq = append(seq,
			previewRenderer,
			handler.HTMLRenderer("<hr />"),
		)
	}

	listRenderer, err := t.renderList(ctx, hndl.(handler.Handler), namespace, input)
	if err != nil {
		return err
	}
	seq = append(seq, listRenderer)

	renderer := handler.SequenceRenderer(seq)
	return renderer.Render(ctx, w)
}

func (t *Things) HandleList(w http.ResponseWriter, req *http.Request) {
	kindParam := chi.URLParam(req, "kind")
	kind, hndl := t.handlers.For(kindParam)
	if hndl == nil {
		http.Error(w, "unknown kind "+kindParam, http.StatusNotFound)
		return
	}

	// args := n.QueryArgs(make([]any, 0, 1)) // TODO: filter by category/first param

	input := kind

	namespace := "test"

	renderer, err := t.renderList(req.Context(), hndl.(handler.Handler), namespace, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageWithContent(w, req, input, renderer)
}

func (t *Things) renderList(ctx context.Context, hndl handler.Handler, namespace string, input string) (handler.Renderer, error) {
	rows, err := hndl.Query(ctx, t.storage, namespace, input)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []handler.Renderer{}
	for rows.Next() {
		var row storage.Row
		err := rows.Scan(&row)
		if err != nil {
			return nil, err
		}

		renderer, err := hndl.Render(ctx, &row)
		if err != nil {
			return nil, err
		}

		res = append(res, renderer)
	}

	return handler.ListRenderer(res), nil
}

func (t *Things) HandleFind(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "not implemented", http.StatusInternalServerError)
}
