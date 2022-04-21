package strava

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type stravaOauthTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Client) RefreshToken() error {
	log.Print("updating access token")

	reqData := url.Values{
		"client_id":     []string{c.clientId},
		"client_secret": []string{c.clientSecret},
		"refresh_token": []string{c.refreshToken},
		"scope":         []string{"activity:read,activity:write"},
		"grant_type":    []string{"refresh_token"},
	}

	req, err := http.NewRequest(http.MethodPost, "https://www.strava.com/oauth/token", strings.NewReader(reqData.Encode()))
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resBody, err := c.call(req, false)
	if err != nil {
		return fmt.Errorf("c.call: %w", err)
	}

	var v stravaOauthTokenResponse
	if err := json.Unmarshal(resBody, &v); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	c.accessToken = v.AccessToken

	return nil
}
