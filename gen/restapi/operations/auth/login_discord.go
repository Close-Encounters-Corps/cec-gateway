// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// LoginDiscordHandlerFunc turns a function with the right signature into a login discord handler
type LoginDiscordHandlerFunc func(LoginDiscordParams) middleware.Responder

// Handle executing the request and returning a response
func (fn LoginDiscordHandlerFunc) Handle(params LoginDiscordParams) middleware.Responder {
	return fn(params)
}

// LoginDiscordHandler interface for that can handle valid login discord params
type LoginDiscordHandler interface {
	Handle(LoginDiscordParams) middleware.Responder
}

// NewLoginDiscord creates a new http.Handler for the login discord operation
func NewLoginDiscord(ctx *middleware.Context, handler LoginDiscordHandler) *LoginDiscord {
	return &LoginDiscord{Context: ctx, Handler: handler}
}

/* LoginDiscord swagger:route GET /0/login/discord auth private loginDiscord

Authenticate using Discord

*/
type LoginDiscord struct {
	Context *middleware.Context
	Handler LoginDiscordHandler
}

func (o *LoginDiscord) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewLoginDiscordParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
