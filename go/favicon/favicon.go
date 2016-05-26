package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"

	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"github.com/golang/groupcache/lru"
)

var port = flag.Int("p", 8080, "port [8080]")
var cacheSize = flag.Int("s", 10000, "cache size [10000]")

var faviconCache *lru.Cache
var lock sync.RWMutex

func main() {
	flag.Parse()

	faviconCache = lru.New(*cacheSize)

	http.HandleFunc("/favicon", HandleGetFavicon)
	if p := os.Getenv("PORT"); p != "" {
		flag.Set("p", p)
	}

	addr := fmt.Sprintf("localhost:%d", *port)
	fmt.Printf("listening on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}

func HandleGetFavicon(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()["url"][0]
	favicon, err := GetFaviconCached(url)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}
	_, noRedirect := r.URL.Query()["no_redirect"]
	if noRedirect {
		w.Write([]byte(favicon))
		return
	}
	w.Header().Set("Location", favicon)
	w.WriteHeader(http.StatusSeeOther)
}

func GetFaviconCached(u string) (string, error) {
	parsed, err := url.Parse(u)
	var host = ""
	if err != nil {
		host = u
	} else {
		host = parsed.Host
	}
	lock.RLock()
	fu, cached := faviconCache.Get(host)
	lock.RUnlock()

	if cached {
		return fu.(string), nil
	}

	faviconUrl, err := GetFavicon(u)
	if err != nil {
		return faviconUrl, err
	}

	lock.Lock()
	faviconCache.Add(host, faviconUrl)
	lock.Unlock()
	return faviconUrl, nil
}

func GetFavicon(url string) (string, error) {
	if favicon, err := GetCanonicalFavicon(url); err == nil {
		fmt.Println("found favicon.ico")
		return favicon, nil
	}

	resp, err := http.Get(url)
	fmt.Println("get html", resp, err)
	if err != nil {
		return "", err
	}
	tree, err := html.Parse(resp.Body)
	fmt.Println("parse html", tree, err)
	if err != nil {
		return "", err
	}

	sel := cascadia.MustCompile("link[rel~=icon]")
	node := sel.MatchFirst(tree)
	if node == nil {
		return "", errors.New("no favicon found")
	}

	favicon, found := FindAttr("href", node.Attr)
	if !found {
		return "", errors.New("no link found")
	}

	return ToAbsolute(resp.Request.URL, favicon.Val), nil
}

func GetCanonicalFavicon(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	faviconUrl := fmt.Sprintf("%s://%s/favicon.ico", parsed.Scheme, parsed.Host)

	resp, err := http.Get(faviconUrl)
	fmt.Println("get favicon.ico", resp, err)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 || resp.Header.Get("Content-Length") == "0" {
		return "", errors.New("no /favicon.ico")
	}
	fmt.Println("favicon.ico", resp.Request.URL.String(), faviconUrl)
	return resp.Request.URL.String(), nil
}

func FindAttr(name string, attributes []html.Attribute) (*html.Attribute, bool) {
	for _, attr := range attributes {
		if attr.Key == name {
			return &attr, true
		}
	}
	return nil, false
}

func ToAbsolute(base *url.URL, rawUrl string) string {
	url, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}
	return base.ResolveReference(url).String()
}
