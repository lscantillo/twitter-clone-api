package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lscantillo/twitter-clone-api/config"
	"github.com/lscantillo/twitter-clone-api/models"
)

/*GeneroJWT genera el encriptado con JWT */
func GenerateJWT(t models.User) (string, error) {

	myKey := []byte(config.GetVariables("SECRET_KEY"))

	payload := jwt.MapClaims{
		"email":     t.Email,
		"name":      t.Name,
		"last_name": t.LastName,
		"birth_day": t.BirthDate,
		"biography": t.Biography,
		"location":  t.Location,
		"web_site":  t.WebSite,
		"_id":       t.ID.Hex(),
		"expires":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(myKey)
	if err != nil {
		return tokenStr, err
	}
	return tokenStr, nil
}
