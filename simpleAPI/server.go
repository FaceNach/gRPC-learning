package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)


func main() {
	
	http.HandleFunc("/orders", func (w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Handling incoming orders")
	})
	
	http.HandleFunc("/users", func (w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Handling users")
	})
	
	
	port :=3000
	
	//Load the TLS cert and key
	
	cert := "cert.pem"
	key := "key.pem"
	
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: nil,
		TLSConfig: tlsConfig,
	}
	
	http2.ConfigureServer(server, &http2.Server{})
	
	fmt.Println("Listening on   port: ", port)
	
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("error creating server", err)
	}
	
	//HTTP 1.1 server without TLS
	// err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	// if err != nil {
	// 	log.Fatalln("error creating server: ", err)
	// }
	
	
}