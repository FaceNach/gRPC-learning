package sqlconnect

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"rest_api_go/internal/models"
	"rest_api_go/pkg/utils"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

func GetExecsDBHandler(r *http.Request) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return nil, utils.ErrorHandler(err, "error retrieving data")
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE 1=1"
	var args []interface{}

	query, args = addFiltersToGetExecsQuery(r, query, args)
	query = addFilterSortByExec(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)

		return nil, utils.ErrorHandler(err, "error retrieving data")
	}

	defer rows.Close()

	execList := make([]models.Exec, 0)

	for rows.Next() {
		var exec models.Exec

		err := rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role)

		if err != nil {
			fmt.Println(err)
			return nil, utils.ErrorHandler(err, "error retrieving data")
		}

		execList = append(execList, exec)
	}

	return execList, nil
}

func addFiltersToGetExecsQuery(r *http.Request, query string, args []any) (string, []interface{}) {
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

func addFilterSortByExec(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]

	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}

			field, order := parts[0], parts[1]

			if !isValidSortOrderExec(order) || !isValidSortFieldExec(field) {
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

func isValidSortOrderExec(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortFieldExec(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"role":       true,
		"username":   true,
	}

	exists := validFields[field]

	return exists
}

func GetOneExecDBHandler(id int) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return models.Exec{}, utils.ErrorHandler(err, "error retrieving data")
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE id = ?", id).Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role)
	if err == sql.ErrNoRows {
		return models.Exec{}, utils.ErrorHandler(err, fmt.Sprintf("no exec found with id: %v", id))
	}

	if err != nil {
		fmt.Println(err)
		return models.Exec{}, utils.ErrorHandler(err, "error retrieving data")
	}

	return exec, nil
}

func AddExecsDBHandler(newExecs []models.Exec) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Printf("error: %v", err)
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO execs (first_name, last_name, email, username, password, role, inactive_status) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer stmt.Close()

	addedExecs := make([]models.Exec, len(newExecs))

	for i, newExec := range newExecs {
		if newExec.Password == "" {
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(errors.New("password can't be empty"), "Password can't be empty")
		}

		salt := make([]byte, 16)
		_, err := rand.Read(salt)
		if err != nil {
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(errors.New("failed to genereate password"), "error adding data")
		}

		hash := argon2.IDKey([]byte(newExec.Password), salt, 1, 64*1024, 4, 32)
		saltBase64  := base64.StdEncoding.EncodeToString(salt)
		hashBase64 := base64.StdEncoding.EncodeToString(hash)
		encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
		
		newExec.Password = encodedHash

		result, err := stmt.Exec(newExec.FirstName, newExec.LastName, newExec.Email, newExec.Username, newExec.Password, newExec.Role, newExec.InactiveStatus)
		if err != nil {
			fmt.Println(err)
			return nil, utils.ErrorHandler(err, "error adding data")
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(err, "error adding data")
		}

		newExec.ID = int(lastID)
		addedExecs[i] = newExec
	}

	return addedExecs, nil
}

func PatchExecsDBHandler(updates []map[string]interface{}) error {
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
	defer tx.Rollback()

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

		var execsFromDB models.Exec
		err = db.QueryRow("SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?", id).Scan(&execsFromDB.ID,
			&execsFromDB.FirstName,
			&execsFromDB.LastName,
			&execsFromDB.Email,
			&execsFromDB.Username,
		)

		if err == sql.ErrNoRows {
			tx.Rollback()
			log.Printf("error: %v", err)
			return utils.ErrorHandler(err, "teacher not found")
		}

		if err != nil {
			log.Printf("error: %v", err)
			return utils.ErrorHandler(err, "invalid id")
		}

		execVal := reflect.ValueOf(&execsFromDB).Elem()
		execType := execVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}

			for i := 0; i < execVal.NumField(); i++ {

				field := execType.Field(i)

				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := execVal.Field(i)
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

		_, err = tx.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ?",
			execsFromDB.FirstName, execsFromDB.LastName, execsFromDB.Email, execsFromDB.Username, execsFromDB.ID)

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

func PatchOneExecDBHandler(updates map[string]any, id int) (models.Exec, error) {

	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return models.Exec{}, err
	}
	defer db.Close()

	var existingExec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?", id).Scan(
		&existingExec.ID,
		&existingExec.FirstName,
		&existingExec.LastName,
		&existingExec.Email,
		&existingExec.Username,
	)

	if err != nil {
		log.Printf("error: %v", err)
		return models.Exec{}, utils.ErrorHandler(err, "error updating data")
	}

	execVal := reflect.ValueOf(&existingExec).Elem()
	execType := execVal.Type()

	for k, v := range updates {
		for i := 0; i < execVal.NumField(); i++ {
			field := execType.Field(i)

			if field.Tag.Get("json") == k+",omitempty" {
				if execVal.Field(i).CanSet() {
					execVal.Field(i).Set(reflect.ValueOf(v).Convert(execVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ?",
		existingExec.FirstName,
		existingExec.LastName,
		existingExec.Email,
		existingExec.Username,
		existingExec.ID)

	if err != nil {
		log.Printf("erro updating exec in the DB: %v", err)
		return models.Exec{}, utils.ErrorHandler(err, "error updating to the DB")
	}

	return existingExec, nil
}

func DeleteOneExecDBHandler(id int) (int, error) {

	db, err := ConnectDb()
	if err != nil {
		fmt.Println(err)
		return 0, utils.ErrorHandler(err, "internal server error")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM execs WHERE id = ?", id)
	if err != nil {
		log.Printf("error deleting exec of the DB: %v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error retreaving result of deleting exec of the DB: %v", err)
		return 0, utils.ErrorHandler(err, "internal server error")
	}

	if rowsAffected == 0 {
		log.Printf("no exec found with the id %v on the DB", id)
		return 0, utils.ErrorHandler(err, fmt.Sprintf("no exec found with id %v", id))
	}

	return id, nil
}

func DeleteExecsDBHandlers(idExecsToDelete []int) ([]int, error) {
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

	stmt, err := tx.Prepare("DELETE FROM execs WHERE id = ?")
	if err != nil {
		log.Printf("error: %v", err)
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "error deleting execs")
	}

	defer stmt.Close()

	var deletedId []int

	for _, id := range idExecsToDelete {

		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(err, "error deleting exec")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			log.Printf("error: %v", err)
			return nil, utils.ErrorHandler(err, "error deleting exec")
		}

		if rowsAffected <= 0 {
			log.Printf("no exec found with id: %v", id)
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("no exec found with ID: %v", id))
		}

		deletedId = append(deletedId, id)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("error:%v", err)
		return nil, utils.ErrorHandler(err, "error deleting exec")
	}

	if len(deletedId) <= 0 {
		log.Printf("error: %v", err)
		return nil, utils.ErrorHandler(err, fmt.Sprintf("no exec found with ID: %v", deletedId))
	}

	return deletedId, nil
}
