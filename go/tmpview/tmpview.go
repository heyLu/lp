package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

var (
	lastVisit time.Time
	srv       *http.Server
)

func main() {
	fileName := os.Args[1]
	addr := freePort()

	if os.Args[1] == "serve" {
		addr := os.Args[2]

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			lastVisit = time.Now()

			fileName := req.URL.Path[1:]
			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				http.Error(w, "could not read file", http.StatusNotFound)
				log.Println(err)
			}

			if strings.HasSuffix(strings.ToLower(fileName), ".pdf") {
				w.Write(data)
				return
			}

			style := `body { max-width: 50em; margin: 0 auto; }`
			fmt.Fprintf(w, "<!doctype html><html><head><meta charset=\"utf-8\" /><title>%s</title><style>%s</style></head><body>\n\n\n", fileName, style)
			switch {
			case strings.HasSuffix(strings.ToLower(fileName), ".md"):
				w.Write(blackfriday.MarkdownCommon(data))
			default:
				fmt.Fprint(w, "<pre>\n")
				template.HTMLEscape(w, data)
				fmt.Fprint(w, "\n</pre>")
			}
			fmt.Fprintf(w, "\n\n\n</body></html>")
		})

		srv = &http.Server{
			Handler: mux,
			Addr:    addr,
		}
		go keepalive()
		log.Fatal(srv.ListenAndServe())
	}

	server := exec.Command(os.Args[0], "serve", addr, fileName)
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("xdg-open", fmt.Sprintf("http://%s/%s", addr, fileName))
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func keepalive() {
	for {
		time.Sleep(1 * time.Minute)

		if time.Since(lastVisit) > 5*time.Minute {
			log.Println("shutting down")
			srv.Shutdown(context.TODO())
		}
	}
}

func freePort() string {
	l, err := net.Listen("tcp4", "localhost:0")
	if err != nil {
		panic(err)
	}
	l.Close()
	return l.Addr().String()
}
