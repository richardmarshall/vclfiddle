package fiddle

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

const baseURL = "https://fiddle.fastly.dev"

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{},
	}
}

func (c *Client) Get(id string) (*Fiddle, error) {
	createURL := baseURL + "/fiddle/" + id
	req, err := http.NewRequest(http.MethodGet, createURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	j := Resp{}
	err = json.Unmarshal(body, &j)
	if err != nil {
		return nil, err
	}

	return &j.Fiddle, nil
}

func (c *Client) Update(f *Fiddle) (*Fiddle, error) {
	path := "/fiddle"
	if f.ID != "" {
		path += "/" + f.ID
	}
	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	resp, err := c.post(path, b)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	j := Resp{}
	err = json.Unmarshal(body, &j)
	if err != nil {
		return nil, err
	}

	return &j.Fiddle, nil
}

func (c *Client) Clone(f *Fiddle) (*Fiddle, error) {
	// TODO: proper deep copy
	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	o := Fiddle{}
	json.Unmarshal(b, &o)
	o.ID = ""
	return c.Update(&o)
}

func (c *Client) Lock(f *Fiddle) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) Freeze(f *Fiddle) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) PurgeKey(f *Fiddle, key string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) post(path string, data []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", baseURL, path)
	var buf io.Reader
	if data != nil {
		buf = bytes.NewBuffer(data)
	}
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.client.Do(req)
}

type ExecuteOptions struct {
	CacheID int
}

func (o *ExecuteOptions) setDefaults() {
	if o.CacheID == 0 {
		o.CacheID = rand.Intn(100000)
	}
}

func (c *Client) Execute(f *Fiddle, opts ExecuteOptions) (*ExecResults, error) {
	opts.setDefaults()
	path := fmt.Sprintf("/fiddle/%s/execute?cacheID=%d", f.ID, opts.CacheID)
	resp, err := c.post(path, nil)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	j := ExecuteResp{}
	err = json.Unmarshal(body, &j)
	if err != nil {
		return nil, err
	}

	resultsURL := baseURL + "/results/" + j.SessionID + "/stream"
	req, err := http.NewRequest(http.MethodGet, resultsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	resp, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(resp.Body)
	event := ExecResults{}

	// Simplistic handling of Fastly's Server-Sent Events stream
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		_, d, _ := strings.Cut(line, ": ")
		if d != "" {
			if err := json.Unmarshal([]byte(d), &event); err != nil {
				return nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &event, nil
}
