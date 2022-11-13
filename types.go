package jquants_api_go

import (
	"bytes"
	"fmt"
	"net/url"
	"time"
)

const (
	codeQueryKey = "code"
	dateQueryKey = "date"
)

type JSONTime int64

/*
https://kenzo0107.github.io/2020/05/19/2020-05-20-go-json-time/
*/

// String converts the unix timestamp into a string
func (t JSONTime) String() string {
	tm := t.Time()
	return fmt.Sprintf("\"%s\"", tm.Format("2006-01-02"))
}

// Time returns a `time.Time` representation of this value.
func (t JSONTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}

// UnmarshalJSON will unmarshal both string and int JSON values
func (t *JSONTime) UnmarshalJSON(buf []byte) error {
	s := bytes.Trim(buf, `"`)
	aa, err := time.Parse("20060102", string(s))
	if err != nil {
		return err
	}

	*t = JSONTime(aa.Unix())
	return nil
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken"`
}

type IDToken struct {
	IDToken string `json:"idToken"`
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

type DailyQuotes struct {
	DailyQuotes []Quote `json:"daily_quotes"`
	Message     string  `json:"message,omitempty"`
}

type PricesDailyQuotesOption interface {
	ApplyPricesDailyQuotesOptions(o *PricesDailyQuotesOptions)
}

type PricesDailyQuotesOptions struct {
	Code *string
	Date *string
}

func (o PricesDailyQuotesOptions) Query() url.Values {
	query := url.Values{}
	if o.Code != nil {
		query.Set(codeQueryKey, *o.Code)
	}
	if o.Date != nil {
		query.Set(dateQueryKey, *o.Date)
	}
	return query
}

func (o PricesDailyQuotesOptions) QueryString() string {
	return o.Query().Encode()
}

type code string

func Code(c string) code {
	return code(c)
}

func (c code) ApplyPricesDailyQuotesOptions(o *PricesDailyQuotesOptions) {
	s := string(c)
	o.Code = &s
}

type date string

func Date(d string) date {
	return date(d)
}

func (d date) ApplyPricesDailyQuotesOptions(o *PricesDailyQuotesOptions) {
	s := string(d)
	o.Date = &s
}
