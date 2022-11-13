package jquants_api_go

import (
	"fmt"
	"net/http"
)

const (
	authorizationHeaderKey = "Authorization"
	tokenType              = "Bearer"
)

type Transport struct {
	Base   http.RoundTripper
	Source TokenSource
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get(authorizationHeaderKey)) != 0 {
		return t.base().RoundTrip(req)
	}

	token, err := t.Source.Token()
	if err != nil {
		return nil, err
	}

	clonedReq := cloneRequest(req)
	clonedReq.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", tokenType, token.IDToken.Token))

	return t.base().RoundTrip(clonedReq)
}

func (t *Transport) base() http.RoundTripper {
	if t.Base == nil {
		return http.DefaultTransport
	}
	return t.Base
}

func cloneRequest(req *http.Request) *http.Request {
	clone := new(http.Request)
	*clone = *req
	clone.Header = make(http.Header, len(req.Header))
	for k, v := range req.Header {
		clone.Header[k] = append([]string(nil), v...)
	}
	return clone
}
