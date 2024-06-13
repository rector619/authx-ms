package external_models

type SendResetPasswordMail struct {
	Email       string `json:"email"`
	RedirectURL string `json:"redirect_url"`
}

type SendWelcomeMail struct {
	Email string `json:"email"`
}

type SendVerificationMail struct {
	Email       string `json:"email"`
	RedirectURL string `json:"redirect_url"`
}
