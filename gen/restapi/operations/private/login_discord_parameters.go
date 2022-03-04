// Code generated by go-swagger; DO NOT EDIT.

package private

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewLoginDiscordParams creates a new LoginDiscordParams object
//
// There are no default values defined in the spec.
func NewLoginDiscordParams() LoginDiscordParams {

	return LoginDiscordParams{}
}

// LoginDiscordParams contains all the bound params for the login discord operation
// typically these are obtained from a http.Request
//
// swagger:parameters loginDiscord
type LoginDiscordParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*URL to redirect on success
	  In: query
	*/
	RedirectURL *string
	/*Second phase: State to fetch from CEC Auth
	  In: query
	*/
	State *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewLoginDiscordParams() beforehand.
func (o *LoginDiscordParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qRedirectURL, qhkRedirectURL, _ := qs.GetOK("redirect_url")
	if err := o.bindRedirectURL(qRedirectURL, qhkRedirectURL, route.Formats); err != nil {
		res = append(res, err)
	}

	qState, qhkState, _ := qs.GetOK("state")
	if err := o.bindState(qState, qhkState, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindRedirectURL binds and validates parameter RedirectURL from query.
func (o *LoginDiscordParams) bindRedirectURL(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.RedirectURL = &raw

	return nil
}

// bindState binds and validates parameter State from query.
func (o *LoginDiscordParams) bindState(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.State = &raw

	return nil
}
