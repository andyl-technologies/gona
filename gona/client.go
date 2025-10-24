// Package gona provides a simple golang interface to the NetActuate
// Rest API at https://vapi2.netactuate.com/
package gona

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Version, BaseEndpoint, ContentType constants
const (
	Version      = "0.2.0"
	BaseEndpoint = "https://vapi2.netactuate.com/api/"
	ContentType  = "application/json"
)

// Client is the main object (struct) to which we attach most
// methods/functions.
// It has the following fields:
// (client, userAgent, endPoint, apiKey)
type Client struct {
	client    *http.Client
	userAgent string
	endPoint  *url.URL
	apiKey    string
}

// GetKeyFromEnv is a simple function to grab the value for
// "NA_API_KEY" from the environment
func GetKeyFromEnv() string {
	return os.Getenv("NA_API_KEY")
}

// NewClientCustom is the main entrypoint for instantiating a Client struct.
// It takes your API Key as it's sole argument
// and returns the Client struct ready to talk to the API
func NewClientCustom(apikey string, apiurl string) *Client {
	useragent := "gona/" + Version
	transport := &http.Transport{
		TLSNextProto: make(
			map[string]func(string, *tls.Conn) http.RoundTripper,
		),
	}
	client := http.DefaultClient
	client.Transport = transport
	endpoint, _ := url.Parse(apiurl)

	return &Client{
		userAgent: useragent,
		client:    client,
		endPoint:  endpoint,
		apiKey:    apikey,
	}
}

// NewClient takes an apikey and calls NewClientCustom with the hardcoded
// BaseEndpoint constant API URL
func NewClient(apikey string) *Client {
	return NewClientCustom(apikey, BaseEndpoint)
}

// apiKeyPath is just a short internal function for appending the key to the url
func apiKeyPath(path, apiKey string) string {
	if strings.Contains(path, "?") {
		return path + "&key=" + url.QueryEscape(apiKey)
	}
	return path + "?key=" + url.QueryEscape(apiKey)
}

func (c *Client) debugLog(format string, v ...any) {
	if os.Getenv("NA_API_DEBUG") == "" {
		return
	}
	log.Printf("[DEBUG] "+format, v...)
}

// get internal method on Client struct for providing the HTTP GET call
func (c *Client) get(path string, data interface{}) error {
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return err
	}
	return c.do(req, data)
}

// post internal method on Client struct for providing the HTTP POST call
func (c *Client) post(path string, values []byte, data interface{}) error {
	c.debugLog("POST data for %s: %s", path, string(values))

	req, err := c.newRequest("POST", path, bytes.NewBuffer(values))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req, data)
}

// delete internal method on Client struct for providing the HTTP DELETE call
func (c *Client) delete(path string, values url.Values, data interface{}) error {
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	return c.do(req, data)
}

// Two functions (newRequest, do) below are used by the http method name functions above
// newRequest internal method on Client struct to be wrapped inside the above http method
// named functions for doing the actual work of the get/post/put/patch/delete methods
func (c *Client) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	relPath, err := url.Parse(apiKeyPath(path, c.apiKey))

	if err != nil {
		return nil, err

	}

	url := c.endPoint.ResolveReference(relPath)

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err

	}

	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Accept", ContentType)

	c.debugLog("making a %s request to %s", method, url)
	return req, nil
}

// apiResponse is a message returned by the API that is used both for successful
// responses and for some error responses.
type apiResponse struct {
	Result  string                 `json:"result"`
	Message string                 `json:"message"`
	Data    interface{}            `json:"data"`
	Code    int                    `json:"code"`
	Fields  map[string]interface{} `json:"fields"`
}

// do internal method on Client struct for making the HTTP calls
func (c *Client) do(req *http.Request, data interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.debugLog("got a response: %s", string(body))

	r := &apiResponse{
		Data: data,
	}
	if err := json.Unmarshal(body, r); err != nil {
		return fmt.Errorf("could not unmarshal response %q: %w", string(body), err)
	}

	// Error Handling - This currently ignores invalid mbpkdgid errors to enable the Terraform Provider
	if (resp.StatusCode == 422 || r.Code == 422) && (r.Fields != nil && r.Fields["mbpkgid"] == nil) {
		fieldStr := ""
		for key, value := range r.Fields {
			fieldStr = fieldStr + fmt.Sprintf("%s: %v, ", key, value)
		}
		return fmt.Errorf("got an ERROR response on %s %s: code %d / %d, response: %s / %s", req.Method, req.URL, resp.StatusCode, r.Code, r.Message, fieldStr)
	}

	if (resp.StatusCode != http.StatusOK && resp.StatusCode != 422) || (r.Code != http.StatusOK && r.Code != 422) {
		return fmt.Errorf("got an error response on %s %s: code %d / %d, response: %s / %+v", req.Method, req.URL, resp.StatusCode, r.Code, r.Message, r.Data)
	}

	return nil
}
