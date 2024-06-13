package notification

import (
	"fmt"

	"github.com/SineChat/auth-ms/external/external_models"
)

func (r *RequestObj) SendWelcomeMail() (interface{}, error) {

	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.SendWelcomeMail)
	if !ok {
		return nil, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	logger.Info("welcome email", data)

	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("welcome email", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("welcome email", outBoundResponse)

	return nil, nil

}
