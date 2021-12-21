package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lscantillo/twitter-clone-api/db"
	"github.com/lscantillo/twitter-clone-api/models"
	"github.com/lscantillo/twitter-clone-api/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
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
