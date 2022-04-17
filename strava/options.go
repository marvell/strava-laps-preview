package strava

import (
	"fmt"
	"net/http"

	"golang.org/x/net/proxy"
)

type Option func(*Client) error

func WithDebugMode(enabled bool) Option {
	return func(c *Client) error {
		c.debugMode = enabled
		return nil
	}
}

func WithSocks5(addr, user, pass string) Option {
	return func(c *Client) error {
		proxy, err := proxy.SOCKS5("tcp", addr, &proxy.Auth{User: user, Password: pass}, proxy.Direct)
		if err != nil {
			return fmt.Errorf("proxy.SOCKS5: %w", err)
		}
		c.httpClient.Transport = &http.Transport{Dial: proxy.Dial}
		return nil
	}
}
