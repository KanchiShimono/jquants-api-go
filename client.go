package jquants_api_go

import (
	"context"
	"encoding/json"
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

func (c *Client) GetPricesDailyQuotes(ctx context.Context, opts ...PricesDailyQuotesOption) (DailyQuotes, error) {
	baseURL, err := c.Config.BaseURL()
	if err != nil {
		return DailyQuotes{}, err
	}

	endpoint := baseURL.JoinPath("prices", "daily_quotes")
	options := PricesDailyQuotesOptions{}
	for _, o := range opts {
		o.ApplyPricesDailyQuotesOptions(&options)
	}
	endpoint.RawQuery = options.QueryString()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return DailyQuotes{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return DailyQuotes{}, err
	}
	defer resp.Body.Close()

	var quotes DailyQuotes
	if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
		return DailyQuotes{}, err
	}

	return quotes, nil
}
