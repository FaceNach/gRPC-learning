package sqlconnect

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"rest_api_go/internal/models"
	"strconv"
	"strings"
)

func GetTeachersDBHandler(w http.ResponseWriter, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	query := "SELECT * FROM teachers WHERE 1=1"
	var args []interface{}

	query, args = addFiltersToGetTeachersQuery(r, query, args)
	query = addFilterSortBy(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "database query error", http.StatusInternalServerError)
		return nil, err
	}

	defer rows.Close()

	teachersList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher

		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "error scanning database results", http.StatusInternalServerError)
			return nil, err
		}

		teachersList = append(teachersList, teacher)
	}

	return teachersList, nil
}

func addFiltersToGetTeachersQuery(r *http.Request, query string, args []any) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}

	return query, args
}

func addFilterSortBy(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]

	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}

			field, order := parts[0], parts[1]

			if !isValidSortOrder(order) || !isValidSortField(field) {
				continue
			}

			if i > 0 {
				query += ","
			}

			query += " " + field + " " + order
		}
	}

	return query
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}

	exists := validFields[field]

	return exists
}

func GetOneTeacherDBHandler(w http.ResponseWriter, r *http.Request) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	defer db.Close()

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "internal server: error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	var teacher models.Teacher
	err = db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return models.Teacher{}, err
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error finding the teacher", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	return teacher, nil
}

func AddedTeachersDBHandler(w http.ResponseWriter, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return nil, err
	}

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error in preparing query", http.StatusInternalServerError)
		return nil, err
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		result, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error inserting data to DB", http.StatusInternalServerError)
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting the ID of the last inserted item in the DB", http.StatusInternalServerError)
			return nil, err
		}

		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}

	return addedTeachers, nil
}

func PatchTeachersDBHandler(w http.ResponseWriter, r *http.Request) error {
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error connecting to the DB", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	var updates []map[string]interface{}

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Printf("error encoding data: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("error with db: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			err := tx.Rollback()
			log.Printf("error converting var: %v", err)
			http.Error(w, "invalid teacher ID in update", http.StatusInternalServerError)
			return err
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("error: %v", err)
			http.Error(w, "error trying to convert from string to int", http.StatusInternalServerError)
			return err
		}

		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(&teacherFromDb.ID,
			&teacherFromDb.FirstName,
			&teacherFromDb.LastName,
			&teacherFromDb.Email,
			&teacherFromDb.Class,
			&teacherFromDb.Subject)

		if err == sql.ErrNoRows {
			tx.Rollback()
			log.Printf("error: %v", err)
			http.Error(w, "teacher not found", http.StatusNotFound)
			return err
		}

		if err != nil {
			log.Printf("error: %v", err)
			http.Error(w, "error retrieving teacher", http.StatusInternalServerError)
			return err
		}

		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}

			for i := 0; i < teacherVal.NumField(); i++ {

				field := teacherType.Field(i)

				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							err = tx.Rollback()
							log.Printf("error: %v", err)
							return err
						}
					}
					break
				}
			}

		}

		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, subject = ? WHERE id = ?",
			teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Subject, teacherFromDb.ID)

		if err != nil {
			log.Printf("error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "error commiting transaction", http.StatusInternalServerError)
		return err
	}

	return nil
}

func PatchOneTeacherDBHandler(w http.ResponseWriter, r *http.Request) (models.Teacher, error) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	var updates map[string]any

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Printf("internal error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error connecting to the DB", http.StatusInternalServerError)
		return models.Teacher{}, err
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

	if err != nil {
		log.Printf("error updating teacher to the DB: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)

			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		existingTeacher.FirstName,
		existingTeacher.LastName,
		existingTeacher.Email,
		existingTeacher.Class,
		existingTeacher.Subject,
		existingTeacher.ID)

	if err != nil {
		log.Printf("erro updating teacher in the DB: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	return existingTeacher, nil
}

func DeleteOneTeacherDBHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return 0, err
	}

	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error connecting to the DB", http.StatusInternalServerError)
		return 0, err
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		log.Printf("error deleting teacher of the DB: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error retreaving result of deleting teacher of the DB: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return 0, err
	}

	if rowsAffected == 0 {
		log.Printf("no teacher found with that id on the DB: %v ", err)
		http.Error(w, "internal server error", http.StatusNotFound)
		return 0, err
	}

	return id, nil
}

func DeleteTeachersDBHandlers(w http.ResponseWriter, r *http.Request)([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Printf("error connecting to the DB: %v", err)
		http.Error(w, "error connecting to the DB", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	var idTeachersToDelete []int

	err = json.NewDecoder(r.Body).Decode(&idTeachersToDelete)
	if err != nil {
		log.Printf("error decoding json: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return nil, err
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		log.Printf("error: %v", err)
		tx.Rollback()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return nil, err
	}

	defer stmt.Close()

	var deletedId []int

	for _, id := range idTeachersToDelete {

		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			log.Printf("error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return nil, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			log.Printf("error:%v", err)
			http.Error(w, "error retreiving delete result", http.StatusInternalServerError)
			return nil, err
		}

		if rowsAffected <= 0 {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("no teacher found with the id %v", id), http.StatusNotFound)
			return nil, err
		}

		deletedId = append(deletedId, id)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("error:%v", err)
		http.Error(w, "error trying to delete teachers from DB", http.StatusInternalServerError)
		return nil, err
	}

	if len(deletedId) <= 0 {
		log.Printf("error: %v", err)
		http.Error(w, "No teachers found with that id", http.StatusNotFound)
		return nil, err
	}
	
	return deletedId, nil
}
