package streetsweeping

import (
	"log"
	"net/url"

	"github.com/dcu/go-authy"
	"github.com/sfreiberg/gotwilio"
)

// MessageServicer is out interface for sending messages. This can be mocked in the tests.
type MessageServicer interface {
	phoneVerifier
	smsMessager
}

type phoneVerifier interface {
	RequestCode(phoneNumber string) (bool, error)
	VerifyCode(phoneNumber, code string) (bool, error)
}

type smsMessager interface {
	Send(from, to, body string) error
}

type TwilioMessageService struct {
	Authy  *authy.Authy
	Twilio *gotwilio.Twilio
}

func (t *TwilioMessageService) Send(from, to, body string) error {
	_, _, err := t.Twilio.SendSMS("+1"+from, "+1"+to, body, "", "") // todo: should not ignore those returns
	return err
}

func (t *TwilioMessageService) RequestCode(phoneNumber string) (bool, error) {
	log.Println("phoneNumber: ", phoneNumber)
	verification, err := t.Authy.StartPhoneVerification(1, phoneNumber, "sms", url.Values{})
	return verification.Success, err
}

func (t *TwilioMessageService) VerifyCode(phoneNumber, code string) (bool, error) {
	verification, err := t.Authy.CheckPhoneVerification(1, phoneNumber, code, url.Values{})
	return verification.Success, err
}
