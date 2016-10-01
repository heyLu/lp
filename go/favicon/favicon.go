package main

import (
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
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
var debug = flag.Bool("debug", false, "Print out debug info")

var faviconCache *lru.Cache
var lock sync.RWMutex

var imageCache *lru.Cache
var imageHashes *lru.Cache
var imageLock sync.RWMutex

func main() {
	flag.Parse()

	faviconCache = lru.New(*cacheSize)
	imageCache = lru.New(*cacheSize)
	imageHashes = lru.New(*cacheSize)

	http.HandleFunc("/favicon", HandleGetFavicon)
	http.HandleFunc("/favicon_proxy", HandleProxy)
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

func HandleProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=2419200")

	url := r.URL.Query()["url"][0]
	favicon, err := GetFaviconCached(url)
	if err != nil {
		fmt.Printf("Error: '%s': %s\n", url, err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}

	image, hash, err := GetImageCached(favicon)
	if err != nil {
		fmt.Printf("Error: '%s': %s\n", url, err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}
	w.Header().Set("ETag", hash)

	ifNoneMatch := r.Header.Get("If-None-Match")
	if ifNoneMatch != "" && hash == ifNoneMatch {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(image)))
	w.Write(image)
}

func GetImageCached(u string) ([]byte, string, error) {
	imageLock.RLock()
	image, cached := imageCache.Get(u)
	hash, _ := imageHashes.Get(u)
	imageLock.RUnlock()

	if cached {
		return image.([]byte), hash.(string), nil
	}

	resp, err := http.Get(u)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	imageHash := fmt.Sprintf("%x", md5.Sum(imageData))

	imageLock.Lock()
	imageCache.Add(u, imageData)
	imageHashes.Add(u, imageHash)
	imageLock.Unlock()
	return imageData, imageHash, nil

}

func HandleGetFavicon(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()["url"][0]
	favicon, err := GetFaviconCached(url)
	if err != nil {
		fmt.Printf("Error: '%s': %s\n", url, err)
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
		switch fu.(type) {
		case string:
			return fu.(string), nil
		case error:
			return "", fu.(error)
		default:
			panic("unexpected type")
		}
	}

	faviconUrl, err := GetFavicon(u)

	lock.Lock()
	if err != nil {
		faviconCache.Add(host, err)
	} else {
		faviconCache.Add(host, faviconUrl)
	}
	lock.Unlock()
	return faviconUrl, err
}

func GetFavicon(url string) (string, error) {
	if favicon, err := GetCanonicalFavicon(url); err == nil {
		fmt.Println("found favicon.ico")
		return favicon, nil
	} else if *debug {
		fmt.Printf("Error: getting /favicon.ico: %s\n", err)
	}

	resp, err := http.Get(url)
	if *debug {
		fmt.Println("get html", resp, err)
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tree, err := html.Parse(resp.Body)
	if *debug {
		fmt.Println("parse html", tree, err)
	}
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
	if *debug {
		fmt.Println("get favicon.ico", resp, err)
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 || resp.Header.Get("Content-Length") == "0" {
		return "", errors.New("no /favicon.ico")
	}
	buf := make([]byte, 0, 1)
	n, err := resp.Body.Read(buf)
	if err != nil || n == 0 {
		return "", errors.New("can't read /favicon.ico")
	}
	if *debug {
		fmt.Println("favicon.ico", resp.Request.URL.String(), faviconUrl)
	}
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
