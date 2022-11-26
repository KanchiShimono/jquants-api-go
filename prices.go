package jquants_api_go

import (
	"context"
	"net/http"
	"net/url"
)

type GetPricesDailyQuotesOption interface {
	ApplyGetPricesDailyQuotesOptions(o *GetPricesDailyQuotesOptions)
}

type GetPricesDailyQuotesOptions struct {
	Code *string
	Date *string
}

func (o GetPricesDailyQuotesOptions) Query() url.Values {
	query := url.Values{}
	if o.Code != nil {
		query.Set(codeQueryKey, *o.Code)
	}
	if o.Date != nil {
		query.Set(dateQueryKey, *o.Date)
	}
	return query
}

func (o GetPricesDailyQuotesOptions) QueryString() string {
	return o.Query().Encode()
}

type Quote struct {
	Code             string   `json:"Code"`
	Close            float64  `json:"Close"`
	Date             JSONTime `json:"Date"`
	AdjustmentHigh   float64  `json:"AdjustmentHigh"`
	Volume           float64  `json:"Volume"`
	TurnoverValue    float64  `json:"TurnoverValue"`
	AdjustmentClose  float64  `json:"AdjustmentClose"`
	AdjustmentLow    float64  `json:"AdjustmentLow"`
	Low              float64  `json:"Low"`
	High             float64  `json:"High"`
	Open             float64  `json:"Open"`
	AdjustmentOpen   float64  `json:"AdjustmentOpen"`
	AdjustmentFactor float64  `json:"AdjustmentFactor"`
	AdjustmentVolume float64  `json:"AdjustmentVolume"`
}

type GetPricesDailyQuotesResponse struct {
	DailyQuotes []Quote `json:"daily_quotes"`
	Message     string  `json:"message,omitempty"`
}

func (c *Client) GetPricesDailyQuotes(ctx context.Context, opts ...GetPricesDailyQuotesOption) (*GetPricesDailyQuotesResponse, error) {
	baseURL, err := c.Config.BaseURL()
	if err != nil {
		return nil, err
	}

	endpoint := baseURL.JoinPath("prices", "daily_quotes")
	options := GetPricesDailyQuotesOptions{}
	for _, o := range opts {
		o.ApplyGetPricesDailyQuotesOptions(&options)
	}
	endpoint.RawQuery = options.QueryString()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	quotes := new(GetPricesDailyQuotesResponse)
	if err := DecodeBody(resp.Body, quotes); err != nil {
		return nil, err
	}

	return quotes, nil
}
