package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/Close-Encounters-Corps/cec-gateway/pkg/config"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/gateway"
	"github.com/Close-Encounters-Corps/cec-gateway/pkg/models"
	"github.com/gin-gonic/gin"
)

type PrincipalsController struct {
	Client *http.Client
	Config *config.Config
}

// @Router /0/users/current [get]
// @Summary Get current user
// @Tags private
// @Param X-Auth-Token header string true "Authorization token"
// @Success 200 {object} models.User
// @Failure 400,500 {object} models.Error
func (ctrl *PrincipalsController) GetCurrentUser(c *gin.Context) {
	ctx, span := gateway.NewSpan(c.Request.Context(), "controller.users.current", nil)
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
	u, err := url.Parse(ctrl.Config.Urls.Core)
	if err != nil {
		internalError(err)
		return
	}
	u.Path = "/v1/users/current"
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		internalError(err)
		return
	}
	req.Header.Add("X-Auth-Token", c.GetHeader("X-Auth-Token"))
	resp, err := ctrl.Client.Do(req)
	if err != nil {
		internalError(err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		internalError(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.Writer.Write(body)
			return
		}
		internalError(err)
		return
	}
	var out models.User
	err = json.Unmarshal(body, &out)
	if err != nil {
		internalError(err)
		return
	}
	c.JSON(http.StatusOK, out)
}
