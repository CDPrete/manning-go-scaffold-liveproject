package handlers

import "net/http"

func Register(mux *http.ServeMux)  {
	mux.HandleFunc("/", HandleIndex)
	mux.HandleFunc("/healthcheck", HandleHealthCheck)
	mux.HandleFunc("/static/", HandleStatic)
}
