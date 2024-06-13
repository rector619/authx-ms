package external_models

type FlutterWavePaymentRequest struct {
	Amount      string `json:"amount" validate:"required"`
	Currency    string `json:"currency"` // make the defalut be NGN if none is provided
	RedirectURL string `json:"redirect_url" validate:"required"`
	Customer    struct {
		Email    string `json:"email"`
		FullName string `json:"name"`
	} `json:"customer"`
	TxRef string `json:"tx_ref"`
}

type FlutterWavePaymentResponse struct {
	Data struct {
		Link string `json:"link"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type FlutterWaveVerifyPaymentResponse struct {
	Data struct {
		ID            int    `json:"id"`
		TxRef         string `json:"tx_ref"`
		Amount        int    `json:"amount"`
		Currency      string `json:"currency"`
		ChargedAmount int    `json:"charged_amount"`
		IP            string `json:"ip"`
		Narration     string `json:"narration"`
		Status        string `json:"status"`
		PaymentType   string `json:"payment_type"`
	} `json:"data"`
}
