package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	// DefaultBaseURL is the local Ollama API URL documented by Ollama.
	DefaultBaseURL = "http://localhost:11434/api"

	// CloudBaseURL is the hosted Ollama API URL used with API-key auth.
	CloudBaseURL = "https://ollama.com/api"

	defaultUserAgent = "ollama-go-sdk"
)

// Client is an Ollama API client.
type Client struct {
	Generate   *GenerateService
	Chat       *ChatService
	Embeddings *EmbeddingService
	Models     *ModelService
	Version    *VersionService

	baseURL    *url.URL
	httpClient *http.Client
	userAgent  string
	headers    http.Header
}

// Option configures a Client.
type Option func(*Client) error

// NewClient creates a client for the local Ollama API by default.
func NewClient(opts ...Option) (*Client, error) {
	baseURL, err := parseBaseURL(DefaultBaseURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		userAgent:  defaultUserAgent,
		headers:    make(http.Header),
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	client.initServices()
	return client, nil
}

// New creates a client for the local Ollama API by default.
func New(opts ...Option) (*Client, error) {
	return NewClient(opts...)
}

// MustNew creates a client and panics if configuration is invalid.
func MustNew(opts ...Option) *Client {
	client, err := NewClient(opts...)
	if err != nil {
		panic(err)
	}
	return client
}

// NewCloud creates a client for https://ollama.com/api.
func NewCloud(opts ...Option) (*Client, error) {
	all := []Option{WithBaseURL(CloudBaseURL)}
	all = append(all, opts...)
	return NewClient(all...)
}

func (c *Client) initServices() {
	c.Generate = &GenerateService{client: c}
	c.Chat = &ChatService{client: c}
	c.Embeddings = &EmbeddingService{client: c}
	c.Models = &ModelService{client: c}
	c.Version = &VersionService{client: c}
}

// WithBaseURL sets the API base URL. If only a host is provided, /api is added.
func WithBaseURL(rawURL string) Option {
	return func(c *Client) error {
		baseURL, err := parseBaseURL(rawURL)
		if err != nil {
			return err
		}
		c.baseURL = baseURL
		return nil
	}
}

// WithHTTPClient sets the HTTP client used for requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		if httpClient == nil {
			c.httpClient = http.DefaultClient
			return nil
		}
		c.httpClient = httpClient
		return nil
	}
}

// WithUserAgent sets the User-Agent header sent with requests.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) error {
		c.userAgent = strings.TrimSpace(userAgent)
		return nil
	}
}

// WithHeader sets a header on every request.
func WithHeader(key, value string) Option {
	return func(c *Client) error {
		if strings.TrimSpace(key) == "" {
			return fmt.Errorf("ollama: header key cannot be empty")
		}
		c.headers.Set(key, value)
		return nil
	}
}

// WithHeaders sets multiple headers on every request.
func WithHeaders(headers map[string]string) Option {
	return func(c *Client) error {
		for key, value := range headers {
			if strings.TrimSpace(key) == "" {
				return fmt.Errorf("ollama: header key cannot be empty")
			}
			c.headers.Set(key, value)
		}
		return nil
	}
}

func parseBaseURL(rawURL string) (*url.URL, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil, fmt.Errorf("ollama: base URL cannot be empty")
	}

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("ollama: parse base URL: %w", err)
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, fmt.Errorf("ollama: base URL must include scheme and host")
	}
	if baseURL.Path == "" || baseURL.Path == "/" {
		baseURL.Path = "/api"
	}
	baseURL.Path = strings.TrimRight(baseURL.Path, "/")
	baseURL.RawQuery = ""
	baseURL.Fragment = ""
	return baseURL, nil
}

func (c *Client) endpoint(path string) string {
	next := *c.baseURL
	next.Path = strings.TrimRight(next.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return next.String()
}

func (c *Client) newRequest(ctx context.Context, method, endpoint string, body any) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("ollama: encode request: %w", err)
		}
		reader = &buf
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint(endpoint), reader)
	if err != nil {
		return nil, fmt.Errorf("ollama: create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	for key, values := range c.headers {
		req.Header.Del(key)
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return req, nil
}

func (c *Client) doJSON(ctx context.Context, method, endpoint string, body any, out any) error {
	req, err := c.newRequest(ctx, method, endpoint, body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ollama: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return decodeAPIError(resp)
	}

	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ollama: read response: %w", err)
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("ollama: decode response: %w", err)
	}
	return nil
}
