package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

var flags struct {
	addr            string
	proxyURL        string
	proxyClientCert string
	proxyClientKey  string
}

func init() {
	flag.StringVar(&flags.addr, "addr", "localhost:8080", "Address to listen on")
	flag.StringVar(&flags.proxyURL, "proxy-url", "", "Proxy requests to this URL")
	flag.StringVar(&flags.proxyClientCert, "proxy-client-cert", "", "Client certificate to use when connecting to proxy")
	flag.StringVar(&flags.proxyClientKey, "proxy-client-key", "", "Client key to use when connecting to proxy")
}

var responses = []Response{}

func main() {
	flag.Parse()

	if flag.NArg() == 1 {
		rs, err := loadResponses(flag.Arg(0))
		if err != nil {
			log.Fatalf("Error: Parsing %s: %s", flag.Arg(0), err)
		}
		responses = rs
	}

	requestLog := make([]LogEntry, 0)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var resp *http.Response
		if flags.proxyURL != "" {
			resp = respondWithProxy(flags.proxyURL, w, req)
		} else {
			resp = respondWithStub(responses, w, req)
		}

		userAgent := req.Header.Get("User-Agent")
		log.Printf("%s %s - %d (%s, %q)", req.Method, req.URL, resp.StatusCode, req.RemoteAddr, userAgent)

		out, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Printf("Error: Dumping request: %s", err)
			return
		}

		if resp.Header.Get("Content-Type") == "application/json" {
			pretty, err := prettyfyJSON(resp.Body)
			if err != nil {
				log.Printf("Error: Prettyfying JSON: %s", err)
			} else {
				resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(pretty)))
				resp.Body = ioutil.NopCloser(bytes.NewReader(pretty))
			}
		}

		requestLog = append(requestLog, LogEntry{
			Request:  out,
			Response: asResponse(req, resp),
		})
	})

	http.HandleFunc("/_log", func(w http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.Header.Get("Accept"), "application/yaml") {
			rs := make([]Response, len(requestLog))
			for i, log := range requestLog {
				rs[i] = log.Response
			}
			err := renderYAML(w, rs)
			if err != nil {
				log.Printf("Error: Render log: %s", err)
			}
			return
		}

		for i := len(requestLog) - 1; i >= 0; i-- {
			w.Write([]byte("------------------------------------------------------\n"))
			w.Write(requestLog[i].Request)
			w.Write([]byte("\n\n"))
			requestLog[i].Response.AsHTTP().Write(w)
			w.Write([]byte("\n"))
		}
	})

	http.HandleFunc("/_stub", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			err := stubTmpl.Execute(w, nil)
			if err != nil {
				log.Printf("Error: Rendering stub template: %s", err)
				return
			}
		case "POST":
			err := req.ParseForm()
			if err != nil {
				log.Printf("Error: Parsing form: %s", err)
				return
			}
			responses = append(responses, readResponse(req.Form))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/_stubs", func(w http.ResponseWriter, req *http.Request) {
		var err error
		if strings.Contains(req.Header.Get("Accept"), "application/yaml") {
			err = renderYAML(w, responses)
		} else {
			err = renderHTML(w, responses)
		}
		if err != nil {
			log.Printf("Error: Rendering stubs: %s", err)
		}
	})

	http.HandleFunc("/_help", func(w http.ResponseWriter, req *http.Request) {
		urls := []struct {
			URL     string
			Summary string
		}{
			{URL: "/_log", Summary: "View all received requests with responses"},
			{URL: "/_stub", Summary: "Add a new response stub"},
			{URL: "/_stubs", Summary: "View all defined stubs"},
			{URL: "/_help", Summary: "This help"},
		}
		fmt.Fprint(w, `<!doctype html>
<html>
	<head>
		<title>/_help</title>
	</head>
	<body>
		<ul>`)
		for _, url := range urls {
			fmt.Fprintf(w, "<li><pre><a href=\"%s\">%s</a> - %s</pre></li>", url.URL, url.URL, url.Summary)
		}
		fmt.Fprint(w, `
		</ul>
	</body>
</html`)
	})

	log.Printf("Listening on http://%s", flags.addr)
	log.Printf("See http://%s/_help", flags.addr)
	log.Fatal(http.ListenAndServe(flags.addr, nil))
}

func matchResponse(req *http.Request, responses []Response) *Response {
	for _, resp := range responses {
		if req.Method == resp.Method && req.URL.Path == resp.Path {
			return &resp
		}
	}
	return nil
}

func respondWithStub(responses []Response, w http.ResponseWriter, req *http.Request) *http.Response {
	resp := matchResponse(req, responses)
	if resp == nil {
		resp = &Response{Status: 404, Body: "Not found"}
	}

	for _, header := range resp.Headers {
		w.Header().Set(header.Name, header.Value)
	}
	if resp.Status == 0 {
		resp.Status = 200
	}
	w.WriteHeader(resp.Status)
	w.Write([]byte(resp.Body))

	return resp.AsHTTP()
}

func respondWithProxy(proxyURL string, w http.ResponseWriter, req *http.Request) *http.Response {
	/*proxyReq, err := http.NewRequest(req.Method, proxyURL+req.URL.Path, nil)
	if err != nil {
		log.Printf("Error: Creating proxy request: %s", err)
		return nil
	}
	for name, vals := range req.Header {
		proxyReq.Header[name] = vals
	}*/

	proxyTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
				if flags.proxyClientCert != "" && flags.proxyClientKey != "" {
					cert, err := tls.LoadX509KeyPair(flags.proxyClientCert, flags.proxyClientKey)
					if err != nil {
						return nil, err
					}
					return &cert, nil
				}
				return &tls.Certificate{}, nil
			},
			InsecureSkipVerify: true,
		},
	}
	proxyClient := &http.Client{Transport: proxyTransport}

	u, err := url.Parse(proxyURL)
	if err != nil {
		log.Printf("Error: Parsing proxy url: %s", err)
		return nil
	}

	req.URL.Scheme = u.Scheme
	req.URL.Host = u.Host
	req.RequestURI = ""
	resp, err := proxyClient.Do(req)
	if err != nil {
		log.Printf("Error: Proxying %s: %s", req.URL.Path, err)
		return nil
	}
	defer resp.Body.Close()

	for name, vals := range resp.Header {
		w.Header()[name] = vals
	}
	w.WriteHeader(resp.StatusCode)

	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
	io.Copy(w, buf)

	return resp
}

func prettyfyJSON(r io.Reader) ([]byte, error) {
	dec := json.NewDecoder(r)
	var val interface{}
	err := dec.Decode(&val)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(val, "", "    ")
}

func renderYAML(w http.ResponseWriter, responses []Response) error {
	out, err := yaml.Marshal(responses)
	if err != nil {
		return err
	}
	w.Write(out)
	return nil
}

func renderHTML(w http.ResponseWriter, responses []Response) error {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<!doctype html><html><head><style>pre{max-width:100vw;padding:0.5em;background-color:#eee;white-space:pre-wrap;}</style></head><body><ul>\n")
	for _, resp := range responses {
		fmt.Fprintf(w, "<li><pre>%s</pre></li>\n", resp.String())
	}
	fmt.Fprintf(w, "\n</ul></body></html>")
	return nil
}

// LogEntry is a request/respond pair for logging.
type LogEntry struct {
	Request  Request
	Response Response
}

// Request is a stored serialized HTTP request.
type Request []byte

// Response is a mocked HTTP response.
type Response struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`

	Status  int      `yaml:"status"`
	Headers []Header `yaml:"headers"`
	Body    string   `yaml:"body"`
}

func (resp Response) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s %s\r\n", resp.Method, resp.Path)
	for _, header := range resp.Headers {
		fmt.Fprintf(buf, "%s: %s\r\n", header.Name, header.Value)
	}
	fmt.Fprintf(buf, "\r\n%s", resp.Body)
	return buf.String()
}

// AsHTTP returns a http.Response representation.
func (resp Response) AsHTTP() *http.Response {
	headers := make(map[string][]string)
	for _, header := range resp.Headers {
		h, ok := headers[header.Name]
		if !ok {
			h = []string{}
		}
		h = append(h, header.Value)
		headers[header.Name] = h
	}
	return &http.Response{
		ProtoMajor: 1,
		ProtoMinor: 1,

		StatusCode: resp.Status,
		Header:     headers,
		Body:       ioutil.NopCloser(strings.NewReader(resp.Body)),
	}
}

func asResponse(req *http.Request, resp *http.Response) Response {
	headers := make([]Header, 0)
	for name, vals := range resp.Header {
		for _, val := range vals {
			headers = append(headers, Header{Name: name, Value: val})
		}
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)
	return Response{
		Method:  req.Method,
		Path:    req.URL.Path,
		Status:  resp.StatusCode,
		Headers: headers,
		Body:    buf.String(),
	}
}

// Header is a single-valued HTTP header name and value
type Header struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// JSONResponse creates a Response with "Content-Type: application/json".
func JSONResponse(method, path, body string) Response {
	return Response{
		Method:  method,
		Path:    path,
		Status:  200,
		Headers: []Header{Header{Name: "Content-Type", Value: "application/json"}},
		Body:    body,
	}
}

func readResponse(form url.Values) Response {
	r := Response{}
	r.Method = form.Get("method")
	r.Path = form.Get("path")
	r.Status = 200
	headers := make([]Header, 0)
	for i, name := range form["header"] {
		headers = append(headers, Header{Name: name, Value: form["value"][i]})
	}
	r.Body = form.Get("body")
	return r
}

func loadResponses(path string) ([]Response, error) {
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(out, &responses)
	if err != nil {
		return nil, err
	}

	return responses, nil
}

var stubTmpl = template.Must(template.New("").Parse(`<!doctype html>
<html>
	<head>
	</head>

	<body>
		<form method="POST" action="/_stub">
			<input type="text" name="method" placeholder="GET" />
			<input type="text" name="path" placeholder="/request/path?query" />
			<ul>
				<li>
					<input type="text" name="header" placeholder="Content-Type" />
					<input type="text" name="value" placeholder="application/json" />
				</li>
			</ul>
			<textarea name="body" placeholder="{}"></textarea>
			<input type="submit" />
		</form>
	</body>
</html>`))
