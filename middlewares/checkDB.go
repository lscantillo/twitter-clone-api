package middlewares

import (
	"net/http"

	"github.com/lscantillo/twitter-clone-api/db"
)

func CheckDB(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db.CheckDB() == 0 {
			http.Error(w, "Error while connecting to the database", 500)
			return
		}
		next.ServeHTTP(w, r)
	}
}
