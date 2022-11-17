package jquants_api_go

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

const (
	credentialFileName       = "jquants-api.toml"
	credentialUserDir        = ".jquants-api"
	credentialFilePathEnvKey = "JQUANTS_API_CLIENT_CONFIG_FILE"
	mailAddressEnvKey        = "JQUANTS_API_MAIL_ADDRESS"
	passwordEnvKey           = "JQUANTS_API_PASSWORD"
	refreshTokenEnvKey       = "JQUANTS_API_REFRESH_TOKEN"
	jsonContentType          = "application/json"
	refreshTokenQueryKey     = "refreshtoken"
	refreshTokenExpiry       = 7 * 24 * time.Hour
	idTokenExpiry            = 24 * time.Hour
)

type Credential struct {
	MailAddress  string `toml:"mail_address"`
	Password     string `toml:"password"`
	RefreshToken string `toml:"refresh_token"`
}

func (c *Credential) overwrite(o *Credential) {
	if o == nil {
		return
	}
	if o.MailAddress != "" {
		c.MailAddress = o.MailAddress
	}
	if o.Password != "" {
		c.Password = o.Password
	}
	if o.RefreshToken != "" {
		c.RefreshToken = o.RefreshToken
	}
}

func LoadCredential() *Credential {
	cred := &Credential{}

	if path, err := userCredentialPath(); err == nil {
		c := readCredentialFile(path)
		cred.overwrite(c)
	}
	if path, err := currentDirCredentialPath(); err == nil {
		c := readCredentialFile(path)
		cred.overwrite(c)
	}
	if path, ok := os.LookupEnv(credentialFilePathEnvKey); ok && path != "" {
		c := readCredentialFile(path)
		cred.overwrite(c)
	}

	if val, ok := os.LookupEnv(mailAddressEnvKey); ok && val != "" {
		cred.MailAddress = val
	}
	if val, ok := os.LookupEnv(passwordEnvKey); ok && val != "" {
		cred.Password = val
	}
	if val, ok := os.LookupEnv(refreshTokenEnvKey); ok && val != "" {
		cred.RefreshToken = val
	}

	return cred
}

func GetRefreshToken(c *Credential) (InternalToken, error) {
	if c == nil {
		return InternalToken{}, errors.New("credential is nil")
	}

	now := time.Now()
	endpoint, err := url.JoinPath(endpointBase, "token", "auth_user")
	if err != nil {
		return InternalToken{}, err
	}

	type params struct {
		MailAddress string `json:"mailaddress"`
		Password    string `json:"password"`
	}

	p := params{
		MailAddress: c.MailAddress,
		Password:    c.Password,
	}

	b, err := json.Marshal(p)
	if err != nil {
		return InternalToken{}, err
	}

	resp, err := http.Post(endpoint, jsonContentType, bytes.NewBuffer(b))
	if err != nil {
		return InternalToken{}, err
	}
	defer resp.Body.Close()

	var rt RefreshToken
	if err := json.NewDecoder(resp.Body).Decode(&rt); err != nil {
		return InternalToken{}, err
	}
	if rt.Message != "" {
		return InternalToken{}, errors.New(rt.Message)
	}

	return InternalToken{
		Token:  rt.RefreshToken,
		Expiry: now.Add(refreshTokenExpiry),
	}, nil
}

func GetIDToken(c *Credential) (InternalToken, error) {
	if c == nil {
		return InternalToken{}, errors.New("credential is nil")
	}

	now := time.Now()
	baseURL, err := url.Parse(endpointBase)
	if err != nil {
		return InternalToken{}, err
	}

	endpoint := baseURL.JoinPath("token", "auth_refresh")
	query := endpoint.Query()
	query.Set(refreshTokenQueryKey, c.RefreshToken)
	endpoint.RawQuery = query.Encode()

	resp, err := http.Post(endpoint.String(), jsonContentType, nil)
	if err != nil {
		return InternalToken{}, err
	}
	defer resp.Body.Close()

	var t IDToken
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return InternalToken{}, err
	}
	if t.Message != "" {
		return InternalToken{}, errors.New(t.Message)
	}

	return InternalToken{
		Token:  t.IDToken,
		Expiry: now.Add(idTokenExpiry),
	}, nil
}

func userCredentialPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, credentialUserDir, credentialFileName), nil
}

func currentDirCredentialPath() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(current, credentialFileName), nil
}

type jquantsConfig struct {
	Credential Credential `toml:"jquants-api-client"`
}

func readCredentialFile(path string) *Credential {
	file, err := os.Open(path)
	if err != nil {
		return &Credential{}
	}
	defer file.Close()

	c := jquantsConfig{}
	if err := toml.NewDecoder(file).Decode(&c); err != nil {
		return &Credential{}
	}

	return &c.Credential
}
