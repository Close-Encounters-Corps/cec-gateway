package gateway

import (
	"context"
	"net/http"

	"github.com/Close-Encounters-Corps/cec-gateway/pkg/config"
)

var APPLICATION = "cec-gw"
var VERSION = "0.1.0"

type Gateway struct {
	Cfg         *config.Config
	Tracer      *Tracer
	Environment string
	Client      *http.Client
}

func (g *Gateway) SetupAll() error {
	var err error
	err = g.SetupTracing(&TracerConfig{
		ServiceName: APPLICATION,
		ServiceVer:  VERSION,
		Jaeger:      g.Cfg.JaegerUrl,
		Environment: g.Environment,
		Disabled:    false,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *Gateway) Close(ctx context.Context) {
	g.Tracer.Close(ctx)
}
