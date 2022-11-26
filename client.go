package jquants_api_go

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	endpointBase = "https://api.jpx-jquants.com/v1"
)

type Endpoint struct {
	URL string
}

type Config struct {
	Base Endpoint
}

func NewConfig() Config {
	return Config{
		Base: Endpoint{
			URL: endpointBase,
		},
	}
}

func (c Config) BaseURL() (*url.URL, error) {
	return url.Parse(c.Base.URL)
}

type Client struct {
	httpClient *http.Client
	Config     Config
}

func NewClient(c Config, s TokenSource) *Client {
	httpClient := &http.Client{
		Transport: &Transport{Source: s},
	}

	return &Client{
		httpClient: httpClient,
		Config:     c,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func DecodeBody(body io.Reader, v interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
