package jquants_api_go

import (
	"errors"
	"time"
)

const expiryDelta = 1 * time.Hour

type Token struct {
	IDToken      InternalToken
	RefreshToken InternalToken
}

type InternalToken struct {
	Token  string
	Expiry time.Time
}

func (t *Token) Valid() bool {
	return t != nil && t.IDToken.Token != "" && !t.expired()
}

func (t *Token) expired() bool {
	if t.IDToken.Expiry.IsZero() {
		return false
	}
	return t.IDToken.Expiry.Round(0).Add(-expiryDelta).Before(time.Now())
}

type TokenSource interface {
	Token() (*Token, error)
}

type StaticTokenSource struct {
	t *Token
}

func NewStaticTokenSource(t *Token) TokenSource {
	return &StaticTokenSource{t: t}
}

func (s *StaticTokenSource) Token() (*Token, error) {
	if s.t == nil || s.t.IDToken.Token == "" {
		return nil, errors.New("token is not set")
	}
	return s.t, nil
}
