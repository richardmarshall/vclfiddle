package fiddle

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"math/rand"
	"net/http"
)

const defaultBaseURL = "https://fiddle.fastly.dev"
const userAgent = "vclfiddle"

type Client struct {
	client  *http.Client
	baseURL string

	UserAgent string
}

func NewClient() *Client {
	return &Client{
		client:    &http.Client{},
		baseURL:   defaultBaseURL,
		UserAgent: userAgent,
	}
}

type RequestOptions struct {
	authorization string
}

type RequestOptionFunc func(*http.Request)

func WithPassword(pw string) RequestOptionFunc {
	return func(r *http.Request) {
		r.Header.Set("Authorization", pw)
	}
}

func (c *Client) NewRequest(method, path string, body any, opts ...RequestOptionFunc) (*http.Request, error) {
	url := c.baseURL + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	for _, fn := range opts {
		fn(req)
	}
	req.Header.Set("Accept", "application/json")
	switch method {
	case http.MethodPost, http.MethodPut:
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		if b != nil {
			buf := bytes.NewBuffer(b)
			req.Body = io.NopCloser(buf)
			req.Header.Set("Content-Type", "application/json")
		}
	}
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v any) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err = checkResponse(resp); err != nil {
		return resp, err
	}
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

var ErrNotFound = errors.New("not found")

type FiddleError struct {
	Response *http.Response
	Message  string
}

func (e *FiddleError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("fiddle error: %s", e.Message)
	}
	return fmt.Sprintf("fiddle error")
}

func checkResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201:
		return nil
	case 404:
		return ErrNotFound
	}

	respErr := &FiddleError{Response: resp}
	body, err := io.ReadAll(resp.Body)
	if err != nil && len(body) > 0 {
		respErr.Message = string(body)
	}
	return respErr
}

func (c *Client) Get(id string) (*Fiddle, bool, Lints, error) {
	req, err := c.NewRequest(http.MethodGet, "/fiddle/"+id, nil)
	if err != nil {
		return nil, false, nil, err
	}
	var fiddleResp Resp
	_, err = c.Do(req, &fiddleResp)
	if err != nil {
		return nil, false, nil, err
	}

	return &fiddleResp.Fiddle, fiddleResp.Valid, fiddleResp.LintStatus, nil
}

func (c *Client) Create(f *Fiddle, opts ...RequestOptionFunc) (*Fiddle, error) {
	req, err := c.NewRequest(http.MethodPost, "/fiddle", f, opts...)
	if err != nil {
		return nil, err
	}

	var fiddleResp Resp
	_, err = c.Do(req, &fiddleResp)
	if err != nil {
		return nil, err
	}
	return &fiddleResp.Fiddle, nil
}

func (c *Client) Update(f *Fiddle, opts ...RequestOptionFunc) (*Fiddle, error) {
	if f.ID == "" {
		return nil, errors.New("Fiddle ID is required")
	}
	path := fmt.Sprintf("/fiddle/%s", f.ID)
	req, err := c.NewRequest(http.MethodPut, path, f, opts...)
	if err != nil {
		return nil, err
	}

	var fiddleResp Resp
	_, err = c.Do(req, &fiddleResp)
	if err != nil {
		return nil, err
	}

	return &fiddleResp.Fiddle, nil
}

func (c *Client) Clone(f *Fiddle) (*Fiddle, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	o := Fiddle{}
	json.Unmarshal(b, &o)
	o.ID = ""
	return c.Create(&o)
}

func (c *Client) Lock(f *Fiddle, password string) (string, error) {
	path := fmt.Sprintf("/fiddle/%s/lock", f.ID)
	req, err := c.NewRequest(http.MethodPost, path, nil, WithPassword(password))
	if err != nil {
		return "", err
	}
	var lockResp string
	_, err = c.Do(req, &lockResp)
	return lockResp, err
}

func (c *Client) PurgeKey(f *Fiddle, key string) (bool, error) {
	path := fmt.Sprintf("/fiddle/%s/purgeKey", f.ID)
	req, err := c.NewRequest(http.MethodPost, path, Purge{Key: key})
	if err != nil {
		return false, err
	}
	var purgeResp bool
	_, err = c.Do(req, &purgeResp)
	return purgeResp, err
}

type ExecuteOptions struct {
	CacheID int
}

func (o *ExecuteOptions) setDefaults() {
	if o.CacheID == 0 {
		o.CacheID = rand.Intn(100000)
	}
}

// Execute Fiddle and return the final StreamEvent.
func (c *Client) Execute(f *Fiddle, opts ExecuteOptions) (*StreamEvent, error) {
	iter, err := c.ExecuteIter(f, opts)
	if err != nil {
		return nil, err
	}

	var event *StreamEvent
	for event, err = range iter {
		if err != nil {
			return nil, err
		}
	}

	return event, nil
}

var (
	eventPrefix = []byte("event:")
	dataPrefix  = []byte("data:")
)

// Execute Fiddle and return and iterator which will yield all the events produced during the execution.
func (c *Client) ExecuteIter(f *Fiddle, opts ExecuteOptions) (iter.Seq2[*StreamEvent, error], error) {
	opts.setDefaults()
	path := fmt.Sprintf("/fiddle/%s/execute?cacheID=%d", f.ID, opts.CacheID)
	req, err := c.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}
	var execResp ExecuteResponse
	resp, err := c.Do(req, &execResp)
	if err != nil {
		return nil, err
	}

	resultsURL := fmt.Sprintf("/results/%s/stream", execResp.SessionID)
	req, err = c.NewRequest(http.MethodGet, resultsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	resp, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Simplistic handling of Fastly's Server-Sent Events stream
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(split)

	return func(yield func(*StreamEvent, error) bool) {
		for scanner.Scan() {
			data := scanner.Bytes()

			var eventType string
			var eventData []byte
			for _, l := range bytes.Split(data, []byte("\n")) {
				if bytes.HasPrefix(l, eventPrefix) {
					eventType = string(bytes.TrimSpace(bytes.TrimPrefix(l, eventPrefix)))
				}
				if bytes.HasPrefix(l, dataPrefix) {
					eventData = bytes.TrimPrefix(l, dataPrefix)
				}
			}

			switch eventType {
			case "updateResult", "waitForSync":
				var execEvent StreamEvent
				err := json.Unmarshal(eventData, &execEvent)
				execEvent.Type = eventType
				if !yield(&execEvent, err) {
					return
				}
			case "timeout":
				return
			default:
				continue
			}
		}
		if err := scanner.Err(); err != nil {
			yield(nil, err)
		}
	}, nil
}

// Split stream data into chunks delimited by double newlines.
func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
		return i + 2, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
