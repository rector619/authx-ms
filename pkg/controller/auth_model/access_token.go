package auth_model

import (
	"net/http"

	"github.com/SineChat/auth-ms/internal/models"
	"github.com/SineChat/auth-ms/utility"
	"github.com/gin-gonic/gin"
)

func (base *Controller) GetAccessTokenByKey(c *gin.Context) {
	var (
		key         = c.Param("key")
		accessToken = models.AccessToken{PrivateKey: key, PublicKey: key, IsLive: true}
	)

	err := accessToken.LiveTokensWithKey(base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", accessToken)
	c.JSON(http.StatusOK, rd)
}
