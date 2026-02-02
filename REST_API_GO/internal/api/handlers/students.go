package handlers

import "net/http"

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET method on students route"))
	case http.MethodPost:
		w.Write([]byte("Hello POST method on students route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT method on students route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE method on students route"))
	}
}
