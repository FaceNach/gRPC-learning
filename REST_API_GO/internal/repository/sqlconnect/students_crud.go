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

func GetStudentsDBHandler(r *http.Request) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return nil, utils.ErrorHandler(err, "error retrieving data")
	}
	defer db.Close()

	query := "SELECT * FROM students WHERE 1=1"
	var args []interface{}

	query, args = addFiltersToGetStudentsQuery(r, query, args)
	query = addFilterSortByStudent(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)

		return nil, utils.ErrorHandler(err, "error retrieving data")
	}

	defer rows.Close()

	studentList := make([]models.Student, 0)

	for rows.Next() {
		var student models.Student

		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)

		if err != nil {
			fmt.Println(err)
			return nil, utils.ErrorHandler(err, "error retrieving data")
		}

		studentList = append(studentList, student)
	}

	return studentList, nil
}

func addFiltersToGetStudentsQuery(r *http.Request, query string, args []any) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
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

func addFilterSortByStudent(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]

	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}

			field, order := parts[0], parts[1]

			if !isValidSortOrderStudent(order) || !isValidSortFieldStudent(field) {
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

func isValidSortOrderStudent(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortFieldStudent(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
	}

	exists := validFields[field]

	return exists
}

func GetOneStudentDBHandler(id int) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return models.Student{}, utils.ErrorHandler(err, "error retrieving data")
	}
	defer db.Close()

	var student models.Student
	err = db.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err == sql.ErrNoRows {
		return models.Student{}, utils.ErrorHandler(err, fmt.Sprintf("no teacher found with id: %v", id))
	}

	if err != nil {
		fmt.Println(err)
		return models.Student{}, utils.ErrorHandler(err, "error retrieving data")
	}

	return student, nil
}

func AddedStudentsDBHandler(newStudents []models.Student) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Printf("error: %v", err)
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO students (first_name, last_name, email, class) VALUES (?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer stmt.Close()

	addedStudents := make([]models.Student, len(newStudents))

	for i, newStudent := range newStudents {
		result, err := stmt.Exec(newStudent.FirstName, newStudent.LastName, newStudent.Email, newStudent.Class)
		if err != nil {
			fmt.Println(err)
			return nil, utils.ErrorHandler(err, "error adding data")
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(err, "error adding data")
		}

		newStudent.ID = int(lastID)
		addedStudents[i] = newStudent
	}

	return addedStudents, nil
}

func PatchStudentsDBHandler(updates []map[string]interface{}, r *http.Request) error {
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

		var studentsFromDb models.Student
		err = db.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(&studentsFromDb.ID,
			&studentsFromDb.FirstName,
			&studentsFromDb.LastName,
			&studentsFromDb.Email,
			&studentsFromDb.Class)

		if err == sql.ErrNoRows {
			tx.Rollback()
			log.Printf("error: %v", err)
			return utils.ErrorHandler(err, "teacher not found")
		}

		if err != nil {
			log.Printf("error: %v", err)
			return utils.ErrorHandler(err, "invalid id")
		}

		studentVal := reflect.ValueOf(&studentsFromDb).Elem()
		studentType := studentVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}

			for i := 0; i < studentVal.NumField(); i++ {

				field := studentType.Field(i)

				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := studentVal.Field(i)
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

		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ? WHERE id = ?",
			studentsFromDb.FirstName, studentsFromDb.LastName, studentsFromDb.Email, studentsFromDb.ID)

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

func PatchOneStudentDBHandler(updates map[string]any, id int) (models.Student, error) {

	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return models.Student{}, err
	}
	defer db.Close()

	var existintStudent models.Student
	err = db.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(
		&existintStudent.ID,
		&existintStudent.FirstName,
		&existintStudent.LastName,
		&existintStudent.Email,
		&existintStudent.Class)

	if err != nil {
		log.Printf("error: %v", err)
		return models.Student{}, utils.ErrorHandler(err, "error updating data")
	}

	studentVal := reflect.ValueOf(&existintStudent).Elem()
	studentType := studentVal.Type()

	for k, v := range updates {
		for i := 0; i < studentVal.NumField(); i++ {
			field := studentType.Field(i)

			if field.Tag.Get("json") == k+",omitempty" {
				if studentVal.Field(i).CanSet() {
					studentVal.Field(i).Set(reflect.ValueOf(v).Convert(studentVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
		existintStudent.FirstName,
		existintStudent.LastName,
		existintStudent.Email,
		existintStudent.Class,
		existintStudent.ID)

	if err != nil {
		log.Printf("erro updating teacher in the DB: %v", err)
		return models.Student{}, utils.ErrorHandler(err, "error updating to the DB")
	}

	return existintStudent, nil
}

func DeleteOneStudentDBHandler(id int) (int, error) {

	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return 0, utils.ErrorHandler(err, "internal server error")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		log.Printf("error deleting student of the DB: %v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error retreaving result of deleting student of the DB: %v", err)
		return 0, utils.ErrorHandler(err, "internal server error")
	}

	if rowsAffected == 0 {
		log.Printf("no student found with the id %v on the DB", id)
		return 0, utils.ErrorHandler(err, fmt.Sprintf("no student found with id %v", id))
	}

	return id, nil
}

func DeleteStudentsDBHandlers(idStudentsToDelete []int) ([]int, error) {
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

	stmt, err := tx.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		log.Printf("error: %v", err)
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "error deleting student")
	}

	defer stmt.Close()

	var deletedId []int

	for _, id := range idStudentsToDelete {

		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(err, "error deleting student")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(err, "error deleting student")
		}

		if rowsAffected <= 0 {
			log.Printf("no teacher found with id: %v", id)
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("no student found with ID: %v", id))
		}

		deletedId = append(deletedId, id)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("error:%v", err)
		return nil, utils.ErrorHandler(err, "error deleting student")
	}

	if len(deletedId) <= 0 {
		log.Printf("error: %v", err)
		return nil, utils.ErrorHandler(err, fmt.Sprintf("no student found with ID: %v", deletedId))
	}

	return deletedId, nil
}
