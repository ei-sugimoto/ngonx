package cmd

import (
	"fmt"

	"github.com/ei-sugimoto/ngonx/pkg/proxy"
)

func Run() {
	mux := proxy.NewReverseProxy()
	mux.NewServe()

	fmt.Println("Server started on :8080")
}
