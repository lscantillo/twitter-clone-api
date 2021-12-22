package utils

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/lscantillo/twitter-clone-api/config"
	"github.com/lscantillo/twitter-clone-api/db"
	"github.com/lscantillo/twitter-clone-api/models"
)

/*Email valor de Email usado en todos los EndPoints */
var Email string

/*IDUsuario es el ID devuelto del modelo, que se usará en todos los EndPoints */
var UserID string

func TokenProcess(token string) (*models.Claim, bool, string, error) {
	myKey := []byte(config.GetVariables("SECRET_KEY"))
	claims := &models.Claim{}
	splitToken := strings.Split(token, "Bearer ")
	if len(splitToken) != 2 {
		return claims, false, "", errors.New("invalid token")
	}
	token = strings.TrimSpace(splitToken[1])
	tokenClaims, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		_, finded, _ := db.UserExists(claims.Email)
		if finded {
			Email = claims.Email
			UserID = claims.ID.Hex()
		}
		return claims, finded, UserID, nil
	}
	if !tokenClaims.Valid {
		return claims, false, "", errors.New("invalid token")
	}
	return claims, false, "", err
}
