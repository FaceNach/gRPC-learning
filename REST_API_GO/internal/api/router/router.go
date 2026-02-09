package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {

	mux := http.NewServeMux()

	eRouter := execsRouter()
	tRouter := teachersRouter()
	sRouter := studentsRouter()

	mux.Handle("/teachers/", http.StripPrefix("/teachers", tRouter))
	mux.Handle("/students/", http.StripPrefix("/students", sRouter))
	mux.Handle("/execs/", http.StripPrefix("/execs", eRouter))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to our scholar API"))
	})

	return mux
}
