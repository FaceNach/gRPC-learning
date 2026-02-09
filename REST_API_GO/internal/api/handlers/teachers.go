package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"rest_api_go/internal/models"
	"rest_api_go/internal/repository/sqlconnect"
	"strconv"
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	teachersList, err := sqlconnect.GetTeachersDBHandler(r)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	teacher, err := sqlconnect.GetOneTeacherDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(teacher)
	if err != nil {
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}

	for _, teacher := range newTeachers {
		// if teacher.FirstName == ""  ||teacher.LastName == "" || teacher.Email == "" || teacher.Class == "" || teacher.Subject == "" {
		// 	http.Error(w, "All fields are obligatory", http.StatusBadRequest)
		// 	return
		// }

		val := reflect.ValueOf(teacher)
		for i := 0; i < val.Type().NumField(); i++ {
			field := val.Field(i)
			if field.Kind() == reflect.String && field.String() == "" {
				http.Error(w, "All fields are obligatory", http.StatusBadRequest)
				return
			}
		}

	}

	addedTeachers, err := sqlconnect.AddedTeachersDBHandler(newTeachers)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusBadRequest)
		return
	}
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

func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}

	err = sqlconnect.PatchTeachersDBHandler(updates, r)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//PATCH /teachers/{id}

func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

	var updates map[string]any

	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing the data", http.StatusInternalServerError)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	updatedTeacher, err := sqlconnect.PatchOneTeacherDBHandler(updates, id)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	idDeleted, err := sqlconnect.DeleteOneTeacherDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "teacher deleted successfully from the DB",
		ID:     idDeleted,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error encoding the response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var idTeachersToDelete []int

	err := json.NewDecoder(r.Body).Decode(&idTeachersToDelete)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachersDBHandlers(idTeachersToDelete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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


func GetStudentsByTeacherId(w http.ResponseWriter, r *http.Request){
	teacherId := r.PathValue("id")
	
	students, err := sqlconnect.GetStudentsByTeacherIdDBHandler(teacherId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := struct {
		Status string `json:"status"`
		Count int `json:"count"`
		Data []models.Student `json:"data"`
	}{
		Status : "success",
		Count : len(students),
		Data : students,
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing the response", http.StatusInternalServerError)
		return
	}
}

func GetStudentsCountByTeacherId(w http.ResponseWriter, r *http.Request){
	
	teacherId := r.PathValue("id")
	
	studentsCount, err := sqlconnect.GetStudentsCountByTeacherIdDBHandler(teacherId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := struct {
		Status string `json:"status"`
		Count int `json:"count"`
	}{
		Status : "success",
		Count : studentsCount,
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing the response", http.StatusInternalServerError)
		return
	}
	
}