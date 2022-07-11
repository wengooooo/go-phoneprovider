package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

type FiveSim struct {
	Provider
}

type NumberDetail struct {
	Taskid           int       `json:"id"`
	Phone            string    `json:"phone"`
	Operator         string    `json:"operator"`
	Product          string    `json:"product"`
	Price            float32   `json:"price"`
	Status           string    `json:"status"`
	Expires          time.Time `json:"expires"`
	SMS              []SMS     `json:"sms"`
	Code             string
	CreatedAt        time.Time `json:"created_at"`
	Forwarding       bool      `json:"forwarding"`
	ForwardingNumber string    `json:"forwarding_number"`
	Country          string    `json:"russia"`
}

// SMS represents info about an incoming SMS
type SMS struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Date      time.Time `json:"date"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Code      string    `json:"code"`
}

func NewFiveSim() *FiveSim {
	return &FiveSim{
		Provider{
			Endpoint: "",
			APIKey:   "",
		},
	}
}

func (f *FiveSim) GetNumber(country, operator, name, forwardingNumber string) (*NumberDetail, error) {
	// If country is empty, it will pass "any" to the service
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
		return &NumberDetail{}, err
	}
	// Check status code
	if resp.StatusCode != 200 {
		return &NumberDetail{}, fmt.Errorf("%s", resp.Status)
	}

	// Read request body
	r, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		resp.Body.Close()
		return &NumberDetail{}, err
	}
	resp.Body.Close()

	// Unmarshal the body into a struct
	var info NumberDetail
	err = json.Unmarshal(r, &info)
	if err != nil {
		return &NumberDetail{}, err
	}

	countryCode, phone := f.GetCountry(info.Phone)
	info.Country = countryCode
	info.Phone = phone

	return &info, nil
}
