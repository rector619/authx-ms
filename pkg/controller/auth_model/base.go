package auth_model

import (
	"github.com/SineChat/auth-ms/external/request"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Db        *mongodb.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}
