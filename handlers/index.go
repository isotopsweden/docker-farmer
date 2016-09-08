package handlers

import (
	"fmt"
	"net/http"
)

// IndexHandler will render the index route.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Docker Farmer!")
}
