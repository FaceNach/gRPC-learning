package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest_api_go/internal/models"
	"rest_api_go/internal/repository/sqlconnect"
	"strconv"
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	teachersList, err := sqlconnect.GetTeachersDBHandler(w, r)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachersList),
		Data:   teachersList,
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error getting all the teachers data", http.StatusBadRequest)
		return
	}

}

func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

	teacher, err := sqlconnect.GetOneTeacherDBHandler(w, r)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(teacher)
	if err != nil {
		http.Error(w, "Error trying to access only one teacher", http.StatusBadRequest)
	}
}



func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {

	addedTeachers, err := sqlconnect.AddedTeachersDBHandler(w, r)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error posting the teacher", http.StatusBadRequest)
	}
}



func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	err := sqlconnect.PatchTeachersDBHandler(w,r)
	if err != nil {
		log.Printf("error: %v",err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// PUT /teachers/
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var updatedTeacher models.Teacher

	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		log.Printf("internal error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error connecting to the DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject)

	if err == sql.ErrNoRows {
		log.Printf("no match found: %v", err)
		http.Error(w, "no match found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName,
		updatedTeacher.LastName,
		updatedTeacher.Email,
		updatedTeacher.Class,
		updatedTeacher.Subject,
		updatedTeacher.ID)

	if err != nil {
		log.Printf("erro updating teacher in the DB: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedTeacher)
	if err != nil {
		log.Printf("error encoding the response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

//PATCH /teachers/{id}

func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

	updatedTeacher, err := sqlconnect.PatchOneTeacherDBHandler(w,r)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedTeacher)
	if err != nil {
		log.Printf("error encoding the response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

	id, err := sqlconnect.DeleteOneTeacherDBHandler(w,r)
	
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "teacher deleted successfully from the DB",
		ID:     id,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error encoding the response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {

	deletedIds, err := sqlconnect.DeleteTeachersDBHandlers(w,r)

	w.Header().Set("Content-type", "application/json")
	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
		Data   []int  `json:"deletedID"`
	}{
		Status: "correct",
		Count:  len(deletedIds),
		Data:   deletedIds,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error encoding json: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
