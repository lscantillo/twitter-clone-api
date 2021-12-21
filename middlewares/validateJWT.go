package middlewares

import (
	"net/http"

	"github.com/lscantillo/twitter-clone-api/utils"
)

func ValidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, _, err := utils.TokenProcess(r.Header.Get("Authorization"))
		if err != nil {
			utils.RespondWithJSON(w, http.StatusUnauthorized, "Unauthorized: Token error"+err.Error(), nil)
			return
		}
		next.ServeHTTP(w, r)
	}
}
