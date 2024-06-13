package models

type FlutterWaveWebhookRequest struct {
	ID               int                               `json:"id"`
	TxRef            string                            `json:"txRef"`
	FlwRef           string                            `json:"flwRef"`
	OrderRef         string                            `json:"orderRef"`
	PaymentPlan      interface{}                       `json:"paymentPlan"`
	PaymentPage      interface{}                       `json:"paymentPage"`
	CreatedAt        string                            `json:"createdAt"`
	Amount           int                               `json:"amount"`
	ChargedAmount    int                               `json:"charged_amount"`
	Status           string                            `json:"status"`
	IP               string                            `json:"IP"`
	Currency         string                            `json:"currency"`
	Appfee           float64                           `json:"appfee"`
	Merchantfee      int                               `json:"merchantfee"`
	Merchantbearsfee int                               `json:"merchantbearsfee"`
	Customer         FlutterWaveWebhookRequestCustomer `json:"customer"`
	Entity           interface{}                       `json:"entity"`
	EventType        string                            `json:"event.type"`
}

type FlutterWaveWebhookRequestCustomer struct {
	ID            int         `json:"id"`
	Phone         interface{} `json:"phone"`
	FullName      string      `json:"fullName"`
	Customertoken interface{} `json:"customertoken"`
	Email         string      `json:"email"`
	CreatedAt     string      `json:"createdAt"`
	UpdatedAt     string      `json:"updatedAt"`
	AccountID     int         `json:"AccountId"`
}
