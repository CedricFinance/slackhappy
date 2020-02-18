package bamboohr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ApiURL = "https://api.bamboohr.com/api/gateway.php"

type Client interface {
	CustomReport(ctx context.Context, title string, fields []string) (*CustomReportResponse, error)
}

type client struct {
	domain     string
	token      string
	httpClient *http.Client
	baseURL    string
}

type Option func(client *client)

func OptionHttpClient(httpClient *http.Client) Option {
	return func(c *client) {
		c.httpClient = httpClient
	}
}

func OptionBaseURL(url string) Option {
	return func(c *client) {
		c.baseURL = url
	}
}

func New(domain string, token string, options ...Option) Client {
	c := &client{
		baseURL:    ApiURL,
		domain:     domain,
		token:      token,
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

type CustomReportRequest struct {
	Title  string   `json:"title"`
	Fields []string `json:"fields"`
}

type CustomReportResponse struct {
	Title     string
	Fields    []string
	Employees []map[string]string
}

func (c *client) CustomReport(ctx context.Context, title string, fields []string) (*CustomReportResponse, error) {
	response, err := c.doPost(ctx, "v1/reports/custom", CustomReportRequest{
		Title:  title,
		Fields: fields,
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 299 {
		return nil, fmt.Errorf("request failed, got status %q", response.Status)
	}

	var customReport CustomReportResponse

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &customReport)
	defer response.Body.Close()

	return &customReport, nil
}

func (c *client) doPost(ctx context.Context, path string, payload interface{}) (*http.Response, error) {
	var buffer bytes.Buffer

	encoder := json.NewEncoder(&buffer)
	encoder.Encode(payload)

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/%s", c.baseURL, c.domain, path), &buffer)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.token, "x")

	request = request.WithContext(ctx)

	return c.httpClient.Do(request)
}
