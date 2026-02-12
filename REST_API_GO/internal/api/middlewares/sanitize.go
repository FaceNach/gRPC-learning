package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"rest_api_go/pkg/utils"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {
	fmt.Println("********Initializing XSSMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("********Inside XSSMiddleware")

		//Sanitize URL Path
		sanitizePath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Original path: ", r.URL.Path)
		fmt.Println("Sanitize path: ", sanitizePath)

		params := r.URL.Query()
		sanitizedQuery := make(map[string][]string)

		for key, values := range params {
			sanitizedKey, err := clean(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var sanitizedValues []string

			for _, value := range values {
				cleanValue, err := clean(value)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				sanitizedValues = append(sanitizedValues, cleanValue.(string))
			}
			sanitizedQuery[sanitizedKey.(string)] = sanitizedValues

			fmt.Printf("Original Query %s: %s \n", key, strings.Join(values, ", "))
			fmt.Printf("Original Query %s: %s \n", sanitizedKey, strings.Join(sanitizedValues, ", "))

		}

		r.URL.Path = sanitizePath.(string)
		r.URL.RawQuery = url.Values(sanitizedQuery).Encode()
		fmt.Println("Updated URL: ", r.URL.String())

		if r.Header.Get("Content-Type") == "application/json" {
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					utils.ErrorHandler(err, "")
					http.Error(w, "Error reading request body", http.StatusBadRequest)
					return
				}

				bodyString := strings.TrimSpace(string(bodyBytes))
				fmt.Println("Original body: ", bodyString)

				r.Body = io.NopCloser(bytes.NewReader([]byte(bodyString)))

				if len(bodyString) > 0 {
					var inputData interface{}
					err = json.NewDecoder(bytes.NewReader([]byte(bodyString))).Decode(&inputData)
					if err != nil {
						utils.ErrorHandler(err, "")
						http.Error(w, "Error reading request body", http.StatusBadRequest)
						return
					}
					
					fmt.Println("Original JSON DATA :", inputData)
					sanitizedData, err := clean(inputData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					
					fmt.Println("Sanitized JSON DATA :", sanitizedData)
					sanitizedBody, err := json.Marshal(sanitizedData)
					if err != nil {
						utils.ErrorHandler(err, "")
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					
					r.Body = io.NopCloser(bytes.NewReader(sanitizedBody))
					fmt.Println("Sanitized body: ", string(sanitizedBody))

				} else {
					log.Println("Request body its empty")
				}

			} else {
				log.Println("No body in the request")
				return
			}
		} else {
			log.Printf("Received request with unsupported Content-Type: %s . Expected application/json", r.Header.Get("Content-Type"))
			http.Error(w, "unsupported content-type, please use application/json", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
		fmt.Println("********Sending response from  XSSMiddleware")
	})
}

func clean(data any) (any, error) {
	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v, nil
	case []any:
		for i, value := range v {
			v[i] = sanitizeValue(value)
		}
		return v, nil
	case string:
		return sanitizeString(v), nil
	default:

		return nil, utils.ErrorHandler(fmt.Errorf("unsupported type: %T", data), fmt.Sprintf("unsupported type: %T", data))
	}
}

func sanitizeValue(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return sanitizeString(v)

	case map[string]any:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v
	case []any:
		for i, value := range v {
			v[i] = sanitizeValue(value)
		}
		return v
	default:
		return v
	}
}

func sanitizeString(value string) string {
	return bluemonday.UGCPolicy().Sanitize(value)
}
