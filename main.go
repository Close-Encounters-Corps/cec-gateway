package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/Close-Encounters-Corps/cec-gateway/docs"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/config"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/controllers"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/gateway"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
 
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var COMMITSHA string

// @title CEC Gateway
// @version 0.1.0

// @Description Gateway endpoint of a CEC Platform v2, serves as a proxy and one swagger to rule them all. Find more at Close Encounters Corps Discord server!
// @BasePath /api

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
	authCtrl := controllers.AuthController{
		Client: app.Client,
		Config: app.Cfg,
	}
	princCtrl := controllers.PrincipalsController{
		Client: app.Client,
		Config: app.Cfg,
	}

	r := gin.Default()
	api0 := r.Group("/api/0")
	api0.Use(otelgin.Middleware("api-0"))
	api0.GET("/login/discord", authCtrl.LoginDiscord)
	api0.GET("/users/current", princCtrl.GetCurrentUser)
	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
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
