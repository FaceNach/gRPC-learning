package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"rest_api_go/internal/models"
	"rest_api_go/pkg/utils"
	"strconv"
	"strings"
)

func GetTeachersDBHandler(r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return nil, utils.ErrorHandler(err, "error retrieving data")
	}
	defer db.Close()

	query := "SELECT * FROM teachers WHERE 1=1"
	var args []interface{}

	query, args = addFiltersToGetTeachersQuery(r, query, args)
	query = addFilterSortBy(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)

		return nil, utils.ErrorHandler(err, "error retrieving data")
	}

	defer rows.Close()

	teachersList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher

		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

		if err != nil {
			fmt.Println(err)
			return nil, utils.ErrorHandler(err, "error retrieving data")
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

func GetOneTeacherDBHandler(id int) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return models.Teacher{}, utils.ErrorHandler(err, "error retrieving data")
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, fmt.Sprintf("no teacher found with id: %v", id))
	}

	if err != nil {
		fmt.Println(err)
		return models.Teacher{}, utils.ErrorHandler(err, "error retrieving data")
	}

	return teacher, nil
}

func AddedTeachersDBHandler(newTeachers[]models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Printf("error: %v", err)
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer db.Close()

	
	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		result, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			fmt.Println(err)
			return nil, utils.ErrorHandler(err, "error adding data")
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(err, "error adding data")
		}

		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}

	return addedTeachers, nil
}

func PatchTeachersDBHandler(updates []map[string]interface{} ,r *http.Request) error {
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return utils.ErrorHandler(err, "error adding data")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("error with db: %v", err)
		return utils.ErrorHandler(err, "error updating data")
	}
	defer db.Close()

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			err := tx.Rollback()
			log.Printf("error : %v", err)
			return utils.ErrorHandler(err, "error updating data")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("error: %v", err)
			return utils.ErrorHandler(err, "error updating data")
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
			return utils.ErrorHandler(err, "teacher not found")
		}

		if err != nil {
			log.Printf("error: %v", err)
			return utils.ErrorHandler(err, "invalid id")
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
							return utils.ErrorHandler(err, "error updating data")
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
			return utils.ErrorHandler(err, "error updating data")
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("error: %v", err)
		return utils.ErrorHandler(err, "error updating data")
	}

	return nil
}

func PatchOneTeacherDBHandler(updates map[string]any, id int) (models.Teacher, error) {
	
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
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
		log.Printf("error: %v", err)
		return models.Teacher{}, utils.ErrorHandler(err, "error updating data")
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
		return models.Teacher{}, utils.ErrorHandler(err, "error updating to the DB")
	}

	return existingTeacher, nil
}

func DeleteOneTeacherDBHandler( id int) (int, error) {

	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return 0, utils.ErrorHandler(err, "internal server error")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		log.Printf("error deleting teacher of the DB: %v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error retreaving result of deleting teacher of the DB: %v", err)
		return 0, utils.ErrorHandler(err, "internal server error")
	}

	if rowsAffected == 0 {
		log.Printf("no teacher found with the id %v on the DB", id)
		return 0, utils.ErrorHandler(err, fmt.Sprintf("no teacher found with id %v", id))
	}

	return id, nil
}

func DeleteTeachersDBHandlers(idTeachersToDelete []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Printf("error : %v", err)
		return nil, utils.ErrorHandler(err, "internal server error")
	}
	defer db.Close()

	

	tx, err := db.Begin()
	if err != nil {
		log.Printf("error: %v", err)
		return nil, utils.ErrorHandler(err, "internal server error")
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		log.Printf("error: %v", err)
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "error deleting teacher")
	}

	defer stmt.Close()

	var deletedId []int

	for _, id := range idTeachersToDelete {

		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(err, "error deleting teacher")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(err, "error deleting teacher")
		}

		if rowsAffected <= 0 {
			log.Printf("no teacher found with id: %v", id)
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("no teacher found with ID: %v", id))
		}

		deletedId = append(deletedId, id)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("error:%v", err)
		return nil, utils.ErrorHandler(err, "error deleting teacher")
	}

	if len(deletedId) <= 0 {
		log.Printf("error: %v", err)
		return nil, utils.ErrorHandler(err, fmt.Sprintf("no teacher found with ID: %v", deletedId))
	}

	return deletedId, nil
}
