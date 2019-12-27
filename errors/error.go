package errors

import (
	"encoding/json"
	"net/http"
)

// Error the super type of all errors used inside cerberus
type Error interface {
	Error() string // Returns a string depicting the reason of error, that should be printed in the console
	Message() string
	StatusCode() int
}

// AuthenticationError invalid authentication error
type AuthenticationError struct {
	Mess string
	Err  string
}

// Error returns the error message for the AuthenticationError
func (e AuthenticationError) Error() string {
	return e.Err
}

// Message returns the error message to be sent with the response for the AuthenticationError
func (e AuthenticationError) Message() string {
	if e.Mess != "" {
		return e.Mess
	}
	return e.Err
}

// StatusCode returns the status code to be sent in the response for the AuthenticationError
func (e AuthenticationError) StatusCode() int {
	return http.StatusForbidden
}

// AuthorizationError invalid authorization error
type AuthorizationError struct {
	Mess string
	Err  string
}

// Error returns the error message for the AuthorizationError
func (e AuthorizationError) Error() string {
	return e.Err
}

// Message returns the error message to be sent with the response for the AuthorizationError
func (e AuthorizationError) Message() string {
	if e.Mess != "" {
		return e.Mess
	}
	return e.Err
}

// StatusCode returns the status code to be sent in the response for the AuthorizationError
func (e AuthorizationError) StatusCode() int {
	return http.StatusUnauthorized
}

// InvalidInputError invalid input error
type InvalidInputError struct {
	Mess string
	Err  string
}

// Error returns the error message for the InvalidInputError
func (e InvalidInputError) Error() string {
	return e.Err
}

// Message returns the error message to be sent with the response for the InvalidInputError
func (e InvalidInputError) Message() string {
	if e.Mess != "" {
		return e.Mess
	}
	return e.Err
}

// StatusCode returns the status code to be sent in the response for the InvalidInputError
func (e InvalidInputError) StatusCode() int {
	return http.StatusBadRequest
}

//InternalError error that has occured internally
type InternalError struct {
	Mess string
	Err  string
}

// Error returns the error message for the InternalError
func (e InternalError) Error() string {
	return e.Err
}

// Message returns the error message to be sent with the response for the InternalError
func (e InternalError) Message() string {
	if e.Mess != "" {
		return e.Mess
	}
	return "Internal Error, please contact admin"
}

// StatusCode returns the status code to be sent in the response for the InternalError
func (e InternalError) StatusCode() int {
	return http.StatusInternalServerError
}

//InvalidMethodError the url does not support the given method
type InvalidMethodError struct {
	Mess string
	Err  string
}

// Error returns the error message for the InvalidMethodError
func (e InvalidMethodError) Error() string {
	return e.Err
}

// Message returns the error message to be sent with the response for the InvalidMethodError
func (e InvalidMethodError) Message() string {
	if e.Mess != "" {
		return e.Mess
	}
	return "Invalid Method"
}

// StatusCode returns the status code to be sent in the response for the InvalidMethodError
func (e InvalidMethodError) StatusCode() int {
	return http.StatusMethodNotAllowed
}

//CustomStatusError an error with a custom message and status code
type CustomStatusError struct {
	Mess   string
	Err    string
	Status int
}

// Error returns the error message for the CustomStatusError
func (e CustomStatusError) Error() string {
	return e.Err
}

// Message returns the error message to be sent with the response for the CustomStatusError
func (e CustomStatusError) Message() string {
	if e.Mess != "" {
		return e.Mess
	}
	return "Invalid Method"
}

// StatusCode returns the status code to be sent in the response for the CustomStatusError
func (e CustomStatusError) StatusCode() int {
	return e.Status
}

// struct to hold complete error messages
func GetMessageBytes(err Error) []byte {
	vals := make(map[string]string)
	vals["message"] = err.Message()
	vals["status"] = "error"
	js, _ := json.Marshal(vals)
	return js
}
