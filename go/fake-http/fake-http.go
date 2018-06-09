package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

var flags struct {
	addr            string
	proxyURL        string
	proxyClientCert string
	proxyClientKey  string

	proxyMinikube bool
}

func init() {
	flag.StringVar(&flags.addr, "addr", "localhost:8080", "Address to listen on")
	flag.StringVar(&flags.proxyURL, "proxy-url", "", "Proxy requests to this URL")
	flag.StringVar(&flags.proxyClientCert, "proxy-client-cert", "", "Client certificate to use when connecting to proxy")
	flag.StringVar(&flags.proxyClientKey, "proxy-client-key", "", "Client key to use when connecting to proxy")

	flag.BoolVar(&flags.proxyMinikube, "proxy-minikube", false, "Shortcut for -proxy-url https://$(minikube ip):8443 -proxy-client-cert ~/.minikube/client.crt -proxy-client-key ~/.minikube/client.key")
}

var responses Responses

// Responses is a list of responses that will be stubbed/faked.
type Responses []Response

func main() {
	flag.Parse()

	if flags.proxyMinikube {
		err := proxyMinikube()
		if err != nil {
			log.Fatalf("Error: Setting up Minikube proxy: %s", err)
		}
	}

	var responsesPath string
	if flag.NArg() == 1 {
		responsesPath = flag.Arg(0)
	}

	requestLog := Log(make([]LogEntry, 0))

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		responses.Load(responsesPath)

		var resp *http.Response
		if flags.proxyURL != "" {
			resp = respondWithProxy(flags.proxyURL, w, req)
		} else {
			resp = respondWithStub(responses, w, req)
		}

		userAgent := req.Header.Get("User-Agent")
		log.Printf("%s %s - %d (%s, %q)", req.Method, req.URL, resp.StatusCode, req.RemoteAddr, userAgent)

		if resp.Header.Get("Content-Type") == "application/json" {
			pretty, err := prettyfyJSON(resp.Body)
			if err != nil {
				log.Printf("Error: Prettyfying JSON: %s", err)
			} else {
				resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(pretty)))
				resp.Body = ioutil.NopCloser(bytes.NewReader(pretty))
			}
		}

		requestLog.Log(req, resp)
	})

	http.HandleFunc("/_log", func(w http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.Header.Get("Accept"), "application/yaml") {
			rs := make([]Response, len(requestLog))
			for i, log := range requestLog {
				rs[i] = log.AsResponse()
			}
			err := renderYAML(w, rs)
			if err != nil {
				log.Printf("Error: Render log: %s", err)
			}
			return
		}

		for i := len(requestLog) - 1; i >= 0; i-- {
			w.Write([]byte("------------------------------------------------------\n"))
			requestLog[i].Request().Write(w)
			w.Write([]byte("\n"))
			requestLog[i].Response().Write(w)
			w.Write([]byte("\n"))
		}
	})

	http.HandleFunc("/_help", func(w http.ResponseWriter, req *http.Request) {
		urls := []struct {
			URL     string
			Summary string
		}{
			{URL: "/_log", Summary: "View all received requests with responses"},
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

func proxyMinikube() error {
	out, err := exec.Command("minikube", "ip").Output()
	if err != nil {
		return err
	}
	flags.proxyURL = fmt.Sprintf("https://%s:8443", strings.TrimSpace(string(out)))

	homeDir := os.Getenv("HOME")
	flags.proxyClientCert = path.Join(homeDir, ".minikube/client.crt")
	flags.proxyClientKey = path.Join(homeDir, ".minikube/client.key")

	return nil
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

// Log collects HTTP requests/responses for later display and
// processing.
type Log []LogEntry

// Log logs the request/response pair.
func (l *Log) Log(req *http.Request, resp *http.Response) {
	e := LogEntry{
		request:      req,
		requestBody:  new(bytes.Buffer),
		response:     resp,
		responseBody: new(bytes.Buffer),
	}
	io.Copy(e.requestBody, req.Body)
	io.Copy(e.responseBody, resp.Body)
	*l = append(*l, e)
}

// LogEntry is a request/response pair for logging.
type LogEntry struct {
	request      *http.Request
	requestBody  *bytes.Buffer
	response     *http.Response
	responseBody *bytes.Buffer
}

// AsResponse returns a Response representation of the entry.
func (e LogEntry) AsResponse() Response {
	return Response{
		Method: e.request.Method,
		Path:   e.request.URL.Path,
		Status: e.response.StatusCode,
		Body:   e.responseBody.String(),
	}
}

// Request returns the stored http.Request.
func (e LogEntry) Request() *http.Request {
	e.request.Body = ioutil.NopCloser(bytes.NewReader(e.requestBody.Bytes()))
	return e.request
}

// Response returns the stored http.Response.
func (e LogEntry) Response() *http.Response {
	e.response.Body = ioutil.NopCloser(bytes.NewReader(e.responseBody.Bytes()))
	return e.response
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

// Load loads responses from the YAML file at path.
func (rs *Responses) Load(path string) {
	if path == "" {
		return
	}
	responses, err := rs.loadFile(path)
	if err != nil {
		log.Printf("Error: Parsing %s: %s", path, err)
	}
	*rs = responses
}

func (rs *Responses) loadFile(path string) ([]Response, error) {
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
