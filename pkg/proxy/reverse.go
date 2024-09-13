package proxy

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ei-sugimoto/ngonx/pkg/parser"
)

func NewServe() (*http.ServeMux, error) {
	p := parser.NewServer()
	err := p.Parse()
	if err != nil {
		return nil, err
	}
	urlList := p.GetURLList()

	proxyURLs := make([]url.URL, 0, len(urlList))
	for _, u := range urlList {
		url, err := url.Parse(u.URL)
		if err != nil {
			return nil, err
		}
		proxyURLs = append(proxyURLs, *url)
	}

	proxyList := make([]*httputil.ReverseProxy, 0, len(proxyURLs))

	for _, u := range proxyURLs {
		proxy := httputil.NewSingleHostReverseProxy(&u)
		proxyList = append(proxyList, proxy)
	}

	mux := http.NewServeMux()

	if len(proxyList) == 0 {
		return nil, errors.New("no proxy server")
	}

	if len(proxyList) != len(urlList) {
		return nil, errors.New("proxy server count not match")
	}

	for idx, proxy := range proxyList {
		mux.HandleFunc(urlList[idx].EndPoint, func(w http.ResponseWriter, r *http.Request) {
			log.Println("proxy to", urlList[idx].URL)
			proxy.ServeHTTP(w, r)
		})
	}

	return mux, nil

}
