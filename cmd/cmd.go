package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ei-sugimoto/ngonx/pkg/proxy"
)

func Run() {
	mux, err := proxy.NewServe()

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("error: %v", err)
	}
}
