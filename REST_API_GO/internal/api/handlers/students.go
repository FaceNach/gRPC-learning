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

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {

	studentList, err := sqlconnect.GetStudentsDBHandler(r)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(studentList),
		Data:   studentList,
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error getting all the students data", http.StatusBadRequest)
		return
	}

}

func GetOneStudentHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	student, err := sqlconnect.GetOneStudentDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}
}

func AddStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var newStudents []models.Student
	err := json.NewDecoder(r.Body).Decode(&newStudents)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}

	for _, student := range newStudents {
		// if teacher.FirstName == ""  ||teacher.LastName == "" || teacher.Email == "" || teacher.Class == "" || teacher.Subject == "" {
		// 	http.Error(w, "All fields are obligatory", http.StatusBadRequest)
		// 	return
		// }

		val := reflect.ValueOf(student)
		for i := 0; i < val.Type().NumField(); i++ {
			field := val.Field(i)
			if field.Kind() == reflect.String && field.String() == "" {
				http.Error(w, "All fields are obligatory", http.StatusBadRequest)
				return
			}
		}

	}

	addedStudents, err := sqlconnect.AddedStudentsDBHandler(newStudents)
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
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusBadRequest)
		return
	}
}

// PUT /teachers/
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var updateStudent models.Student

	err = json.NewDecoder(r.Body).Decode(&updateStudent)
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

	var existingStudent models.Student
	err = db.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(
		&existingStudent.ID,
		&existingStudent.FirstName,
		&existingStudent.LastName,
		&existingStudent.Email,
		&existingStudent.Class)

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

	updateStudent.ID = existingStudent.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updateStudent.FirstName,
		updateStudent.LastName,
		updateStudent.Email,
		updateStudent.Class,
		updateStudent.ID)

	if err != nil {
		log.Printf("erro updating teacher in the DB: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updateStudent)
	if err != nil {
		log.Printf("error encoding the response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}

	err = sqlconnect.PatchStudentsDBHandler(updates, r)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//PATCH /teachers/{id}

func PatchOneStudentHandler(w http.ResponseWriter, r *http.Request) {

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

	updatedStudent, err := sqlconnect.PatchOneStudentDBHandler(updates, id)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedStudent)
	if err != nil {
		log.Printf("error encoding the response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func DeleteOneStudentHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	idDeleted, err := sqlconnect.DeleteOneStudentDBHandler(id)
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

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var idStudentsToDelete []int

	err := json.NewDecoder(r.Body).Decode(&idStudentsToDelete)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachersDBHandlers(idStudentsToDelete)
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
