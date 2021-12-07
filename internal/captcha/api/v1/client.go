// Package apiv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.3 DO NOT EDIT.
package apiv1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GenerateCaptcha request
	GenerateCaptcha(ctx context.Context, params *GenerateCaptchaParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// VerifyCaptcha request with any body
	VerifyCaptchaWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	VerifyCaptcha(ctx context.Context, body VerifyCaptchaJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GenerateCaptcha(ctx context.Context, params *GenerateCaptchaParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGenerateCaptchaRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) VerifyCaptchaWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewVerifyCaptchaRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) VerifyCaptcha(ctx context.Context, body VerifyCaptchaJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewVerifyCaptchaRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGenerateCaptchaRequest generates requests for GenerateCaptcha
func NewGenerateCaptchaRequest(server string, params *GenerateCaptchaParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/captcha")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if params.ChallengeMode != nil {

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "challenge_mode", runtime.ParamLocationQuery, *params.ChallengeMode); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewVerifyCaptchaRequest calls the generic VerifyCaptcha builder with application/json body
func NewVerifyCaptchaRequest(server string, body VerifyCaptchaJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewVerifyCaptchaRequestWithBody(server, "application/json", bodyReader)
}

// NewVerifyCaptchaRequestWithBody generates requests for VerifyCaptcha with any type of body
func NewVerifyCaptchaRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/captcha")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GenerateCaptcha request
	GenerateCaptchaWithResponse(ctx context.Context, params *GenerateCaptchaParams, reqEditors ...RequestEditorFn) (*GenerateCaptchaResponse, error)

	// VerifyCaptcha request with any body
	VerifyCaptchaWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*VerifyCaptchaResponse, error)

	VerifyCaptchaWithResponse(ctx context.Context, body VerifyCaptchaJSONRequestBody, reqEditors ...RequestEditorFn) (*VerifyCaptchaResponse, error)
}

type GenerateCaptchaResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Code string `json:"code"`
		Data struct {
			Challenge Challenge `json:"challenge"`
		} `json:"data"`
		Message string `json:"message"`
	}
	JSON400 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r GenerateCaptchaResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GenerateCaptchaResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type VerifyCaptchaResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Code string `json:"code"`
		Data struct {
			Ok bool `json:"ok"`
		} `json:"data"`
		Message string `json:"message"`
	}
	JSON400 *ErrorResponse
	JSON500 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r VerifyCaptchaResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r VerifyCaptchaResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GenerateCaptchaWithResponse request returning *GenerateCaptchaResponse
func (c *ClientWithResponses) GenerateCaptchaWithResponse(ctx context.Context, params *GenerateCaptchaParams, reqEditors ...RequestEditorFn) (*GenerateCaptchaResponse, error) {
	rsp, err := c.GenerateCaptcha(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGenerateCaptchaResponse(rsp)
}

// VerifyCaptchaWithBodyWithResponse request with arbitrary body returning *VerifyCaptchaResponse
func (c *ClientWithResponses) VerifyCaptchaWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*VerifyCaptchaResponse, error) {
	rsp, err := c.VerifyCaptchaWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseVerifyCaptchaResponse(rsp)
}

func (c *ClientWithResponses) VerifyCaptchaWithResponse(ctx context.Context, body VerifyCaptchaJSONRequestBody, reqEditors ...RequestEditorFn) (*VerifyCaptchaResponse, error) {
	rsp, err := c.VerifyCaptcha(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseVerifyCaptchaResponse(rsp)
}

// ParseGenerateCaptchaResponse parses an HTTP response from a GenerateCaptchaWithResponse call
func ParseGenerateCaptchaResponse(rsp *http.Response) (*GenerateCaptchaResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GenerateCaptchaResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Code string `json:"code"`
			Data struct {
				Challenge Challenge `json:"challenge"`
			} `json:"data"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseVerifyCaptchaResponse parses an HTTP response from a VerifyCaptchaWithResponse call
func ParseVerifyCaptchaResponse(rsp *http.Response) (*VerifyCaptchaResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &VerifyCaptchaResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Code string `json:"code"`
			Data struct {
				Ok bool `json:"ok"`
			} `json:"data"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}