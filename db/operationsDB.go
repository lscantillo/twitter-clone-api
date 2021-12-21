package db

import (
	"context"
	"time"

	"github.com/lscantillo/twitter-clone-api/models"
	"github.com/lscantillo/twitter-clone-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(user models.User) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("users")
	user.Password, _ = utils.EncriptPassword(user.Password)

	result, err := col.InsertOne(ctx, user)
	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil
}

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
