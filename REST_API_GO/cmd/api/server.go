package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	mw "rest_api_go/internal/api/middlewares"
)

func main() {
	port := 3000

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/teachers", teachersHandler)
	mux.HandleFunc("/students", studentsHandler)
	mux.HandleFunc("/execs", execsHandler)

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	//rl := mw.NewRateLimiter(5, time.Minute)
	
	// hppOptions := mw.HPPOptions{
	// 	CheckQuery: true,
	// 	CheckBody: true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	WhiteList: []string{"sortBy", "sortOrder", "name", "age", "class"},
		
	// }
	
	//secureMux := mw.Cors(rl.Middleware(mw.ResponseTime(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	//secureMux := applyMidlewares(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHeaders, mw.ResponseTime, rl.Middleware, mw.Cors)
	secureMux := mw.SecurityHeaders(mux)
	
	server := &http.Server{
		Addr:      ":3000",
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Listening server on port:", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error creating server: ", err)
	}
}

type Middleware func(http.Handler) http.Handler

func applyMidlewares (handler http.Handler, middlewares ...Middleware) http.Handler{
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	
	return handler
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello root route \n")
	w.Write([]byte("Hello root route"))
	fmt.Println("Hello root route")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET method on teachers route"))
	case http.MethodPost:
		w.Write([]byte("Hello POST method on teachers route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT method on teachers route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE method on teachers route"))
	}
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
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

func execsHandler(w http.ResponseWriter, r *http.Request) {
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
