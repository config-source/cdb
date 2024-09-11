package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	token   string
	baseURL string
	client  http.Client
}

type requestSpec struct {
	method string
	url    string
	body   interface{}
	params map[string]string
}

func New(token, baseURL string) *Client {
	return &Client{token: token, baseURL: baseURL, client: http.Client{}}
}

func (c *Client) Do(ctx context.Context, spec requestSpec, output interface{}) (*http.Response, error) {
	if c.token == "" && !strings.HasPrefix(spec.url, "/auth") {
		return nil, errors.New(
			"Unable to determine auth token for CDB instance." +
				"Try setting $CDB_TOKEN or setting up a config file.",
		)
	}

	fullURL := fmt.Sprintf("%s%s", c.baseURL, spec.url)

	body := bytes.NewBuffer(nil)
	if spec.body != nil {
		err := json.NewEncoder(body).Encode(spec.body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, spec.method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accepts", "application/json")
	req.Header.Add("Content-Type", "application/json")

	query := url.Values{}
	for key, value := range spec.params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	httpResp, err := c.client.Do(req)
	if err != nil {
		return httpResp, fmt.Errorf("bad http response: %w", err)
	}
	defer httpResp.Body.Close()

	decoder := json.NewDecoder(httpResp.Body)
	if httpResp.StatusCode >= 400 {
		var errResponse struct {
			Message string `json:"message"`
		}

		err = decoder.Decode(&errResponse)
		if err != nil {
			return nil, fmt.Errorf("unexpected error decoding an error: %w", err)
		}

		err = fmt.Errorf("failure response from API: %s", errResponse.Message)
	}

	if output != nil {
		err = decoder.Decode(&output)
		if err != nil {
			err = fmt.Errorf("error decoding response to output: %w", err)
		}
	}

	return httpResp, err
}
