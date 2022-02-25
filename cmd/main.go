package main

import (
	"encoding/json"
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
	"github.com/Close-Encounters-Corps/cec-gateway/gen/restapi/operations/auth"
	"github.com/Close-Encounters-Corps/cec-gateway/gen/restapi/operations/private"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/config"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

var COMMITSHA string

func main() {
	log.Println("Commit:", COMMITSHA)
	cfg := config.Config{
		Urls: &config.UrlSet{
			Core: requireEnv("CEC_URLS_CORE"),
		},
	}
	listenport := requireEnv("CEC_LISTENPORT")
	port, err := strconv.Atoi(listenport)
	if err != nil {
		log.Fatalln(err)
	}
	swagger, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	api := operations.NewCecGwAPI(swagger)
	api.AuthLoginDiscordHandler = auth.LoginDiscordHandlerFunc(
		func(ldp auth.LoginDiscordParams) middleware.Responder {
			state := swag.StringValue(ldp.State)
			u, err := url.Parse(cfg.Urls.Core)
			if err != nil {
				log.Println(err)
				return auth.NewLoginDiscordInternalServerError()
			}
			u.Query().Add("state", state)
			req, err := http.NewRequest("GET", u.String(), nil)
			if err != nil {
				log.Println(err)
				return auth.NewLoginDiscordInternalServerError()
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return auth.NewLoginDiscordInternalServerError()
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return auth.NewLoginDiscordInternalServerError()
			}
			var out core_models.AuthPhaseResult
			err = json.Unmarshal(body, &out)
			if err != nil {
				log.Println(err)
				return auth.NewLoginDiscordInternalServerError()
			}
			return auth.NewLoginDiscordOK().WithPayload(&models.AuthPhaseResult{
				NextURL: out.NextURL,
				Token:   out.Token,
				Phase:   out.Phase,
				User: &models.User{
					ID: out.User.ID,
					Principal: &models.Principal{
						ID:        out.User.Principal.ID,
						Admin:     out.User.Principal.Admin,
						LastLogin: out.User.Principal.LastLogin,
						CreatedOn: out.User.Principal.CreatedOn,
						State:     out.User.Principal.State,
					},
				},
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
