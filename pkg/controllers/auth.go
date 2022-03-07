package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Close-Encounters-Corps/cec-gateway/pkg/config"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/gateway"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/models"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/util"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthController struct {
	Config *config.Config
	Client *http.Client
}

// @Router /0/login/discord [get]
// @Summary Authenticate using Discord
// @Tags private
// @Param state query string false "Second phase: State to fetch from CEC Auth"
// @Param redirect_url query string false "URL to redirect after second phase"
// @Success 200 {object} models.AuthPhaseResult
// @Failure 400,500 {object} models.Error
func (ctrl *AuthController) LoginDiscord(c *gin.Context) {
	ctx, span := gateway.NewSpan(c.Request.Context(), "controller.login.discord", nil)
	traceId := span.SpanContext().TraceID().String()
	c.Header("X-Trace-Id", traceId)
	defer span.End()
	internalError := func(err error) {
		gateway.AddSpanError(span, err)
		gateway.FailSpan(span, "internal error")
		log.Printf("[%s] error: %s", traceId, err)
		c.JSON(http.StatusInternalServerError, models.Error{
			RequestID: traceId,
		})
	}
	state := c.Query("state")
	if state != "" {
		gateway.AddSpanTags(span, map[string]string{"request.state": state})
	}
	redirectUrl := ctrl.Config.Urls.External
	if param := c.Query("redirect_url"); param != "" {
		redirectUrl = param
	}
	req, err := util.V1LoginDiscordRequest(
		ctx,
		ctrl.Config.Urls.Core,
		"/v1/login/discord",
		redirectUrl,
		state,
	)
	if err != nil {
		internalError(err)
		return
	}
	span.AddEvent("making request")
	resp, err := ctrl.Client.Do(req)
	if err != nil {
		internalError(err)
		return
	}
	span.AddEvent("request done")
	span.SetAttributes(attribute.Int("response.status", resp.StatusCode))
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		internalError(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		gateway.AddSpanEvents(span, "non-OK response", map[string]string{"responde.body": string(body)})
		if resp.StatusCode == http.StatusBadRequest {
			var message models.Error
			err := json.Unmarshal(body, &message)
			if err != nil {
				internalError(err)
				return
			}
			c.JSON(http.StatusOK, message)
			return
		}
		c.JSON(http.StatusOK, models.Error{RequestID: traceId})
		return
	}
	var out models.AuthPhaseResult
	err = json.Unmarshal(body, &out)
	if err != nil {
		gateway.AddSpanTags(span, map[string]string{"response.body": string(body)})
		internalError(err)
		return
	}
	if out.User != nil {
		if out.User.Principal != nil {
			span.AddEvent("user.created", trace.WithAttributes(
				attribute.Int64("principal.id", out.User.Principal.ID),
			))
		}
	}
	span.SetAttributes(attribute.Int("login.phase", int(out.Phase)))
	c.JSON(http.StatusOK, out)
}
