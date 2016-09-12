package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// stringInSlice returns true if a string exists or false if not.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}

// write writes to both log and http response writer.
func write(w http.ResponseWriter, text string) {
	log.Println(text)
	fmt.Fprintf(w, text)
}
