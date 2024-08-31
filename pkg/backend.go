package bankend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ServerConfig struct {
	rp  *httputil.ReverseProxy
	url *url.URL
}

func NewBackendConfig(BackendUrl string) (*ServerConfig, error) {
	rp := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = ":8081"
		},
	}
	url, err := url.Parse(BackendUrl)
	if err != nil {
		return nil, err
	}
	return &ServerConfig{
		rp:  &rp,
		url: url,
	}, nil
}
