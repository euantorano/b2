// Package stringutil provides an interface to the Backblaze B2 file storage service.
package b2

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	V1_AUTH_ACCOUNT_URL string = "https://api.backblaze.com/b2api/v1/b2_authorize_account"
)

// Credentials represent the user credentials used to authenticate with the B2 API.
type Credentials struct {
	AccountId      string
	ApplicationKey string
}

// A Client is a client to the B2 API.
type Client struct {
	AccountId   string `json:"accountId"`
	ApiUrl      string `json:"apiUrl"`
	AuthToken   string `json:"authorizationToken"`
	DownloadUrl string `json:"downloadUrl"`

	httpClient *http.Client
}

// An ApiError wraps an error response received from the API.
type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Error implements the `error` interface for `ApiError`, providing the error message that was returned by the API.
func (e ApiError) Error() string {
	return e.Message
}

// NewClient initialises a new API client using the given credentials.
func NewClient(credentials Credentials) (*Client, error) {
	if req, err := http.NewRequest("GET", V1_AUTH_ACCOUNT_URL, nil); err != nil {
		return nil, err
	} else {
		req.SetBasicAuth(credentials.AccountId, credentials.ApplicationKey)

		var c Client
		c.httpClient = &http.Client{}

		err = c.requestJson(req, &c)

		if err != nil {
			return nil, err
		}

		return &c, nil
	}
}

// setHeaders applies the required `Authorization` and `Content-Type` headers on a given HTTP request.
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", c.AuthToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
}

// buildRequestUrl builds the full URL for a given path.
func (c *Client) buildRequestUrl(reqPath string) string {
	return c.ApiUrl + reqPath
}

// buildFileRequestUrl buils the full URL for the given file path.
func (c *Client) buildFileRequestUrl(reqPath string) string {
	return c.DownloadUrl + reqPath
}

// requestJson simplifies the requesting of a JSON response via HTTP, returning an ApiError instance on status codes that don't equal `200 OK`.
func (c *Client) requestJson(req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		// TODO: This is an API error
		var errorResult ApiError

		if json.Unmarshal(body, &errorResult); err != nil {
			return err
		} else {
			return errorResult
		}
	} else {
		return json.Unmarshal(body, result)
	}
}

// requestJson simplifies the requesting of a JSON response via HTTP, returning an ApiError instance on status codes that don't equal `200 OK`.
func (c *Client) requestBytes(req *http.Request) ([]byte, http.Header, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		// TODO: This is an API error
		var errorResult ApiError

		if json.Unmarshal(body, &errorResult); err != nil {
			return nil, nil, err
		} else {
			return nil, nil, errorResult
		}
	} else {
		return body, resp.Header, nil
	}
}
