package models

type ValidateOnDBReq struct {
	Table string                 `validate:"required" json:"table"`
	Type  string                 `validate:"required" json:"type"`
	Query map[string]interface{} `validate:"required" json:"query"`
}

type ValidateAuthorizationReq struct {
	Type               string `validate:"required" json:"type"`
	AuthorizationToken string `json:"authorization-token"`
	AppKey             string `json:"app-key"`
	PrivateKey         string `json:"private-key"`
	PublicKey          string `json:"public-key"`
}
