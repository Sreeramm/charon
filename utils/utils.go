package utils

import (
	"encoding/base64"
	"net/http"
)

//GetBasicAuthHeader returns the header value given the username and password
func GetBasicAuthHeader(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

//IsSuccess checks if the given status code is a success code
func IsSuccess(statusCode int) bool {
	switch statusCode {
	case http.StatusAccepted, http.StatusOK, http.StatusCreated:
		return true
	default:
		return false
	}
}
