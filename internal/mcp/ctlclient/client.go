package ctlclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	do func(method, path string, body io.Reader) ([]byte, error)
}

func New() (*Client, error) {
	dial, err := pipeDialer()
	if err != nil {
		return nil, err
	}
	return &Client{
		do: func(method, path string, body io.Reader) ([]byte, error) {
			req, err := http.NewRequest(method, "http://ctl"+path, body)
			if err != nil {
				return nil, err
			}
			if body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			transport := &http.Transport{DialContext: dial}
			resp, err := (&http.Client{Transport: transport}).Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			out, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode >= 400 {
				return nil, fmt.Errorf("%s", out)
			}
			return out, nil
		},
	}, nil
}

func (c *Client) Status() (map[string]any, error) {
	out, err := c.do(http.MethodGet, "/api/status", nil)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	return result, json.Unmarshal(out, &result)
}

func (c *Client) Operations() (any, error) {
	out, err := c.do(http.MethodGet, "/api/operations", nil)
	if err != nil {
		return nil, err
	}
	var result any
	return result, json.Unmarshal(out, &result)
}

func (c *Client) Execute(name string, args map[string]any) (any, error) {
	payload, _ := json.Marshal(map[string]any{"name": name, "args": args})
	out, err := c.do(http.MethodPost, "/api/execute", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	var result any
	return result, json.Unmarshal(out, &result)
}
