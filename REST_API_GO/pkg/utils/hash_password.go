package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/argon2"
)


func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		log.Printf("error: %v", err)
		return "", ErrorHandler(errors.New("failed to genereate password"), "error adding data")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	encodedHashNewPassword := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	
	
	return encodedHashNewPassword, nil
}