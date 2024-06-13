package request

import (
	"fmt"

	"github.com/SineChat/auth-ms/external/microservice/notification"
	"github.com/SineChat/auth-ms/external/mocks"
	"github.com/SineChat/auth-ms/external/thirdparty/flutterwave"
	"github.com/SineChat/auth-ms/external/thirdparty/ipstack"
	"github.com/SineChat/auth-ms/internal/config"
	"github.com/SineChat/auth-ms/utility"
)

type ExternalRequest struct {
	Logger *utility.Logger
	Test   bool
}

var (
	JsonDecodeMethod    string = "json"
	PhpSerializerMethod string = "phpserializer"

	// microservice
	SendWelcomeMail       string = "send_welcome_mail"
	SendResetPasswordMail string = "send_reset_password_mail"
	SendVerificationMail  string = "send_verification_mail"

	// third party
	IpstackResolveIp string = "ipstack_resolve_ip"

	// flutterwave
	FlutterWavePaymentRequest string = "flutter_wave_payment_request"
	FlutterWavePaymentVerify  string = "flutter_wave_payment_verify"
)

func (er ExternalRequest) SendExternalRequest(name string, data interface{}) (interface{}, error) {
	var (
		config = config.GetConfig()
	)
	if !er.Test {
		switch name {
		case IpstackResolveIp:
			obj := ipstack.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v", config.IPStack.BaseUrl),
				Method:       "GET",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.IpstackResolveIp()
		case SendWelcomeMail:
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v1/send/%s", config.Microservices.Notification, SendWelcomeMail),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.SendWelcomeMail()
		case SendResetPasswordMail:
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v1/send/%s", config.Microservices.Notification, SendResetPasswordMail),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.SendResetPasswordMail()
		case SendVerificationMail:
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v1/send/%s", config.Microservices.Notification, SendVerificationMail),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.SendVerificationMail()
		case FlutterWavePaymentRequest:
			obj := flutterwave.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v3/payments", config.FlutterWave.BaseUrl),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.FlutterWavePayment()
		case FlutterWavePaymentVerify:
			obj := flutterwave.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v3/transactions/verify_by_reference?tx_ref=%s", config.FlutterWave.BaseUrl, data),
				Method:       "GET",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.FlutterWaveVerifyPaymentByTxRef()
		default:
			return nil, fmt.Errorf("request not found")
		}

	} else {
		mer := mocks.ExternalRequest{Logger: er.Logger, Test: true}
		return mer.SendExternalRequest(name, data)
	}
}
