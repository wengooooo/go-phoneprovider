package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type SmsMan struct {
	Provider
}

type SmsManNumberDetail struct {
	CountryId     interface{} `json:"country_id"`
	ApplicationId interface{} `json:"application_id"`
	Phone         string      `json:"number"`
	ID            int         `json:"request_id"`
	ErrorCode     string      `json:"error_code"`
	ErrorMSG      string      `json:"error_msg"`

	DialingCode string
	CountryCode string
}

// SMS represents info about an incoming SMS
type SmsManSMSDetail struct {
	CountryId     interface{} `json:"country_id"`
	ApplicationId interface{} `json:"application_id"`
	Code          string      `json:"sms_code"`
	Phone         string      `json:"number"`
	Taskid        interface{} `json:"request_id"`
	ErrorCode     string      `json:"error_code"`
	ErrorMSG      string      `json:"error_msg"`
}

func NewSmsMan(APIKey string) *SmsMan {
	return &SmsMan{
		Provider{
			Endpoint:           "http://api.sms-man.com/control",
			APIKey:             APIKey,
			VerificationMethod: "token",
		},
	}
}

func (f *SmsMan) GetNumber(countryId, applicationId int) (*SmsManNumberDetail, error) {
	// Check if any additional query values could be encapsulated
	// Check if any additional query values could be encapsulated
	queryValues := url.Values{}

	queryValues.Add("country_id", strconv.Itoa(countryId))
	queryValues.Add("application_id", strconv.Itoa(applicationId))

	// Make request
	resp, err := f.makeGetRequest(
		fmt.Sprintf("%s/get-number", f.Endpoint),
		&queryValues,
	)

	fmt.Println(resp.String())
	if err != nil {
		return nil, err
	}

	var info SmsManNumberDetail
	err = json.Unmarshal(resp.Body(), &info)
	if err != nil {
		return nil, err
	}

	if info.ErrorCode != "" {
		return nil, errors.New(info.ErrorMSG)
	}

	countryCode, phone := f.GetCountry(info.Phone)
	info.CountryCode = countryCode
	info.Phone = phone

	return &info, nil
}

func (f *SmsMan) GetSms(ID int) (*SmsManSMSDetail, error) {
	// Make request
	queryValues := url.Values{}
	queryValues.Add("request_id", strconv.Itoa(ID))

	// Make request
	resp, err := f.makeGetRequest(
		fmt.Sprintf("%s/get-sms", f.Endpoint),
		&queryValues,
	)

	if err != nil {
		return nil, err
	}

	var info SmsManSMSDetail
	err = json.Unmarshal(resp.Body(), &info)

	if err != nil {
		return nil, err
	}

	// Check status code
	//if info.ErrorCode != "wait_sms" {
	//	return nil, fmt.Errorf("%s", info.ErrorMSG)
	//}

	return &info, nil

}
