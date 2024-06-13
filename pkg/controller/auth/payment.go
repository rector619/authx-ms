package auth

import (
	"fmt"
	"net/http"

	"github.com/SineChat/auth-ms/external/external_models"
	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/services/auth"
	"github.com/SineChat/auth-ms/utility"
	"github.com/gin-gonic/gin"
)

func (base *Controller) PaymentRequest(c *gin.Context) {
	var (
		req = external_models.FlutterWavePaymentRequest{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	vr := mongodb.ValidateRequestM{Logger: base.Logger, Test: false}
	err = vr.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	response, code, err := auth.PaymentRequestService(base.Logger, base.Db, req, user.Email)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	data, ok := response.(external_models.FlutterWavePaymentResponse)
	if !ok {
		rd := utility.BuildErrorResponse(code, "error", "response data format error", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, data.Message, data.Data)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) PaymentVerify(c *gin.Context) {

	txRef := c.Param("tx_ref")

	data, code, err := auth.PaymentVerifyService(base.Logger, base.Db, txRef)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	resp, ok := data.(external_models.FlutterWaveVerifyPaymentResponse)
	if !ok {
		rd := utility.BuildErrorResponse(code, "error", "response data format error", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	code, err = auth.UpdatePaymentStatus(base.Logger, base.Db, resp)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	// rd := utility.BuildSuccessResponse(http.StatusOK, "", "payment has been verified")
	rd := utility.BuildSuccessResponse(http.StatusOK, "", fmt.Sprintf("payment status is %s", resp.Data.Status))
	c.JSON(http.StatusOK, rd)
}
