package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	core_models "github.com/Close-Encounters-Corps/cec-core/gen/models"
	"github.com/Close-Encounters-Corps/cec-gateway/gen/models"
	"github.com/Close-Encounters-Corps/cec-gateway/gen/restapi"
	"github.com/Close-Encounters-Corps/cec-gateway/gen/restapi/operations"
	"github.com/Close-Encounters-Corps/cec-gateway/gen/restapi/operations/private"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/config"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/gateway"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/util"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var COMMITSHA string

func main() {
	log.Println("Commit:", COMMITSHA)
	cfg := config.Config{
		JaegerUrl: requireEnv("CEC_JAEGER"),
		Urls: &config.UrlSet{
			External: requireEnv("CEC_URLS_EXTERNAL"),
			Core:     requireEnv("CEC_URLS_CORE"),
		},
	}
	listenport := requireEnv("CEC_LISTENPORT")
	port, err := strconv.Atoi(listenport)
	if err != nil {
		log.Fatalln(err)
	}
	app := gateway.Gateway{
		Cfg:         &cfg,
		Environment: requireEnv("CEC_ENVIRONMENT"),
		Client: &http.Client{
			Timeout:   1 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
	app.SetupAll()
	defer app.Close(context.Background())
	swagger, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewCecGwAPI(swagger)
	api.PrivateLoginDiscordHandler = private.LoginDiscordHandlerFunc(
		func(ldp private.LoginDiscordParams) middleware.Responder {
			ctx, span := gateway.NewSpan(ldp.HTTPRequest.Context(), "controller.login.discord", nil)
			defer span.End()
			internalError := func(err error) middleware.Responder {
				gateway.AddSpanError(span, err)
				gateway.FailSpan(span, "internal error")
				return private.NewLoginDiscordInternalServerError().WithPayload(&models.Error{
					RequestID: span.SpanContext().TraceID().String(),
				})
			}
			state := swag.StringValue(ldp.State)
			gateway.AddSpanTags(span, map[string]string{"request.state": state})
			u, err := url.Parse(cfg.Urls.External)
			if err != nil {
				return internalError(err)
			}
			u.Path = ldp.HTTPRequest.URL.String()
			req, err := util.V1LoginDiscordRequest(
				ctx,
				cfg.Urls.Core,
				"/v1/login/discord",
				u.String(),
				state,
			)
			if err != nil {
				return internalError(err)
			}
			span.AddEvent("making request")
			resp, err := app.Client.Do(req)
			span.AddEvent("request done")
			if err != nil {
				return internalError(err)
			}
			gateway.AddSpanTags(span, map[string]string{"response.status": fmt.Sprint(resp.StatusCode)})
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return internalError(err)
			}
			log.Println(string(body))
			var out core_models.AuthPhaseResult
			err = json.Unmarshal(body, &out)

			if err != nil {
				gateway.AddSpanTags(span, map[string]string{"response.body": string(body)})
				return internalError(err)
			}
			var user *models.User = nil
			if out.User != nil {
				user = &models.User{}
				user.ID = out.User.ID
				if out.User.Principal != nil {
					pid := out.User.Principal.ID
					gateway.AddSpanEvents(span, "user.created", map[string]string{"principal.id": fmt.Sprint(pid)})
					user.Principal = &models.Principal{
						ID:        out.User.Principal.ID,
						Admin:     out.User.Principal.Admin,
						LastLogin: out.User.Principal.LastLogin,
						CreatedOn: out.User.Principal.CreatedOn,
						State:     out.User.Principal.State,
					}
				}
			}
			gateway.AddSpanTags(span, map[string]string{"login.phase": fmt.Sprint(out.Phase)})
			return private.NewLoginDiscordOK().WithPayload(&models.AuthPhaseResult{
				NextURL: out.NextURL,
				Token:   out.Token,
				Phase:   out.Phase,
				User:    user,
			})
		},
	)
	api.PrivateCurrentUserHandler = private.CurrentUserHandlerFunc(
		func(cup private.CurrentUserParams) middleware.Responder {
			return private.NewCurrentUserOK().WithPayload(&models.User{})
		},
	)
	api.UseSwaggerUI()
	server := restapi.NewServer(api)
	server.Port = port
	server.ConfigureFlags()
	defer server.Shutdown()
	log.Println("Ready.")
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func requireEnv(name string) string {
	out := os.Getenv(name)
	if out == "" {
		log.Fatalln("variable", name, "is unset")
	}
	return out
}
