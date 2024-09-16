package proxy

import (
	"log"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/ei-sugimoto/ngonx/pkg/parser"
	"github.com/gin-gonic/gin"
)

type ReverseProxy struct{}

func NewReverseProxy() *ReverseProxy {
	return &ReverseProxy{}
}

func (r *ReverseProxy) NewServe() {
	p := parser.NewServer()
	if err := p.Parse(); err != nil {
		panic(err)
	}

	g := gin.Default()

	for _, server := range p {
		g.GET(server.EndPoint, r.makeReserveProxy(server.Host, strconv.Itoa(server.Port), server.EndPoint))
	}

	if err := g.Run(":8080"); err != nil {
		panic(err)
	}
}

func (r *ReverseProxy) makeReserveProxy(host string, port string, endpoint string) func(*gin.Context) {
	return func(c *gin.Context) {
		remote, _ := url.Parse("http://" + host + ":" + port)

		log.Println("Proxy to", remote)
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, endpoint)
		proxy := httputil.NewSingleHostReverseProxy(remote)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
