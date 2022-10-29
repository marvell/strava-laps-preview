package strava

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	authScope       = "activity:read,activity:write"
	authRedirectUrl = "http://localhost"
)

type stravaOauthTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Client) oauthToken(reqData url.Values) (*stravaOauthTokenResponse, error) {
	reqData["client_id"] = []string{c.clientId}
	reqData["client_secret"] = []string{c.clientSecret}

	req, err := http.NewRequest(http.MethodPost, "https://www.strava.com/oauth/token", strings.NewReader(reqData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resBody, err := c.call(req, false)
	if err != nil {
		return nil, fmt.Errorf("c.call: %w", err)
	}

	var v stravaOauthTokenResponse
	if err := json.Unmarshal(resBody, &v); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &v, nil
}

func (c *Client) AccessToken() (string, error) {
	log.Print("get access token")

	reqData := url.Values{
		"refresh_token": []string{c.refreshToken},
		"scope":         []string{authScope},
		"grant_type":    []string{"refresh_token"},
	}

	v, err := c.oauthToken(reqData)
	if err != nil {
		return "", err
	}

	return v.AccessToken, nil
}

func (c *Client) AuthorizeUrl() (*url.URL, error) {
	u, err := url.Parse("https://www.strava.com/oauth/authorize")
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	u.RawQuery = url.Values{
		"client_id":     []string{c.clientId},
		"redirect_uri":  []string{authRedirectUrl},
		"response_type": []string{"code"},
		"scope":         []string{authScope},
	}.Encode()

	return u, nil
}

func (c *Client) RefreshToken(code string) (string, error) {
	log.Print("get refresh token")

	reqData := url.Values{
		"grant_type": []string{"authorization_code"},
		"code":       []string{code},
	}

	v, err := c.oauthToken(reqData)
	if err != nil {
		return "", err
	}

	return v.RefreshToken, nil
}
