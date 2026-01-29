package middlewares

import (
	"net/http"
	"slices"
)

//api is hosted at www.myapi.com
// frontend server is at 

//Allowed Origins
var allowedOrigins =[] string {
	"https://my-origin-url.com",
	"https://www.myfrontend.com",
	"https://localhost:3000",
}

func Cors(next http.Handler) http.Handler{
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request){
		
		origin := r.Header.Get("Origin")
		
		if !slices.Contains(allowedOrigins, origin){
			http.Error(w, "Forbidden domain", http.StatusForbidden)
			return
		}
		
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authoritazion")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")
		
		if r.Method == http.MethodOptions{
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

