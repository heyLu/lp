package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

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
		handlers: []handler.Handler{
			// TODO: note ...
			// TODO: bookmark <url> note
			handler.ReminderHandler{},
			handler.TrackHandler{},
			handler.MathHandler{},
			// TODO: HandleSummary
			// HandleHelp,
		},
		storage: dbStorage,
	}

	http.HandleFunc("/", things.HandleIndex)
	http.HandleFunc("/thing", things.HandleThing)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Listening on http://%s", settings.Addr)
	log.Fatal(http.ListenAndServe(settings.Addr, nil))
}

type Things struct {
	handlers []handler.Handler

	storage storage.Storage
}

type Handler func(ctx context.Context, storage storage.Storage, namespace string, w http.ResponseWriter, input string, save bool) error

var ErrNotHandled = errors.New("not handled")

func (t *Things) HandleIndex(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, `<!doctype html>
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
				hx-get="/thing"
				hx-trigger="load, input changed delay:250ms"
				hx-target="#answer"
				hx-indicator="#waiting" />
			<input name="save" value="yes" hidden />
			<input type="submit" value="💾" />
		    <img id="waiting" class="htmx-indicator" src="/static/three-dots.svg" />
	    </form>

		<section id="answer">
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
	if tellMe == "" {
		fmt.Fprintln(w)
		return
	}

	namespace := "test" // FIXME: get from path/cookie/stuff (like trackl does)

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	save := req.Method == http.MethodPost

	for _, handler := range t.handlers {
		err := handle(ctx, handler, t.storage, namespace, w, tellMe, save)
		if err == ErrNotHandled {
			continue
		}

		if err != nil {
			fmt.Fprintln(w, err)
		}

		break
	}
}

func handle(ctx context.Context, handler handler.Handler, storage storage.Storage, namespace string, w http.ResponseWriter, input string, save bool) error {
	kind, ok := handler.CanHandle(input)
	if !ok {
		return ErrNotHandled
	}

	thing, err := handler.Parse(input)
	if err != nil {
		return err
	}

	if save {
		args := thing.Args(make([]any, 0, 10))
		_, err = storage.Insert(ctx, namespace, kind, args...)
		if err != nil {
			return err
		}

		fmt.Fprintln(w, "saved!")
	}

	renderer, err := thing.Render(ctx, storage, namespace, input)
	if err != nil {
		return err
	}

	return renderer.Render(ctx, w)
}

func HandleHelp(ctx context.Context, _ storage.Storage, _ string, w http.ResponseWriter, input string, _ bool) error {
	if input != "help" {
		fmt.Fprint(w, "don't know that thing, sorry.  ")
	}

	fmt.Fprintln(w, "try math, echo, ...")

	return nil
}
