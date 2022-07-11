package providers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

//go:embed countries.json
var countries []byte

type Provider struct {
	Endpoint string
	APIKey   string
}

func (p *Provider) makeGetRequest(url string, queryValues *url.Values) (*http.Response, error) {
	// Craft the header
	header := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", p.APIKey),
	}

	// Creates a client
	client := &http.Client{}
	// Creates a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, err
	}

	// Incapsulate header elements into the request
	for k, v := range header {
		req.Header.Set(k, v)
	}

	// Encode the query values (if any)
	req.URL.RawQuery = queryValues.Encode()
	return client.Do(req)
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
