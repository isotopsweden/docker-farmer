package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func write(w http.ResponseWriter, text string) {
	log.Println(text)
	fmt.Fprintf(w, text)
}
