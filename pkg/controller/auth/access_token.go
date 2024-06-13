package auth

import (
	"fmt"
	"net/http"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/services/auth"
	"github.com/SineChat/auth-ms/utility"
	"github.com/gin-gonic/gin"
)

func (base *Controller) CreateAccessToken(c *gin.Context) {

	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	data, code, err := auth.CreateAccessTokenService(base.Logger, base.Db, user.ID)

	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "access tokens successfully generated", data)
	c.JSON(http.StatusOK, rd)

}
