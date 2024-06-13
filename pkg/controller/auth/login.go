package auth

import (
	"fmt"
	"net/http"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/services/auth"
	"github.com/SineChat/auth-ms/utility"
	"github.com/gin-gonic/gin"
)

func (base *Controller) Login(c *gin.Context) {
	var (
		req = models.LoginRequest{}
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

	data, code, err := auth.LoginService(base.Logger, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "login successful", data)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) Logout(c *gin.Context) {
	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	token := models.IdentityLoginToken
	if token == nil {
		msg := "error retrieving authentication token"
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	err := token.DeleteLoginToken(base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "logout successful", nil)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) LogoutAllSessions(c *gin.Context) {
	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	token := models.IdentityLoginToken
	if token == nil {
		msg := "error retrieving authentication token"
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	err := token.DeleteLoginTokensByAccountID(base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "logout successful", nil)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) ChangeAuthUserPassword(c *gin.Context) {
	var (
		req = models.ChangeAuthUserPassword{}
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

	code, err := auth.ChangeAuthUserPasswordService(base.Logger, base.Db, req, user.Email)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "password changed successfully", nil)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) AdminLogin(c *gin.Context) {
	var (
		req = models.LoginRequest{}
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

	data, code, err := auth.AdminLoginService(base.Logger, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "login successful", data)
	c.JSON(http.StatusOK, rd)
}
