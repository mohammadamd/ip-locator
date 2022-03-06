package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Middleware struct {
	UserAgent string
	Next      http.RoundTripper
}

//NewClient creates a new instance of invoice Client
func NewClient(app string, version string, baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Transport: &Middleware{
				UserAgent: fmt.Sprintf("%s/%s (%s/%s)", app, version, runtime.GOOS, runtime.GOARCH),
				Next:      http.DefaultTransport,
			},
			Timeout: timeout,
		},
	}
}

func (m Middleware) RoundTrip(req *http.Request) (res *http.Response, e error) {
	req.Header.Set("User-Agent", m.UserAgent)
	return m.Next.RoundTrip(req)
}

// call will do the http request and decode the response into the v
func (c *Client) call(request *http.Request, v interface{}) error {
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	} else if resp.StatusCode > 400 {
		return fmt.Errorf(http.StatusText(resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(v)

	if err != nil {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("body is not Json: %s", string(data))
	}

	return nil
}

// encode will encode the i into a buffer and return it. default encoding format is json
func (c *Client) encode(i interface{}) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(i)
	if err != nil {
		return nil, err
	}
	return buffer, err
}
