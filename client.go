package gprint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	OAuthClient *http.Client
}

type Job struct {
	ID string `json:"id"`
}

func (c *Client) Jobs() ([]Job, error) {
	req, err := http.NewRequest("GET", "https://www.google.com/cloudprint/jobs", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.OAuthClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var obj struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Jobs    []Job  `json:"jobs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&obj); err != nil {
		return nil, err
	} else if !obj.Success {
		return nil, fmt.Errorf("unsuccessful: %s", obj.Message)
	}

	return obj.Jobs, nil
}

func (c *Client) DeleteJob(id string) error {
	v := &url.Values{
		"jobid": []string{id},
	}

	req, err := http.NewRequest("POST", "https://www.google.com/cloudprint/deletejob", strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	resp, err := c.OAuthClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var obj struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&obj); err != nil {
		return err
	} else if !obj.Success {
		return fmt.Errorf("unsuccessful: %s", obj.Message)
	}

	return nil
}
