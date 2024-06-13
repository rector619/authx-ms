package router

import (
	"fmt"

	"github.com/SineChat/auth-ms/pkg/controller/health"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Health(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *mongodb.Database, logger *utility.Logger) *gin.Engine {
	health := health.Controller{Db: db, Validator: validator, Logger: logger}

	healthUrl := r.Group(fmt.Sprintf("%v/", ApiVersion))
	{
		healthUrl.GET("/health", health.Get)
	}
	return r
}
