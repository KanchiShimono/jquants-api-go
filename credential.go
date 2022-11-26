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
	credentialUserDirName    = ".jquants-api"
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
	cred := new(Credential)

	if path, err := userCredentialPath(); err == nil {
		if c, err := readCredentialFile(path); err == nil {
			cred.overwrite(c)
		}
	}
	if path, err := currentDirCredentialPath(); err == nil {
		if c, err := readCredentialFile(path); err == nil {
			cred.overwrite(c)
		}
	}
	if path, ok := os.LookupEnv(credentialFilePathEnvKey); ok {
		if c, err := readCredentialFile(path); err == nil {
			cred.overwrite(c)
		}
	}

	if val, ok := os.LookupEnv(mailAddressEnvKey); ok {
		cred.MailAddress = val
	}
	if val, ok := os.LookupEnv(passwordEnvKey); ok {
		cred.Password = val
	}
	if val, ok := os.LookupEnv(refreshTokenEnvKey); ok {
		cred.RefreshToken = val
	}

	return cred
}

type RefreshTokenResponse struct {
	RefreshToken string `json:"refreshToken"`
	Message      string `json:"message,omitempty"`
}

type getRefreshTokenParameters struct {
	MailAddress string `json:"mailaddress"`
	Password    string `json:"password"`
}

func GetRefreshToken(c *Credential) (InternalToken, error) {
	now := time.Now()
	endpoint, err := url.JoinPath(endpointBase, "token", "auth_user")
	if err != nil {
		return InternalToken{}, err
	}

	p := getRefreshTokenParameters{
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

	var token RefreshTokenResponse
	if err := DecodeBody(resp.Body, &token); err != nil {
		return InternalToken{}, err
	}
	if token.Message != "" {
		return InternalToken{}, errors.New(token.Message)
	}

	return InternalToken{
		Token:  token.RefreshToken,
		Expiry: now.Add(refreshTokenExpiry),
	}, nil
}

type IDTokenResponse struct {
	IDToken string `json:"idToken"`
	Message string `json:"message,omitempty"`
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

	var token IDTokenResponse
	if err := DecodeBody(resp.Body, &token); err != nil {
		return InternalToken{}, err
	}
	if token.Message != "" {
		return InternalToken{}, errors.New(token.Message)
	}

	return InternalToken{
		Token:  token.IDToken,
		Expiry: now.Add(idTokenExpiry),
	}, nil
}

func userCredentialPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, credentialUserDirName, credentialFileName), nil
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

func readCredentialFile(path string) (*Credential, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c := jquantsConfig{}
	if err := toml.NewDecoder(file).Decode(&c); err != nil {
		return nil, err
	}

	return &c.Credential, nil
}
