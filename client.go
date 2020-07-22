// Package gerrit provides functionality of for talking to Gerrit APIs.
package gerrit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// NewClient creates a new gerrit client with the given root (no trailing slash)
// and user/password to use for basic HTTP auth.
func NewClient(rootPath, user, password string) *Client {
	return &Client{
		Client: http.DefaultClient,
		root:   rootPath,
		user:   user,
		pass:   password,
	}
}

// Client provides methods for making requests to the Gerrit REST API.
type Client struct {
	*http.Client
	root       string
	user, pass string
}

type emptyReader struct{}

func (emptyReader) Read(p []byte) (n int, err error) { return 0, io.EOF }

// CallError is returned from Call if the response failed.
type CallError struct {
	Err      error
	Response []byte
}

func (c *CallError) Error() string { return c.Err.Error() }

// Call a url using the given method and body.
func (c *Client) Call(ctx context.Context, method, url string, body, resp interface{}) error {
	if strings.HasPrefix(url, "/a/") {
		return fmt.Errorf("invalid url: must not begin with /a/: %q", url)
	}
	url = strings.TrimPrefix(url, "/") // remove leading /

	var r io.Reader = emptyReader{}
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.root+"/a/"+url, r)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	}
	req.SetBasicAuth(c.user, c.pass)

	response, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, _ := ioutil.ReadAll(response.Body)
		return &CallError{
			Err:      fmt.Errorf("response status != 200 (%v)", response.Status),
			Response: responseBody,
		}
	}

	// Remove the prefix at the beginning of each response.
	var prefix [5]byte
	if _, err = io.ReadFull(response.Body, prefix[:]); err != nil || !bytes.Equal(prefix[:], invalidPrefix) {
		return fmt.Errorf("expected prefix %q, got %q", invalidPrefix, prefix)
	}
	return json.NewDecoder(response.Body).Decode(resp)
}

// invalidPrefix is the junk that gerrit spews out first.
var invalidPrefix = []byte(")]}'\n")
