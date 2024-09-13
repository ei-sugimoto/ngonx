package proxy

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

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
		parsedURL, err := url.Parse(u.URL)
		if err != nil {
			return nil, err
		}
		proxyURLs = append(proxyURLs, *parsedURL)
	}

	proxyList := make([]*httputil.ReverseProxy, 0, len(proxyURLs))

	for _, u := range proxyURLs {
		proxy := httputil.NewSingleHostReverseProxy(&u)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("proxy error: %v", err)
			http.Error(w, "Bad Gateway: Unable to reach the backend server", http.StatusBadGateway)
		}
		proxy.Transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 2 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 2 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
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
			log.Printf("proxy to %s for request %s", urlList[idx].URL, r.URL.Path)
			time.Sleep(1 * time.Second)
			proxy.ServeHTTP(w, r)
		})
	}

	return mux, nil
}
