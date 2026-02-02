package handlers

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello root route \n")
	w.Write([]byte("Hello root route"))
	fmt.Println("Hello root route")
}
