package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
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
		fileName := os.Args[3]

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			lastVisit = time.Now()

			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				http.Error(w, "could not read file", http.StatusNotFound)
				log.Println(err)
			}

			style := `body { max-width: 50em; margin: 0 auto; }`
			fmt.Fprintf(w, "<!doctype html><html><head><meta charset=\"utf-8\" /><title>%s</title><style>%s</style></head><body>\n\n\n", fileName, style)
			w.Write(blackfriday.MarkdownCommon(data))
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
