package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lscantillo/twitter-clone-api/db"
	"github.com/lscantillo/twitter-clone-api/models"
	"github.com/lscantillo/twitter-clone-api/utils"
)

var UserID string

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, "Welcome to Twitter Clone API", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var t models.User
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while fetching data"+err.Error(), nil)
		return
	}
	if len(t.Email) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Email is required", nil)
		return
	}
	if len(t.Password) < 6 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Password is required minimum with 6 characters", nil)
		return
	}
	_, finded, _ := db.UserExists(t.Email)
	if finded {
		utils.RespondWithJSON(w, http.StatusBadRequest, "User already exists", nil)
		return
	}
	_id, status, err := db.CreateUser(t)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving data"+err.Error(), nil)
		return
	}
	if !status {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving data"+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, "User created successfully", _id)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	var t models.User

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "User or/and password is invalid"+err.Error(), nil)
		return
	}
	if len(t.Email) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Email is required", nil)
		return
	}
	document, exists := db.AttempLogin(t.Email, t.Password)
	if !exists {
		utils.RespondWithJSON(w, http.StatusBadRequest, "User or/and password is invalid", nil)
		return
	}
	jwtKey, err := utils.GenerateJWT(document)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while generating token"+err.Error(), nil)
		return
	}
	resp := models.Token{
		Token: jwtKey,
	}
	utils.RespondWithJSON(w, http.StatusOK, "Login success", resp)
	expirationTime := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: expirationTime,
	})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}
	profile, err := db.SearchProfile(ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while searching user"+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "User found", profile)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var t models.User
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while fetching data"+err.Error(), nil)
		return
	}
	var status bool

	UserID = utils.GetUserID(r.Header.Get("Authorization"))
	status, err = db.UpdateRegister(t, UserID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving data"+err.Error(), nil)
		return
	}
	if !status {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving data", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "User updated successfully", nil)
}

func CreateTweet(w http.ResponseWriter, r *http.Request) {
	var message models.Tweet
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while fetching data"+err.Error(), nil)
		return
	}

	UserID = utils.GetUserID(r.Header.Get("Authorization"))

	register := models.SaveTweet{
		UserID:  UserID,
		Message: message.Message,
		Date:    time.Now(),
	}
	_, status, err := db.CreateTweet(register)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving data"+err.Error(), nil)
		return
	}

	if !status {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving data", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, "Tweet created successfully", nil)
}

func GetTweets(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}
	if len(r.URL.Query().Get("page")) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Page is required", nil)
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Page is invalid", nil)
		return
	}
	pag := int64(page)
	response, correct, _ := db.GetTweets(ID, pag)
	if !correct {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while fetching tweets", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Tweets found", response)
}

func DeleteTweet(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Tweet ID is required", nil)
		return
	}
	UserID = utils.GetUserID(r.Header.Get("Authorization"))
	err := db.DeleteTweet(ID, UserID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while deleting tweet"+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Tweet deleted successfully", nil)
}

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	UserID = utils.GetUserID(r.Header.Get("Authorization"))

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while fetching data"+err.Error(), nil)
		return
	}
	var extension = strings.Split(handler.Filename, ".")[1]
	var file_route string = "uploads/avatars/" + UserID + "." + extension
	f, err := os.OpenFile(file_route, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving avatar"+err.Error(), nil)
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while copying avatar"+err.Error(), nil)
		return
	}
	var user models.User
	var status bool

	user.Avatar = UserID + "." + extension
	status, err = db.UpdateRegister(user, UserID)
	if err != nil || !status {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving avatar in DB"+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Avatar saved successfully", nil)

}

func UploadBanner(w http.ResponseWriter, r *http.Request) {
	UserID = utils.GetUserID(r.Header.Get("Authorization"))

	file, handler, err := r.FormFile("banner")
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while fetching data"+err.Error(), nil)
		return
	}
	var extension = strings.Split(handler.Filename, ".")[1]
	var file_route string = "uploads/banners/" + UserID + "." + extension
	f, err := os.OpenFile(file_route, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving banner"+err.Error(), nil)
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while copying banner"+err.Error(), nil)
		return
	}
	var user models.User
	var status bool

	user.Banner = UserID + "." + extension
	status, err = db.UpdateRegister(user, UserID)
	if err != nil || !status {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving banner in DB"+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "banner saved successfully", nil)

}

func GetAvatar(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}
	profile, err := db.SearchProfile(ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "User not found"+err.Error(), nil)
		return
	}
	OpenFile, err := os.Open("uploads/avatars/" + profile.Avatar)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while opening file"+err.Error(), nil)
		return
	}
	_, err = io.Copy(w, OpenFile)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while copying file: "+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Avatar found", OpenFile)
}

func GetBanner(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}
	profile, err := db.SearchProfile(ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "User not found"+err.Error(), nil)
		return
	}

	OpenFile, err := os.Open("uploads/banners/" + profile.Banner)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while opening file: "+err.Error(), nil)
		return
	}

	_, err = io.Copy(w, OpenFile)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while copying file: "+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Banner found", OpenFile)
}

func CreateRelation(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	if len(ID) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}
	UserID = utils.GetUserID(r.Header.Get("Authorization"))
	t := models.Relation{
		UserID:         UserID,
		UserRelationID: ID,
	}
	status, err := db.InsertRelation(t)
	if err != nil || !status {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while saving relation"+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, "Relation saved successfully", nil)

}

func DeleteRelation(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	UserID = utils.GetUserID(r.Header.Get("Authorization"))
	t := models.Relation{
		UserID:         UserID,
		UserRelationID: ID,
	}
	status, err := db.DeleteRelation(t)
	if err != nil || !status {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while deleting relation"+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, "Relation deleted successfully", nil)
}

func GetRelation(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	UserID = utils.GetUserID(r.Header.Get("Authorization"))
	t := models.Relation{
		UserID:         UserID,
		UserRelationID: ID,
	}
	var resp models.ResponseRelation

	status, err := db.GetRelation(t)
	if err != nil || !status {
		resp.Status = false
	} else {
		resp.Status = true
	}
	utils.RespondWithJSON(w, http.StatusCreated, "Relation found", resp)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	typeUser := r.URL.Query().Get("type")
	page := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")

	pagTemp, err := strconv.Atoi(page)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while converting page"+err.Error(), nil)
		return
	}
	pag := int64(pagTemp)
	UserID = utils.GetUserID(r.Header.Get("Authorization"))
	result, status := db.GetUsers(UserID, pag, search, typeUser)
	if !status {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while getting users"+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, "Users found", result)
}

func ReadTweets(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("page")) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Page is required", nil)
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while converting page"+err.Error(), nil)
		return
	}
	UserID = utils.GetUserID(r.Header.Get("Authorization"))
	response, correct := db.GetFollowersTweets(UserID, page)
	if !correct {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Error while getting tweets"+err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, "Tweets found", response)

}
