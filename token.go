package jquants_api_go

import (
	"errors"
	"sync"
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

type tokenRefresher struct {
	t *Token
}

func (s *tokenRefresher) Token() (*Token, error) {
	if s.t == nil || s.t.RefreshToken.Token == "" {
		return nil, errors.New("refresh token is not set")
	}

	cred := &Credential{RefreshToken: s.t.RefreshToken.Token}
	idToken, err := GetIDToken(cred)
	if err != nil {
		return nil, err
	}

	return &Token{
		IDToken:      idToken,
		RefreshToken: s.t.RefreshToken,
	}, nil
}

type reuseTokenSource struct {
	mu  sync.Mutex
	new TokenSource
	t   *Token
}

func (s *reuseTokenSource) Token() (*Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.t.Valid() {
		return s.t, nil
	}

	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}

	s.t = t
	return t, nil
}
