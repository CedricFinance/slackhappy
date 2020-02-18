package bamboohr

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const ValidToken = "validtoken"
const InvalidToken = "invalidtoken"
const CompanyDomain = "mycompany"

func startWith(str string, prefix string) bool {
	if len(prefix) > len(str) {
		return false
	}

	return str[:len(prefix)] == prefix
}

func extractLogin(authHeader string) (string, error) {
	if !startWith(authHeader, "Basic ") {
		fmt.Println("Invalid header")
		return "", fmt.Errorf("invalid header %q, header does not start with 'Basic '", authHeader)
	}

	parts := strings.Split(authHeader, " ")

	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(parts[1]))
	decodedData, err := ioutil.ReadAll(decoder)
	if err != nil {
		fmt.Println("failed to decode header")
		return "", fmt.Errorf("failed to decode header %q", parts[1])
	}

	userAndPass := strings.Split(string(decodedData), ":")

	return userAndPass[0], nil
}

func validateAuthorizationHeader(r *http.Request) (bool, error) {
	auhtHeader := r.Header.Get("Authorization")
	if auhtHeader == "" {
		return false, fmt.Errorf("authorization header is missing")
	}

	user, err := extractLogin(auhtHeader)
	if err != nil {
		return false, err
	}

	return user == ValidToken, nil
}

func TestClient_CustomReport(t *testing.T) {
	testserver := startFakeServer("{\n  \"title\": \"title\",\n  \"fields\": [\n    { \"id\": \"field1\", \"type\": \"test\", \"name\": \"field1\" }, \n    { \"id\": \"field2\", \"type\": \"test\", \"name\": \"field2\" } \n  ],\n  \"employees\": [\n    { \"id\": \"1\", \"field1\": \"1\", \"field2\": \"2\" }\n  ]\n}")
	c := New(CompanyDomain, ValidToken, OptionHttpClient(testserver.Client()), OptionBaseURL(testserver.URL))

	report, err := c.CustomReport(context.Background(), "title", []string{"field1", "field2"})

	assert.Nilf(t, err, "Unexpected error %q", err)

	if err == nil {
		assert.NotNil(t, report)
		assert.Equal(t, "title", report.Title)
		assert.Equal(t, 2, len(report.Fields))
		assert.Equal(t, 1, len(report.Employees))
		assert.Equal(t, "1", report.Employees[0]["id"])
		assert.Equal(t, "1", report.Employees[0]["field1"])
		assert.Equal(t, "2", report.Employees[0]["field2"])
	}
}

func TestClient_CustomReport_InvalidAuth(t *testing.T) {
	testserver := startFakeServer("")

	c := New(CompanyDomain, InvalidToken, OptionHttpClient(testserver.Client()), OptionBaseURL(testserver.URL))

	report, err := c.CustomReport(context.Background(), "title", []string{"field1", "field2"})

	assert.Nil(t, report)
	assert.NotNil(t, err)
}

func startFakeServer(responsePayload string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validAuth, err := validateAuthorizationHeader(r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(404)
			return
		}
		if !validAuth {
			fmt.Println("invalid token")
			w.WriteHeader(404)
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(responsePayload))
	}))
}
