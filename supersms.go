package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type SuperSms struct {
	Provider
}

type SuperSmsNumberDetail struct {
	ID      int         `json:"taskid"`
	PID     interface{} `json:"pid"`
	Phone   string      `json:"phone"`
	Cost    int         `json:"cost"`
	Message string      `json:"message"`

	DialingCode string
	CountryCode string
}

// SMS represents info about an incoming SMS
type SuperSmsSMSDetail struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Date      time.Time `json:"date"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
}

func NewSuperSms(APIKey string) *SuperSms {
	return &SuperSms{
		Provider{
			Endpoint:           "https://www.supersms.ml/api",
			APIKey:             APIKey,
			VerificationMethod: "secret",
		},
	}
}

func (f *SuperSms) GetNumber(channel, country, pid string) (*SuperSmsNumberDetail, error) {
	// Check if any additional query values could be encapsulated
	queryValues := url.Values{}

	queryValues.Add("channel", channel)
	queryValues.Add("country", country)
	queryValues.Add("pid", pid)
	// Make request
	resp, err := f.makeGetRequest(
		fmt.Sprintf("%s/getnumber", f.Endpoint),
		&queryValues,
	)

	if err != nil {
		return nil, err
	}

	var info SuperSmsNumberDetail
	err = json.Unmarshal(resp.Body(), &info)
	if err != nil {
		return nil, err
	}

	if info.Message != "" {
		return nil, errors.New(resp.String())
	}

	countryCode, phone := f.GetCountry(info.Phone)
	info.CountryCode = countryCode
	info.Phone = phone

	return &info, nil
}

func (f *SuperSms) GetSms(ID int) (*SuperSmsSMSDetail, error) {
	// Make request
	queryValues := url.Values{}
	queryValues.Add("taskid", strconv.Itoa(ID))

	// Make request
	resp, err := f.makeGetRequest(
		fmt.Sprintf("%s/getcode", f.Endpoint),
		&queryValues,
	)

	if err != nil {
		return nil, err
	}

	// Check status code
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("%d", resp.StatusCode())
	}

	var info SuperSmsSMSDetail
	err = json.Unmarshal(resp.Body(), &info)

	if err != nil {
		return nil, err
	}

	return &info, nil

}
