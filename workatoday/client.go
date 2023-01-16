package workatoday

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WorkatodayClient interface {
	GetWorkersContext(ctx context.Context, requestType string) (*WorkersResponse, error)
}

type Client struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

type Worker struct {
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	HireDate   string `json:"hire_date,omitempty"`
	Birthday   string `json:"birthday,omitempty"`
	Department string `json:"department,omitempty"`
	Active     string `json:"is_active,omitempty"`
}

type WorkersResponse struct {
	Records []Worker
}

func (c *Client) GetWorkersContext(ctx context.Context, requestType string) (*WorkersResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/get_workers?request_type=%s", c.BaseURL, requestType), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the http request: %v", err)
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("API-Token", c.APIKey)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make the http request: %v", err)
	}

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("bad http status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}
	defer res.Body.Close()

	var response WorkersResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the response: %v", err)
	}

	return &response, nil
}
