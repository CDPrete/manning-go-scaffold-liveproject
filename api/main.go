package main

import (
	"log"
	"main/handlers"
	"net/http"
)

func main() {
	var mux = http.NewServeMux()
	handlers.Register(mux)

	if err := http.ListenAndServe("", mux); err != nil {
		log.Fatalf("Cannot provide API server: \n%s", err)
	}
}
