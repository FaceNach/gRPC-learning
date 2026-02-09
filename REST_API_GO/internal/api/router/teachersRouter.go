package router

import (
	"net/http"
	"rest_api_go/internal/api/handlers"
)

func teachersRouter() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /{$}", handlers.AddTeacherHandler)
	mux.HandleFunc("DELETE /{$}", handlers.DeleteTeachersHandler)
	mux.HandleFunc("PATCH /{$}", handlers.PatchTeachersHandler)

	mux.HandleFunc("GET /{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("PUT /{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /{id}", handlers.PatchOneTeacherHandler)
	mux.HandleFunc("DELETE /{id}", handlers.DeleteOneTeacherHandler)

	mux.HandleFunc("GET /{id}/students", handlers.GetStudentsByTeacherId)
	mux.HandleFunc("GET /{id}/studentsCount", handlers.GetStudentsCountByTeacherId)

	return mux
}
