package handlers

import (
	"encoding/json"
	"net/http"

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
		http.Error(w, "Error while fetching data"+err.Error(), http.StatusBadRequest)
		return
	}
	if len(t.Email) == 0 {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if len(t.Password) < 6 {
		http.Error(w, "Password is required minimum with 6 characters", http.StatusBadRequest)
		return
	}
	_, finded, _ := db.UserExists(t.Email)
	if finded {
		//http.Error(w, "User already exists", http.StatusBadRequest)
		utils.RespondWithJSON(w, http.StatusBadRequest, "User already exists", nil)
		return
	}
	_id, status, err := db.CreateUser(t)
	if err != nil {
		http.Error(w, "Error while saving data"+err.Error(), http.StatusBadRequest)
		return
	}
	if !status {
		http.Error(w, "Error while saving data", http.StatusBadRequest)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, "User created successfully", _id)

}
