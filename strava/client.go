package strava

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
}

func (c *Client) makeUrl(path string) url.URL {
	u, _ := url.Parse(path)
	u.Scheme = "https"
	u.Host = "www.strava.com"
	u.Path = "/api/v3" + u.Path

	return *u
}

func (c *Client) call(u url.URL) ([]byte, error) {
	log.Printf("URL: %s", u.String())

	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.httpClient.Get: %w", err)
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
	b, err := c.call(c.makeUrl(fmt.Sprintf("/athlete/activities?per_page=%d", limit)))
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
	b, err := c.call(c.makeUrl(fmt.Sprintf("/activities/%d/laps", id)))
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
