package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	mw "rest_api_go/internal/api/middlewares"
	"rest_api_go/internal/api/router"
)

func main() {
	port := 3000

	router := router.Router()

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
		Addr:      ":3000",
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Listening server on port:", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error creating server: ", err)
	}
}
