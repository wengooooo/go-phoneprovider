package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type FiveSim struct {
	Provider
}

type NumberDetail struct {
	ID               int         `json:"id"`
	Phone            string      `json:"phone"`
	Operator         string      `json:"operator"`
	Product          string      `json:"product"`
	Price            float32     `json:"price"`
	Status           string      `json:"status"`
	Expires          time.Time   `json:"expires"`
	SMS              []SMSDetail `json:"sms"`
	CreatedAt        time.Time   `json:"created_at"`
	Forwarding       bool        `json:"forwarding"`
	ForwardingNumber string      `json:"forwarding_number"`
	Country          string      `json:"russia"`

	DialingCode string
	CountryCode string
}

// SMS represents info about an incoming SMS
type SMSDetail struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Date      time.Time `json:"date"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Code      string    `json:"code"`
}

func NewFiveSim(APIKey string) *FiveSim {
	return &FiveSim{
		Provider{
			Endpoint:           "https://5sim.net/v1",
			APIKey:             APIKey,
			VerificationMethod: "Authorization",
		},
	}
}

func (f *FiveSim) GetNumber(params map[string]interface{}) (ID int, countryCode, phoneNumber string, err error) {
	// If country is empty, it will pass "any" to the service
	country := params["country"].(string)
	operator := params["operator"].(string)
	forwardingNumber := params["forwardingNumber"].(string)
	name := params["name"].(string)

	if country == "" {
		country = "any"
	}
	// If operator is empty, it will pass "any" to the service
	if operator == "" {
		operator = "any"
	}

	// Check if any additional query values could be encapsulated
	queryValues := url.Values{}
	if forwardingNumber != "" {
		queryValues.Add("forwarding", "1")
		queryValues.Add("number", forwardingNumber)
	}

	// Make request
	resp, err := f.makeGetRequest(
		fmt.Sprintf("%s/user/buy/activation/%s/%s/%s",
			f.Endpoint, country, operator, name,
		),
		&queryValues,
	)

	if err != nil {
		return 0, "", "", err
	}

	if resp.StatusCode() != 200 {
		switch resp.String() {
		case "bad operator":
			return 0, "", "", errors.New("5sim: 错误的运营商代号")
		case "not enough user balance":
			return 0, "", "", errors.New("5sim: 账户余额不足")
		case "not enough rating":
			return 0, "", "", errors.New("5sim: 账户评级过低")
		case "select country":
			return 0, "", "", errors.New("5sim: 没有选择国家")
		case "select operator":
			return 0, "", "", errors.New("5sim: 没有选择运营商")
		case "bad country":
			return 0, "", "", errors.New("5sim: 国家无效")
		case "no product":
			return 0, "", "", errors.New(fmt.Sprintf("5sim: 没有%s的项目", name))
		case "no server offline":
			return 0, "", "", errors.New("5sim: 服务器已下线")
		}
	}

	var info NumberDetail
	err = json.Unmarshal(resp.Body(), &info)
	if err != nil {
		return 0, "", "", err
	}

	countryCode, phoneNumber = f.GetCountry(info.Phone)

	return info.ID, countryCode, phoneNumber, nil
}

func (f *FiveSim) GetSms(orderID int) (code string, err error) {
	// Make request
	resp, err := f.makeGetRequest(
		fmt.Sprintf("%s/user/check/%d",
			f.Endpoint, orderID,
		),
		&url.Values{},
	)

	if err != nil {
		return "", err
	}

	// Check status code
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("%s", resp.Status)
	}

	var info NumberDetail
	err = json.Unmarshal(resp.Body(), &info)
	fmt.Println(resp.String())
	if err != nil {
		return "", err
	}

	if len(info.SMS) > 0 {
		code = info.SMS[0].Code
	}

	return code, nil

}

func (f *FiveSim) ReleaseNumber(ID interface{}) (err error) {
	// Make request
	_, err = f.makeGetRequest(
		fmt.Sprintf("%s/cancel/%d", f.Endpoint, ID),
		&url.Values{},
	)

	return nil
}
