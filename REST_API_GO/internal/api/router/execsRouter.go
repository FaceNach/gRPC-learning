package router

import (
	"net/http"
	"rest_api_go/internal/api/handlers"
)

func execsRouter() *http.ServeMux{
	 mux := http.NewServeMux()
		
	mux.HandleFunc("GET /{$}", handlers.ExecsHandler)
	mux.HandleFunc("POST /{$}", handlers.ExecsHandler)
	mux.HandleFunc("PATCH /{$}", handlers.ExecsHandler)
	
	mux.HandleFunc("GET /{id}", handlers.ExecsHandler)
	mux.HandleFunc("PATCH /{id}", handlers.ExecsHandler)
	mux.HandleFunc("DELETE /{id}", handlers.ExecsHandler)
	mux.HandleFunc("POST /{id}/updatepassword", handlers.ExecsHandler)
	
	mux.HandleFunc("POST /login", handlers.ExecsHandler)
	mux.HandleFunc("POST /logout", handlers.ExecsHandler)
	mux.HandleFunc("POST /forgotpassword", handlers.ExecsHandler)
	mux.HandleFunc("POST /resetpassword/reset/{resetcode}", handlers.ExecsHandler)
	
	return mux
}