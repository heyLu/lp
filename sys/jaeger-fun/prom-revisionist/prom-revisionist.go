package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/prometheus/prometheus/promql/parser"
)

var config struct {
	Addr        string
	UpstreamURL string
}

func main() {
	flag.StringVar(&config.Addr, "addr", "localhost:19090", "Address to listen on")
	flag.StringVar(&config.UpstreamURL, "upstream", "http://localhost:9090", "Upstream Prometheus url")
	flag.Parse()

	if flag.NArg() == 1 {
		expr, err := parser.ParseExpr(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("before:", expr.Pretty(0))

		err = parser.Walk(Revisionist{}, expr, nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("after: ", expr.Pretty(0))
		return
	}

	upstreamURL, err := url.Parse(config.UpstreamURL)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("api call: %q %s %s", req.URL.String(), req.Header.Get("Content-Type"), req.Header.Get("Content-Length"))

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, req.Body)
		if err != nil {
			log.Printf("could not read request body: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		body := buf.Bytes()

		req.Body = io.NopCloser(buf)
		err = req.ParseForm()
		if err != nil {
			log.Printf("could not read form: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Println(req.PostForm)

		u := req.URL
		u.Scheme = upstreamURL.Scheme
		u.Host = upstreamURL.Host
		proxyReq, err := http.NewRequest(req.Method, u.String(), bytes.NewBuffer(body))
		if err != nil {
			log.Printf("failed to created request: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		proxyReq.Header = req.Header

		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			log.Printf("failed to created request: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for name, vals := range resp.Header {
			for _, val := range vals {
				w.Header().Add(name, val)
			}
		}
		w.WriteHeader(resp.StatusCode)

		var out io.Writer = w
		if resp.StatusCode != 200 {
			out = io.MultiWriter(w, os.Stdout)
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Printf("could not write body: %s", err)
			return
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		u := req.URL
		u.Scheme = upstreamURL.Scheme
		u.Host = upstreamURL.Host
		proxyReq, err := http.NewRequest(req.Method, u.String(), req.Body)
		if err != nil {
			log.Printf("failed to created request: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		proxyReq.Header = req.Header

		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			log.Printf("failed to created request: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for name, vals := range resp.Header {
			for _, val := range vals {
				w.Header().Add(name, val)
			}
		}
		w.WriteHeader(resp.StatusCode)

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("could not write body: %s", err)
			return
		}
	})

	log.Printf("Listening on http://%s", config.Addr)
	log.Fatal(http.ListenAndServe(config.Addr, nil))
}

type Revisionist struct{}

func (r Revisionist) Visit(node parser.Node, path []parser.Node) (parser.Visitor, error) {
	if node == nil && path == nil {
		return nil, nil
	}

	switch val := node.(type) {
	case *parser.AggregateExpr:
		for i, label := range val.Grouping {
			if label == "service_name" {
				val.Grouping[i] = "service"
			}
		}
	case *parser.VectorSelector:
		if val.Name == "calls_total" {
			val.Name = "my_calls_total"
		}

		for _, label := range val.LabelMatchers {
			if label.Name == "__name__" {
				label.Value = "my_calls_total"
			}

			if label.Name == "service_name" {
				label.Name = "service"
			}
		}
	}
	return r, nil
}
