package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lscantillo/twitter-clone-api/models"
)

func RespondWithJSON(w http.ResponseWriter, code int, message string, payload interface{}) {

	resp := handleResponse(code, message, payload)
	response, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Can not convert payload to JSON - %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func handleResponse(code int, msg string, data interface{}) models.Response {
	return models.Response{
		Code:     code,
		Message:  msg,
		Response: data,
	}
}
