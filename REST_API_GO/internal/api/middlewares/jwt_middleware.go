package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"rest_api_go/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)



func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := r.Cookie("Bearer")
		if err != nil {
			log.Printf("error: %v", err)
			http.Error(w, "authorization header missing", http.StatusUnauthorized)
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (any, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(jwtSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {

			if errors.Is(err, jwt.ErrTokenExpired) {
				log.Printf("error: %v", err)
				http.Error(w, "token expired", http.StatusUnauthorized)
				return
			}

			if errors.Is(err, jwt.ErrTokenMalformed) {
				log.Printf("error: %v", err)
				http.Error(w, "token malformed", http.StatusUnauthorized)
				return
			}

			log.Printf("error: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if parsedToken.Valid {
			//fmt.Println("Valid JWT")

		} else {
			http.Error(w, "invalid login Token", http.StatusUnauthorized)
			fmt.Println("Invalid JWT", token.Value)
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(w, "invalid login Token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextKey("role"), claims["role"])
		ctx = context.WithValue(ctx, utils.ContextKey("expiresAt"), claims["exp"])
		ctx = context.WithValue(ctx, utils.ContextKey("username"), claims["user"])
		ctx = context.WithValue(ctx, utils.ContextKey("userId"), claims["uid"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
