package handlers

import (
	"io/ioutil"
	"net/http"
)

// TestHandler will handle the payload from any services and log the payload.
func TestHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	content, err := ioutil.ReadAll(r.Body)

	if err != nil {
		write(w, err.Error())
	} else {
		write(w, string(content))
	}
}
