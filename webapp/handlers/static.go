package handlers

import "net/http"

func HandleStatic(w http.ResponseWriter, req *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, req)
}
