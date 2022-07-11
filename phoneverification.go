package providers

import (
	"net/http"
	"net/url"
)

type IPhoneVerification interface {
	GetNumber(url string, queryValues *url.Values) (*http.Response, error)
	GetSms(string)
}

type SmsProvider interface {
	FiveSim
}

type PhoneVerification[T SmsProvider] struct {
	Client T
}

func NewPhoneVerification[T SmsProvider]() *PhoneVerification[T] {
	return &PhoneVerification[T]{}
}

func (c *PhoneVerification[T]) GetNumber(country, operator, name, forwardingNumber string) (*NumberDetail, error) {
	return c.Client.GetNumber(country, operator, name, forwardingNumber)
}
