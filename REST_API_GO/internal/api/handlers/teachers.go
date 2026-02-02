package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest_api_go/internal/models"
	"strconv"
	"strings"
	"sync"
)

func TeachersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		addTeacherHandler(w, r)
	case http.MethodPut:
		w.Write([]byte("Hello PUT method on teachers route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE method on teachers route"))
	}
}

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextID   = 1
)

func init() {
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Jhon",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}

	nextID++

	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10A",
		Subject:   "Algebra",
	}
	nextID++

	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "11A",
		Subject:   "Bioligy",
	}
	nextID++
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	teachersSlice := make([]models.Teacher, 0, len(teachers))

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teachersSlice = append(teachersSlice, teacher)
			}
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teachersSlice),
			Data:   teachersSlice,
		}

		w.Header().Set("Content-type", "application/json")

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error getting all the teachers data", http.StatusBadRequest)
		}

		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error converting string to int")
		return
	}

	teacher, exists := teachers[id]

	if !exists {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	resOneTeacher := struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}{
		Status: "success",
		Data:   teacher,
	}

	err = json.NewEncoder(w).Encode(resOneTeacher)
	if err != nil {
		http.Error(w, "Error trying to access only one teacher", http.StatusBadRequest)
	}

}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	for i, teacher := range newTeachers {
		teacher.ID = nextID
		newTeachers[i].ID = nextID
		teachers[nextID] = teacher
		nextID++
	}

	w.Header().Set("Content-type", "applicatio/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(newTeachers),
		Data:   newTeachers,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error posting the teacher", http.StatusBadRequest)
	}
}
