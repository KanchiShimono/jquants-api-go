package jquants_api_go

type code string

var _ GetPricesDailyQuotesOption = code("")

func Code(c string) code {
	return code(c)
}

func (c code) ApplyGetPricesDailyQuotesOptions(o *GetPricesDailyQuotesOptions) {
	s := string(c)
	o.Code = &s
}

type date string

var _ GetPricesDailyQuotesOption = date("")

func Date(d string) date {
	return date(d)
}

func (d date) ApplyGetPricesDailyQuotesOptions(o *GetPricesDailyQuotesOptions) {
	s := string(d)
	o.Date = &s
}
