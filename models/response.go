package models

// Response struct is used to return the response to the client
type Response struct {
	Message  string      `json:"message"`
	Code     int         `json:"code"`
	Response interface{} `json:"response"`
}
