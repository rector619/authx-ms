package router

import (
	"fmt"

	"github.com/SineChat/auth-ms/external/request"
	"github.com/SineChat/auth-ms/pkg/controller/auth"
	"github.com/SineChat/auth-ms/pkg/middleware"
	"github.com/SineChat/auth-ms/pkg/repository/storage/mongodb"
	"github.com/SineChat/auth-ms/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Auth(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *mongodb.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	authC := auth.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	authUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		authUrl.POST("/signup", authC.Signup)
		authUrl.POST("/login", authC.Login)
		authUrl.POST("/admin/login", authC.AdminLogin)
	}

	authTypeUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, middleware.AuthType))
	{
		authTypeUrl.POST("/logout", authC.Logout)
		authTypeUrl.POST("/logout_all_sessions", authC.LogoutAllSessions)
		authTypeUrl.POST("/password/change", authC.ChangeAuthUserPassword)
		authTypeUrl.POST("/api-key/create", authC.CreateAccessToken)
		authTypeUrl.POST("/user/update", authC.UpdateUser)
		authTypeUrl.GET("/subscriptions", authC.GetActiveSubscriptions)
	}

	// admin group
	adminUrl := r.Group(fmt.Sprintf("%v/admin", ApiVersion), middleware.Authorize(db, middleware.AdminType))
	{
		adminUrl.POST("/flag", authC.FlagAccount)
		adminUrl.PUT("/unflag", authC.UnFlagAccount)
		adminUrl.GET("/flagged-accounts", authC.GetFlaggedAccounts)
		adminUrl.DELETE("/unban", authC.UnBanAccount)

		// subscription route
		adminUrl.POST("/subscription/create", authC.CreateSubscription)
		adminUrl.GET("/subscription/get", authC.GetSubscriptions)
		adminUrl.PUT("/subscription/update/:id", authC.UpdateSubscription)
		adminUrl.DELETE("/subscription/delete/:id", authC.DeleteSubscription)
	}

	// password group
	pUrl := r.Group(fmt.Sprintf("%v/password", ApiVersion))
	{
		pUrl.POST("/request-reset", authC.PasswordResetRequest)
		pUrl.POST("/reset", authC.PasswordReset)
	}

	// email verification group
	emailUrl := r.Group(fmt.Sprintf("%v/email", ApiVersion))
	{
		emailUrl.POST("/verify", authC.EmailVerification)
		emailUrl.POST("/request-verification", authC.EmailVerificationRequest)
	}

	// payment group
	paymentUrl := r.Group(fmt.Sprintf("%v/payment", ApiVersion))
	{
		paymentUrl.POST("/request", authC.PaymentRequest, middleware.Authorize(db, middleware.AuthType))
		paymentUrl.GET("/verify/:tx_ref", authC.PaymentVerify)
		paymentUrl.POST("/webhook/flutterwave", authC.FlutterWaveWebhook)
	}

	return r
}
