package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

func Execute() {
	fmt.Println("hello")

}

func TestServe() {
	dst := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = ":8082"
	}

	rp := &httputil.ReverseProxy{
		Director: dst,
	}

	s := http.Server{
		Addr:    ":8081",
		Handler: rp,
	}

	g := gin.Default()

	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	if err := g.Run(":8082"); err != nil {
		log.Fatal(err)
	}
	if err := g.Run(":8083"); err != nil {
		log.Fatal(err)
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
