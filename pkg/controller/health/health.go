package health

import (
	"net/http"

	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Db        *mongodb.Database
	Validator *validator.Validate
	Logger    *utility.Logger
}

func (base *Controller) Get(c *gin.Context) {
	base.Logger.Info("ping successfull")
	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successful", nil)
	c.JSON(http.StatusOK, rd)

}
