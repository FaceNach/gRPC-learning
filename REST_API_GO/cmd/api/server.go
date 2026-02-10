package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "rest_api_go/internal/api/middlewares"
	"rest_api_go/internal/api/router"
	"rest_api_go/internal/repository/sqlconnect"
	"rest_api_go/pkg/utils"

	"github.com/joho/godotenv"
)

func main() {
	
	err := godotenv.Load()
	if err != nil {
		return 
	}

	fmt.Println("Connecting to DB")
	_, err = sqlconnect.ConnectDb()
	if err != nil {
		utils.ErrorHandler(err, "error connecting to DB")
		return
	}
	
	fmt.Println("Connected to DB/MariaDB successfully")
	
	router := router.MainRouter()

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// rl := mw.NewRateLimiter(5, time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery: true,
	// 	CheckBody: true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	WhiteList: []string{"sortBy", "sortOrder", "name", "age", "class"},

	// }

	//secureMux := mw.Cors(rl.Middleware(mw.ResponseTime(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	//secureMux := utils.ApplyMidlewares(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHeaders, mw.ResponseTime, rl.Middleware, mw.Cors)
	secureMux := mw.SecurityHeaders(router)

	server := &http.Server{
		Addr:      os.Getenv("API_PORT"),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Listening server on port", os.Getenv("API_PORT"))
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error creating server: ", err)
	}
}
