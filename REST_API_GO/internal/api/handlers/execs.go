package handlers

import "net/http"

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET method on execs route"))
	case http.MethodPost:
		w.Write([]byte("Hello POST method on execs route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT method on execs route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE method on execs route"))
	}
}
