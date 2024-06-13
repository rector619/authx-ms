package auth

import (
	"encoding/json"
	"net/http"

	"github.com/SineChat/auth-ms/internal/config"
	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/services/auth"
	"github.com/SineChat/auth-ms/utility"
	"github.com/gin-gonic/gin"
)

func (base *Controller) FlutterWaveWebhook(c *gin.Context) {
	var (
		req models.FlutterWaveWebhookRequest
	)

	if config.GetConfig().FlutterWave.WebhookSecret != utility.GetHeader(c, "verif-hash") {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid webhook secret", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	requestBody, err := c.GetRawData()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = json.Unmarshal(requestBody, &req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	code, err := auth.FlutterWaveWebhookService(base.Logger, base.Db, req, requestBody)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	code, err = auth.FlutterWaveWebhookVerifyPaymentService(base.Logger, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "ok", nil)
	c.JSON(http.StatusOK, rd)
}
