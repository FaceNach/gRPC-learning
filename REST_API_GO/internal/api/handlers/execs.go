package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"rest_api_go/internal/models"
	"rest_api_go/internal/repository/sqlconnect"
	"rest_api_go/pkg/utils"
	"strconv"
	"time"
)

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {

	execList, err := sqlconnect.GetExecsDBHandler(r)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(execList),
		Data:   execList,
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error getting all the execs data", http.StatusBadRequest)
		return
	}

}

func GetOneExecHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	exec, err := sqlconnect.GetOneExecDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(exec)
	if err != nil {
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}
}

func AddExecsHandler(w http.ResponseWriter, r *http.Request) {

	var newExecs []models.Exec
	err := json.NewDecoder(r.Body).Decode(&newExecs)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}

	for _, exec := range newExecs {
		// if teacher.FirstName == ""  ||teacher.LastName == "" || teacher.Email == "" || teacher.Class == "" || teacher.Subject == "" {
		// 	http.Error(w, "All fields are obligatory", http.StatusBadRequest)
		// 	return
		// }

		val := reflect.ValueOf(exec)
		for i := 0; i < val.Type().NumField(); i++ {
			field := val.Field(i)
			if field.Kind() == reflect.String && field.String() == "" {
				http.Error(w, "All fields are obligatory", http.StatusBadRequest)
				return
			}
		}

	}

	addedExecs, err := sqlconnect.AddExecsDBHandler(newExecs)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusBadRequest)
		return
	}
}

func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error parsing data", http.StatusInternalServerError)
		return
	}

	err = sqlconnect.PatchExecsDBHandler(updates)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//PATCH /teachers/{id}

func PatchOneExecHandler(w http.ResponseWriter, r *http.Request) {

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

	updatedExec, err := sqlconnect.PatchOneExecDBHandler(updates, id)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedExec)
	if err != nil {
		log.Printf("error encoding the response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func DeleteOneExecHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	idDeleted, err := sqlconnect.DeleteOneExecDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "exec deleted successfully from the DB",
		ID:     idDeleted,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error encoding the response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var exec models.Exec

	err := json.NewDecoder(r.Body).Decode(&exec)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := sqlconnect.LoginDBHandler(exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := strconv.Itoa(user.ID)
	tokenString, err := utils.SignToken(userID, user.Username, user.Role)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "test",
		Value:    "testing",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func LogoutHandler (w http.ResponseWriter, r *http.Request){
	http.SetCookie(w, &http.Cookie{
		Name: "Bearer",
		Value:"",
		Path: "/",
		HttpOnly: true,
		Secure: true,
		Expires: time.Unix(0,0),
		SameSite: http.SameSiteStrictMode,
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message" : "logged out successfully"}`))
}