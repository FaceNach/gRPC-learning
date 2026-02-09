package router

import (
	"net/http"
	"rest_api_go/internal/api/handlers"
)

func execsRouter() *http.ServeMux{
	 mux := http.NewServeMux()
		
	mux.HandleFunc("GET /{$}", handlers.GetExecsHandler)
	mux.HandleFunc("POST /{$}", handlers.AddExecsHandler)
	mux.HandleFunc("PATCH /{$}", handlers.PatchExecsHandler)
	
	mux.HandleFunc("GET /{id}", handlers.GetOneExecHandler)
	mux.HandleFunc("PATCH /{id}", handlers.PatchOneExecHandler)
	mux.HandleFunc("DELETE /{id}", handlers.DeleteOneExecHandler)
	mux.HandleFunc("POST /{id}/updatepassword", handlers.GetExecsHandler)
	
	mux.HandleFunc("POST /login", handlers.GetExecsHandler)
	mux.HandleFunc("POST /logout", handlers.GetExecsHandler)
	mux.HandleFunc("POST /forgotpassword", handlers.GetExecsHandler)
	mux.HandleFunc("POST /resetpassword/reset/{resetcode}", handlers.GetExecsHandler)
	
	return mux
}