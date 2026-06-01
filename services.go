package ollama

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

func boolPtr(v bool) *bool {
	return &v
}

// Bool returns a bool pointer for optional params.
func Bool(v bool) *bool {
	return &v
}

// Int returns an int pointer for optional params.
func Int(v int) *int {
	return &v
}

func validModel(model string) bool {
	return strings.TrimSpace(model) != ""
}

func requireNonStreaming(stream *bool) error {
	if stream != nil && *stream {
		return errors.New("ollama: use the streaming method for streaming responses")
	}
	return nil
}

// GenerateService handles /api/generate.
type GenerateService struct {
	client *Client
}

// GenerateNewParams contains params for POST /api/generate.
type GenerateNewParams struct {
	Model       string   `json:"model"`
	Prompt      string   `json:"prompt,omitempty"`
	Suffix      string   `json:"suffix,omitempty"`
	Images      []string `json:"images,omitempty"`
	Format      any      `json:"format,omitempty"`
	System      string   `json:"system,omitempty"`
	Stream      *bool    `json:"stream,omitempty"`
	Think       any      `json:"think,omitempty"`
	Raw         *bool    `json:"raw,omitempty"`
	KeepAlive   any      `json:"keep_alive,omitempty"`
	Options     Options  `json:"options,omitempty"`
	LogProbs    *bool    `json:"logprobs,omitempty"`
	TopLogProbs *int     `json:"top_logprobs,omitempty"`
}

func validateGenerateParams(params GenerateNewParams) error {
	if !validModel(params.Model) {
		return missingField("model")
	}
	return nil
}

// New creates a non-streaming generate request.
func (s *GenerateService) New(ctx context.Context, params GenerateNewParams) (*GenerateResponse, error) {
	if err := validateGenerateParams(params); err != nil {
		return nil, err
	}
	if err := requireNonStreaming(params.Stream); err != nil {
		return nil, err
	}

	params.Stream = boolPtr(false)
	var out GenerateResponse
	if err := s.client.doJSON(ctx, http.MethodPost, "generate", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// NewStreaming creates a streaming generate request and calls fn for each chunk.
func (s *GenerateService) NewStreaming(ctx context.Context, params GenerateNewParams, fn func(GenerateResponse) error) error {
	if err := validateGenerateParams(params); err != nil {
		return err
	}
	params.Stream = boolPtr(true)
	return doStream(s.client, ctx, "generate", params, fn)
}

// ChatService handles /api/chat.
type ChatService struct {
	client *Client
}

// ChatNewParams contains params for POST /api/chat.
type ChatNewParams struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Tools       []Tool    `json:"tools,omitempty"`
	Format      any       `json:"format,omitempty"`
	Options     Options   `json:"options,omitempty"`
	Stream      *bool     `json:"stream,omitempty"`
	Think       any       `json:"think,omitempty"`
	KeepAlive   any       `json:"keep_alive,omitempty"`
	LogProbs    *bool     `json:"logprobs,omitempty"`
	TopLogProbs *int      `json:"top_logprobs,omitempty"`
}

func validateChatParams(params ChatNewParams) error {
	if !validModel(params.Model) {
		return missingField("model")
	}
	if len(params.Messages) == 0 {
		return missingField("messages")
	}
	return nil
}

// New creates a non-streaming chat request.
func (s *ChatService) New(ctx context.Context, params ChatNewParams) (*ChatResponse, error) {
	if err := validateChatParams(params); err != nil {
		return nil, err
	}
	if err := requireNonStreaming(params.Stream); err != nil {
		return nil, err
	}

	params.Stream = boolPtr(false)
	var out ChatResponse
	if err := s.client.doJSON(ctx, http.MethodPost, "chat", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// NewStreaming creates a streaming chat request and calls fn for each chunk.
func (s *ChatService) NewStreaming(ctx context.Context, params ChatNewParams, fn func(ChatResponse) error) error {
	if err := validateChatParams(params); err != nil {
		return err
	}
	params.Stream = boolPtr(true)
	return doStream[ChatResponse](s.client, ctx, "chat", params, fn)
}

// EmbeddingService handles /api/embed.
type EmbeddingService struct {
	client *Client
}

// EmbeddingNewParams contains params for POST /api/embed.
type EmbeddingNewParams struct {
	Model      string  `json:"model"`
	Input      any     `json:"input"`
	Truncate   *bool   `json:"truncate,omitempty"`
	Dimensions *int    `json:"dimensions,omitempty"`
	Options    Options `json:"options,omitempty"`
}

func validateEmbeddingParams(params EmbeddingNewParams) error {
	if !validModel(params.Model) {
		return missingField("model")
	}
	if params.Input == nil {
		return missingField("input")
	}
	return nil
}

// New creates an embedding request.
func (s *EmbeddingService) New(ctx context.Context, params EmbeddingNewParams) (*EmbedResponse, error) {
	if err := validateEmbeddingParams(params); err != nil {
		return nil, err
	}

	var out EmbedResponse
	if err := s.client.doJSON(ctx, http.MethodPost, "embed", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ModelService handles Ollama model endpoints.
type ModelService struct {
	client *Client
}

// ModelShowParams contains params for POST /api/show.
type ModelShowParams struct {
	Model   string `json:"model"`
	Verbose *bool  `json:"verbose,omitempty"`
}

// ModelCreateParams contains params for POST /api/create.
type ModelCreateParams struct {
	Model      string         `json:"model"`
	From       string         `json:"from,omitempty"`
	Template   string         `json:"template,omitempty"`
	License    any            `json:"license,omitempty"`
	System     string         `json:"system,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty"`
	Messages   []Message      `json:"messages,omitempty"`
	Quantize   string         `json:"quantize,omitempty"`
	Stream     *bool          `json:"stream,omitempty"`
}

// ModelCopyParams contains params for POST /api/copy.
type ModelCopyParams struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// ModelPullParams contains params for POST /api/pull.
type ModelPullParams struct {
	Model    string `json:"model"`
	Insecure *bool  `json:"insecure,omitempty"`
	Stream   *bool  `json:"stream,omitempty"`
}

// ModelPushParams contains params for POST /api/push.
type ModelPushParams struct {
	Model    string `json:"model"`
	Insecure *bool  `json:"insecure,omitempty"`
	Stream   *bool  `json:"stream,omitempty"`
}

// ModelDeleteParams contains params for DELETE /api/delete.
type ModelDeleteParams struct {
	Model string `json:"model"`
}

// List calls GET /api/tags.
func (s *ModelService) List(ctx context.Context) (*ListModelsResponse, error) {
	var out ListModelsResponse
	if err := s.client.doJSON(ctx, http.MethodGet, "tags", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListRunning calls GET /api/ps.
func (s *ModelService) ListRunning(ctx context.Context) (*ListRunningModelsResponse, error) {
	var out ListRunningModelsResponse
	if err := s.client.doJSON(ctx, http.MethodGet, "ps", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Show calls POST /api/show.
func (s *ModelService) Show(ctx context.Context, params ModelShowParams) (*ShowResponse, error) {
	if !validModel(params.Model) {
		return nil, missingField("model")
	}

	var out ShowResponse
	if err := s.client.doJSON(ctx, http.MethodPost, "show", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create calls POST /api/create.
func (s *ModelService) Create(ctx context.Context, params ModelCreateParams) (*StatusResponse, error) {
	if !validModel(params.Model) {
		return nil, missingField("model")
	}
	if err := requireNonStreaming(params.Stream); err != nil {
		return nil, err
	}

	params.Stream = boolPtr(false)
	var out StatusResponse
	if err := s.client.doJSON(ctx, http.MethodPost, "create", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateStreaming calls POST /api/create with stream enabled.
func (s *ModelService) CreateStreaming(ctx context.Context, params ModelCreateParams, fn func(StatusResponse) error) error {
	if !validModel(params.Model) {
		return missingField("model")
	}
	params.Stream = boolPtr(true)
	return doStream[StatusResponse](s.client, ctx, "create", params, fn)
}

// Copy calls POST /api/copy.
func (s *ModelService) Copy(ctx context.Context, params ModelCopyParams) error {
	if strings.TrimSpace(params.Source) == "" {
		return missingField("source")
	}
	if strings.TrimSpace(params.Destination) == "" {
		return missingField("destination")
	}
	return s.client.doJSON(ctx, http.MethodPost, "copy", params, nil)
}

// Pull calls POST /api/pull.
func (s *ModelService) Pull(ctx context.Context, params ModelPullParams) (*StatusResponse, error) {
	if !validModel(params.Model) {
		return nil, missingField("model")
	}
	if err := requireNonStreaming(params.Stream); err != nil {
		return nil, err
	}

	params.Stream = boolPtr(false)
	var out StatusResponse
	if err := s.client.doJSON(ctx, http.MethodPost, "pull", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PullStreaming calls POST /api/pull with stream enabled.
func (s *ModelService) PullStreaming(ctx context.Context, params ModelPullParams, fn func(StatusResponse) error) error {
	if !validModel(params.Model) {
		return missingField("model")
	}
	params.Stream = boolPtr(true)
	return doStream[StatusResponse](s.client, ctx, "pull", params, fn)
}

// Push calls POST /api/push.
func (s *ModelService) Push(ctx context.Context, params ModelPushParams) (*StatusResponse, error) {
	if !validModel(params.Model) {
		return nil, missingField("model")
	}
	if err := requireNonStreaming(params.Stream); err != nil {
		return nil, err
	}

	params.Stream = boolPtr(false)
	var out StatusResponse
	if err := s.client.doJSON(ctx, http.MethodPost, "push", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PushStreaming calls POST /api/push with stream enabled.
func (s *ModelService) PushStreaming(ctx context.Context, params ModelPushParams, fn func(StatusResponse) error) error {
	if !validModel(params.Model) {
		return missingField("model")
	}
	params.Stream = boolPtr(true)
	return doStream[StatusResponse](s.client, ctx, "push", params, fn)
}

// Delete calls DELETE /api/delete.
func (s *ModelService) Delete(ctx context.Context, params ModelDeleteParams) error {
	if !validModel(params.Model) {
		return missingField("model")
	}
	return s.client.doJSON(ctx, http.MethodDelete, "delete", params, nil)
}

// VersionService handles /api/version.
type VersionService struct {
	client *Client
}

// Get calls GET /api/version.
func (s *VersionService) Get(ctx context.Context) (*VersionResponse, error) {
	var out VersionResponse
	if err := s.client.doJSON(ctx, http.MethodGet, "version", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
