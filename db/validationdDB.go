package db

import (
	"context"
	"time"

	"github.com/lscantillo/twitter-clone-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func UserExists(email string) (models.User, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("users")

	var user models.User
	condition := bson.M{"email": email}
	err := col.FindOne(ctx, condition).Decode(&user)
	ID := user.ID.Hex()
	if err != nil {
		return user, false, ID
	}

	return user, true, ID
}

func AttempLogin(email string, password string) (models.User, bool) {
	user, finded, _ := UserExists(email)
	if !finded {
		return user, false
	}
	passwordBytes := []byte(password)
	passwordDB := []byte(user.Password)
	err := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes)
	if err != nil {
		return user, false
	}
	return user, true
}
