package main

import (
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var faviconCache = make(map[string]string)

func main() {
	http.HandleFunc("/favicon", HandleGetFavicon)
	err := http.ListenAndServe(":8080", nil)

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
	faviconUrl, cached := faviconCache[host]

	if cached {
		return faviconUrl, nil
	}

	faviconUrl, err = GetFavicon(u)
	if err != nil {
		return faviconUrl, err
	}

	faviconCache[host] = faviconUrl
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
