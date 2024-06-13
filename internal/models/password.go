package models

type PasswordResetRequest struct {
	Email       string `json:"email" validate:"required" mgvalidate:"exists=auth$users$email,email"`
	RedirectURL string `json:"redirect_url" validate:"required"`
}

type PasswordReset struct {
	Password string `json:"password" validate:"required"`
	Token    int    `json:"token" validate:"required"`
}

type ChangeAuthUserPassword struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}
