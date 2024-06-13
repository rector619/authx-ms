package router

import (
	"fmt"

	"github.com/SineChat/auth-ms/external/request"
	"github.com/SineChat/auth-ms/pkg/controller/auth_model"
	"github.com/SineChat/auth-ms/pkg/middleware"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func AuthModel(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *mongodb.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	auth_model := auth_model.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	modelTypeUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, middleware.AppType))
	{
		modelTypeUrl.POST("/get_user", auth_model.GetUser)
		modelTypeUrl.GET("/get_access_token_by_key/:key", auth_model.GetAccessTokenByKey)
		modelTypeUrl.POST("/validate_on_db", auth_model.ValidateOnDB)
		modelTypeUrl.POST("/validate_authorization", auth_model.ValidateAuthorization)
	}

	return r
}
