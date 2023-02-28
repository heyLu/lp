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
	"strings"

	"github.com/prometheus/prometheus/model/labels"
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

		err = parser.Walk(&Revisionist{}, expr, nil)
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

	// TODO: support match[]?

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

		wasRewrite := false
		if len(req.PostForm) > 0 {
			query := req.PostForm.Get("query")
			if query != "" {
				expr, err := parser.ParseExpr(query)
				if err != nil {
					log.Printf("invalid query %q: %s", query, err)
				} else {
					before := expr.String()
					rev := &Revisionist{}
					err = parser.Walk(rev, expr, nil)
					if err != nil {
						log.Printf("could not rewrite: %s", err)
					} else {
						if rev.foundBucket {
							expr = &parser.BinaryExpr{
								Op:  parser.MUL,
								LHS: &parser.NumberLiteral{Val: 1000},
								RHS: expr,
							}
						}

						log.Printf("rewriting!\n%s\n=>\n%s", before, expr.String())
						req.PostForm.Set("query", expr.String())

						wasRewrite = true
					}
				}
			}

			body = []byte(req.PostForm.Encode())
		}

		u := req.URL
		u.Scheme = upstreamURL.Scheme
		u.Host = upstreamURL.Host
		// TODO: modify query in url.Query/url.RawQuery
		proxyReq, err := http.NewRequest(req.Method, u.String(), bytes.NewBuffer(body))
		if err != nil {
			log.Printf("failed to created request: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		proxyReq.Header = req.Header
		if wasRewrite {
			// TODO: allow keeping gzip and other encodings (handle them transparently)
			proxyReq.Header.Del("Accept-Encoding")
		}

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
			log.Printf("error %d", resp.StatusCode)
			out = io.MultiWriter(w, os.Stdout)
		}

		var in io.Reader = resp.Body
		if wasRewrite {
			log.Println("rewriting body")

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, resp.Body)
			if err != nil {
				log.Printf("could not write body: %s", err)
				return
			}

			if strings.Contains(buf.String(), `"service"`) {
				log.Println("rewriting service in response")
			}
			// TODO: rewrite by using streaming in some way
			res := strings.Replace(buf.String(), `"service"`, `"service_name"`, -1)
			res = strings.Replace(res, `"uri"`, `"operation"`, -1)

			buf.Reset()
			_, err = buf.WriteString(res)
			if err != nil {
				log.Printf("could not rewrite: %s", err)
				return
			}

			in = buf
		}

		_, err = io.Copy(out, in)
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

type Revisionist struct {
	foundBucket bool
}

func (r *Revisionist) Visit(node parser.Node, path []parser.Node) (parser.Visitor, error) {
	if node == nil && path == nil {
		return nil, nil
	}

	switch val := node.(type) {
	case *parser.AggregateExpr:
		for i, label := range val.Grouping {
			if label == "service_name" {
				val.Grouping[i] = "service"
			}

			if label == "operation" {
				val.Grouping[i] = "uri"
			}
		}
	case *parser.VectorSelector:
		if val.Name == "calls_total" {
			val.Name = "http_server_requests_seconds_count"
		}
		if val.Name == "latency_bucket" {
			val.Name = "http_server_requests_seconds_bucket"

			r.foundBucket = true
		}

		matchers := make([]*labels.Matcher, 0, len(val.LabelMatchers))
		for _, label := range val.LabelMatchers {
			if label.Name == "__name__" && label.Value == "calls_total" {
				label.Value = "http_server_requests_seconds_count"
			}
			if label.Name == "__name__" && label.Value == "latency_bucket" {
				label.Value = "http_server_requests_seconds_bucket"

				r.foundBucket = true
			}

			if label.Name == "service_name" {
				label.Name = "service"
			}

			if label.Name == "status_code" && label.Value == "STATUS_CODE_ERROR" {
				// label.Type = labels.MatchNotEqual
				// label.Name = "outcome"
				// label.Value = "SUCCESS"
				label.Type = labels.MatchEqual
				label.Value = "outcome"
				label.Value = "SERVER_ERROR"
			}

			if label.Name == "span_kind" {
				continue
			}

			matchers = append(matchers, label)
		}
		val.LabelMatchers = matchers
	}
	return r, nil
}
