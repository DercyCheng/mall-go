// Package http provides utilities for HTTP communication between services
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a wrapper around http.Client with additional functionality
type Client struct {
	client      *http.Client
	baseURL     string
	headers     map[string]string
	credentials Credentials
}

// Credentials represents authentication credentials
type Credentials struct {
	Type  string
	Token string
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// NewClient creates a new HTTP client with the given options
func NewClient(baseURL string, options ...ClientOption) *Client {
	c := &Client{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Apply options
	for _, option := range options {
		option(c)
	}

	return c
}

// WithTimeout sets the client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

// WithHeader sets a header for all requests
func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// WithBasicAuth sets basic authentication credentials
func WithBasicAuth(username, password string) ClientOption {
	return func(c *Client) {
		c.credentials = Credentials{
			Type:  "Basic",
			Token: fmt.Sprintf("%s:%s", username, password),
		}
	}
}

// WithBearerToken sets a bearer token for authentication
func WithBearerToken(token string) ClientOption {
	return func(c *Client) {
		c.credentials = Credentials{
			Type:  "Bearer",
			Token: token,
		}
	}
}

// Request represents an HTTP request
type Request struct {
	Method  string
	Path    string
	Query   map[string]string
	Headers map[string]string
	Body    interface{}
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// Do executes an HTTP request
func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	// Build URL with query parameters
	url := c.baseURL + req.Path
	if len(req.Query) > 0 {
		url += "?"
		for k, v := range req.Query {
			url += fmt.Sprintf("%s=%s&", k, v)
		}
		url = url[:len(url)-1] // Remove trailing &
	}

	// Marshal body if present
	var body io.Reader
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewBuffer(jsonBody)
	}

	// Create request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add default headers
	for k, v := range c.headers {
		httpReq.Header.Set(k, v)
	}

	// Add request-specific headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Add authentication if present
	if c.credentials.Type != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("%s %s", c.credentials.Type, c.credentials.Token))
	}

	// Execute request
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer httpResp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create response
	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       respBody,
	}

	return resp, nil
}

// Get sends a GET request
func (c *Client) Get(ctx context.Context, path string, query map[string]string) (*Response, error) {
	return c.Do(ctx, Request{
		Method: "GET",
		Path:   path,
		Query:  query,
	})
}

// Post sends a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.Do(ctx, Request{
		Method: "POST",
		Path:   path,
		Body:   body,
	})
}

// Put sends a PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.Do(ctx, Request{
		Method: "PUT",
		Path:   path,
		Body:   body,
	})
}

// Delete sends a DELETE request
func (c *Client) Delete(ctx context.Context, path string) (*Response, error) {
	return c.Do(ctx, Request{
		Method: "DELETE",
		Path:   path,
	})
}

// ParseJSON parses the response body as JSON
func (r *Response) ParseJSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}
