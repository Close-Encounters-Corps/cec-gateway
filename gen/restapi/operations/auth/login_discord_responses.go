// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/Close-Encounters-Corps/cec-gateway/gen/models"
)

// LoginDiscordOKCode is the HTTP code returned for type LoginDiscordOK
const LoginDiscordOKCode int = 200

/*LoginDiscordOK Phase successful

swagger:response loginDiscordOK
*/
type LoginDiscordOK struct {

	/*
	  In: Body
	*/
	Payload *models.AuthPhaseResult `json:"body,omitempty"`
}

// NewLoginDiscordOK creates LoginDiscordOK with default headers values
func NewLoginDiscordOK() *LoginDiscordOK {

	return &LoginDiscordOK{}
}

// WithPayload adds the payload to the login discord o k response
func (o *LoginDiscordOK) WithPayload(payload *models.AuthPhaseResult) *LoginDiscordOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the login discord o k response
func (o *LoginDiscordOK) SetPayload(payload *models.AuthPhaseResult) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *LoginDiscordOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// LoginDiscordInternalServerErrorCode is the HTTP code returned for type LoginDiscordInternalServerError
const LoginDiscordInternalServerErrorCode int = 500

/*LoginDiscordInternalServerError Internal error

swagger:response loginDiscordInternalServerError
*/
type LoginDiscordInternalServerError struct {
}

// NewLoginDiscordInternalServerError creates LoginDiscordInternalServerError with default headers values
func NewLoginDiscordInternalServerError() *LoginDiscordInternalServerError {

	return &LoginDiscordInternalServerError{}
}

// WriteResponse to the client
func (o *LoginDiscordInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
