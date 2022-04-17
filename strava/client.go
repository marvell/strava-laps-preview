package strava

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func NewClient(token string, opts ...Option) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		token: token,
	}

	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

type Client struct {
	httpClient *http.Client
	token      string
	debugMode  bool
}

func (c *Client) makeRequest(method, path string, body io.Reader) (*http.Request, error) {
	u, _ := url.Parse(path)
	u.Scheme = "https"
	u.Host = "www.strava.com"
	u.Path = "/api/v3" + u.Path

	return http.NewRequest(method, u.String(), body)
}

func (c *Client) call(req *http.Request) ([]byte, error) {
	if c.debugMode {
		reqDump, _ := httputil.DumpRequestOut(req, true)
		log.Printf("REQ: %s", reqDump)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.httpClient.Get: %w", err)
	}

	if c.debugMode {
		resDump, _ := httputil.DumpResponse(res, true)
		log.Printf("RES: %s", resDump)
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 (%d) response: %s", res.StatusCode, b)
	}

	return b, nil
}

func (c *Client) GetAthleteActivities(limit int) ([]Activity, error) {
	req, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/athlete/activities?per_page=%d", limit), nil)
	if err != nil {
		return nil, fmt.Errorf("c.makeRequest: %w", err)
	}

	b, err := c.call(req)
	if err != nil {
		return nil, fmt.Errorf("c.call: %w", err)
	}

	var v []Activity
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return v, nil
}

func (c *Client) GetActivityLaps(id int) ([]Lap, error) {
	req, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/activities/%d/laps", id), nil)
	if err != nil {
		return nil, fmt.Errorf("c.makeRequest: %w", err)
	}

	b, err := c.call(req)
	if err != nil {
		return nil, fmt.Errorf("c.call: %w", err)
	}

	var v []Lap
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return v, nil
}

func (c *Client) UpdateActivityDescription(id int, desc string) error {
	reqBody := strings.NewReader(url.Values{"description": []string{desc}}.Encode())
	req, err := c.makeRequest(http.MethodPut, fmt.Sprintf("/activities/%d", id), reqBody)
	if err != nil {
		return fmt.Errorf("c.makeRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = c.call(req)
	if err != nil {
		return fmt.Errorf("c.call: %w", err)
	}

	return nil
}
