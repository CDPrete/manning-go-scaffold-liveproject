package handlers

import (
	"fmt"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, _ *http.Request)  {
	_, _ = fmt.Fprint(w, "Hello World")
}
