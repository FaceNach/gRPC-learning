package router

import (
	"net/http"
	"rest_api_go/internal/api/handlers"
)

func studentsRouter() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /{$}", handlers.AddStudentsHandler)
	mux.HandleFunc("DELETE /{$}", handlers.DeleteStudentsHandler)
	mux.HandleFunc("PATCH /{$}", handlers.PatchStudentsHandler)

	mux.HandleFunc("GET /{id}", handlers.GetOneStudentHandler)
	mux.HandleFunc("PUT /{id}", handlers.UpdateStudentHandler)
	mux.HandleFunc("PATCH /{id}", handlers.PatchOneStudentHandler)
	mux.HandleFunc("DELETE /{id}", handlers.DeleteOneStudentHandler)

	return mux
}
