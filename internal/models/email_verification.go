package models

type EmailVerificationRequest struct {
	Email       string `json:"email" validate:"required" mgvalidate:"exists=auth$users$email,email"`
	RedirectURL string `json:"redirect_url" validate:"required"`
}

type EmailVerification struct {
	Token int `json:"token" validate:"required"`
}
