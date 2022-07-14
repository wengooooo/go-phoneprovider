package providers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type IPhoneProvider interface {
	GetNumber(map[string]interface{}) (int, string, string, error)
	GetSms(ID int) (code string, err error)
	ReleaseNumber(ID interface{}) (err error)
}

//go:embed countries.json
var countries []byte

type Provider struct {
	Endpoint           string
	APIKey             string
	VerificationMethod string
}

func (p *Provider) makeGetRequest(url string, queryValues *url.Values) (*Response, error) {
	// Craft the header
	header := map[string]string{}

	switch p.VerificationMethod {
	case "Authorization":
		header = map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", p.APIKey),
		}
	case "secret":
		queryValues.Set("secret_key", p.APIKey)
	case "token":
		queryValues.Set("token", p.APIKey)
	}

	// Creates a client
	client := &http.Client{}
	// Creates a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Incapsulate header elements into the request
	for k, v := range header {
		req.Header.Set(k, v)
	}

	// Encode the query values (if any)
	req.URL.RawQuery = queryValues.Encode()

	resp, err := client.Do(req)

	response := &Response{
		RawResponse: resp,
	}

	defer func(v interface{}) {
		if c, ok := v.(io.Closer); ok {
			func(_ ...interface{}) {}(c.Close())
		}
	}(resp.Body)

	body := resp.Body

	if response.body, err = ioutil.ReadAll(body); err != nil {
		return response, err
	}

	response.size = int64(len(response.body))

	return response, err
}

func (p *Provider) GetCountry(phone string) (conuntryCode, newPhone string) {
	var country map[string]interface{}
	err := json.Unmarshal(countries, &country)
	if err != nil {

	}

	for key, value := range country {
		info := value.(map[string]interface{})
		if strings.Contains(phone, info["code"].(string)) {
			conuntryCode = key
			newPhone = strings.Replace(phone, info["code"].(string), "", -1)
			return
		}
	}

	return "", ""
}
