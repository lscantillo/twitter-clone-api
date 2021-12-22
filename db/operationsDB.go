package db

import (
	"context"
	"time"

	"github.com/lscantillo/twitter-clone-api/jwt"
	"github.com/lscantillo/twitter-clone-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUser(user models.User) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("users")
	user.Password, _ = jwt.EncriptPassword(user.Password)

	result, err := col.InsertOne(ctx, user)
	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil
}

func SearchProfile(ID string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("users")

	var user models.User
	objID, _ := primitive.ObjectIDFromHex(ID)
	condition := bson.M{
		"_id": objID,
	}
	err := col.FindOne(ctx, condition).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func UpdateRegister(u models.User, ID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("users")

	register := make(map[string]interface{})

	register["birth_date"] = u.BirthDate
	if len(u.Name) > 0 {
		register["name"] = u.Name
	}
	if len(u.LastName) > 0 {
		register["last_name"] = u.LastName
	}
	if len(u.Avatar) > 0 {
		register["avatar"] = u.Avatar
	}
	if len(u.Banner) > 0 {
		register["banner"] = u.Banner
	}
	if len(u.Biography) > 0 {
		register["biography"] = u.Biography
	}
	if len(u.Location) > 0 {
		register["location"] = u.Location
	}
	if len(u.WebSite) > 0 {
		register["web_site"] = u.WebSite
	}

	updateString := bson.M{
		"$set": register,
	}
	objID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{
		"_id": bson.M{
			"$eq": objID,
		},
	}
	_, err := col.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateTweet(t models.SaveTweet) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("tweets")

	register := bson.M{
		"user_id": t.UserID,
		"message": t.Message,
		"date":    t.Date,
	}
	result, err := col.InsertOne(ctx, register)
	if err != nil {
		return "", false, err
	}
	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.String(), true, nil
}

func GetTweets(ID string, page int64) ([]*models.ReturnTweets, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("tweets")

	var tweets []*models.ReturnTweets

	condition := bson.M{
		"user_id": ID,
	}
	opts := options.Find()
	opts.SetLimit(20)
	opts.SetSort(bson.D{{Key: "date", Value: -1}})
	opts.SetSkip((page - 1) * 20)

	cursor, err := col.Find(ctx, condition, opts)
	if err != nil {
		return tweets, false, err
	}

	for cursor.Next(context.TODO()) {
		var tweet models.ReturnTweets
		err := cursor.Decode(&tweet)
		if err != nil {
			return tweets, false, err
		}
		tweets = append(tweets, &tweet)
	}
	return tweets, true, nil
}

func DeleteTweet(TweetID string, UserID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("tweets")

	objID, _ := primitive.ObjectIDFromHex(TweetID)

	condition := bson.M{
		"_id":     objID,
		"user_id": UserID,
	}
	_, err := col.DeleteOne(ctx, condition)
	return err
}

func InsertRelation(t models.Relation) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("relations")

	_, err := col.InsertOne(ctx, t)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteRelation(t models.Relation) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("relations")

	_, err := col.DeleteOne(ctx, t)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetRelation(t models.Relation) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("relations")

	condition := bson.M{
		"user_id":          t.UserID,
		"user_relation_id": t.UserRelationID,
	}
	var relation models.Relation
	err := col.FindOne(ctx, condition).Decode(&relation)
	if err != nil {
		return false, err
	}
	return true, nil

}

func GetUsers(ID string, page int64, search string, typ string) ([]*models.User, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("users")

	var results []*models.User
	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * 20)
	findOptions.SetLimit(20)

	query := bson.M{
		"name": bson.M{"$regex": `(?i)` + search},
	}
	cur, err := col.Find(ctx, query, findOptions)
	if err != nil {
		return results, false
	}

	var finded, include bool

	for cur.Next(ctx) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return results, false
		}
		var r models.Relation
		r.UserID = ID
		r.UserRelationID = user.ID.Hex()

		include = false

		finded, err = GetRelation(r)
		if typ == "new" && !finded {
			include = true
		}
		if typ == "follow" && finded {
			include = true
		}
		if r.UserRelationID == ID {
			include = false
		}

		if include {
			user.Password = ""
			user.Biography = ""
			user.WebSite = ""
			user.Location = ""
			user.Banner = ""
			user.Email = ""

			results = append(results, &user)
		}
	}
	err = cur.Err()
	if err != nil {
		return results, false
	}
	cur.Close(ctx)
	return results, true
}

func GetFollowersTweets(ID string, page int) ([]models.ResponseFollowersTweets, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("twitter")
	col := db.Collection("relations")

	skip := (page - 1) * 20

	conditions := make([]bson.M, 0)
	conditions = append(conditions, bson.M{"$match": bson.M{"usuarioid": ID}})
	conditions = append(conditions, bson.M{
		"$lookup": bson.M{
			"from":         "tweets",
			"localField":   "user_relation_id",
			"foreignField": "user_id",
			"as":           "tweets",
		}})
	conditions = append(conditions, bson.M{"$unwind": "$tweets"})
	conditions = append(conditions, bson.M{"$sort": bson.M{"tweets.date": -1}})
	conditions = append(conditions, bson.M{"$skip": skip})
	conditions = append(conditions, bson.M{"$limit": 20})

	cursor, err := col.Aggregate(ctx, conditions)
	var result []models.ResponseFollowersTweets
	err = cursor.All(ctx, &result)
	if err != nil {
		return result, false
	}
	return result, true

}
