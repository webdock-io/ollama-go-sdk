package ollama

import "time"

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"
)

// Message is a chat message accepted and returned by /api/chat.
type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content,omitempty"`
	Thinking  string     `json:"thinking,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Images    []string   `json:"images,omitempty"`
}

// NewMessage creates a chat message.
func NewMessage(role, content string) Message {
	return Message{Role: role, Content: content}
}

// Tool describes a callable function tool for chat requests.
type Tool struct {
	Type     string       `json:"type,omitempty"`
	Function ToolFunction `json:"function"`
}

// ToolFunction describes a function exposed to the model.
type ToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
}

// ToolCall is a function call requested by the model.
type ToolCall struct {
	Function FunctionCall `json:"function"`
}

// FunctionCall contains function call details returned by the model.
type FunctionCall struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Arguments   map[string]any `json:"arguments,omitempty"`
}

// Metrics contains timing and token usage fields returned by generation endpoints.
type Metrics struct {
	TotalDuration      int64 `json:"total_duration,omitempty"`
	LoadDuration       int64 `json:"load_duration,omitempty"`
	PromptEvalCount    int   `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64 `json:"prompt_eval_duration,omitempty"`
	EvalCount          int   `json:"eval_count,omitempty"`
	EvalDuration       int64 `json:"eval_duration,omitempty"`
}

// LogProb contains log probability information for an output token.
type LogProb struct {
	Token       string       `json:"token"`
	LogProb     float64      `json:"logprob"`
	Bytes       []int        `json:"bytes,omitempty"`
	TopLogProbs []TopLogProb `json:"top_logprobs,omitempty"`
}

// TopLogProb is an alternate token probability.
type TopLogProb struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []int   `json:"bytes,omitempty"`
}

// GenerateResponse is returned by /api/generate.
type GenerateResponse struct {
	Model      string    `json:"model"`
	CreatedAt  time.Time `json:"created_at"`
	Response   string    `json:"response"`
	Thinking   string    `json:"thinking,omitempty"`
	Done       bool      `json:"done"`
	DoneReason string    `json:"done_reason,omitempty"`
	Metrics
	LogProbs []LogProb `json:"logprobs,omitempty"`
}

// ChatResponse is returned by /api/chat.
type ChatResponse struct {
	Model      string    `json:"model"`
	CreatedAt  time.Time `json:"created_at"`
	Message    Message   `json:"message"`
	Done       bool      `json:"done"`
	DoneReason string    `json:"done_reason,omitempty"`
	Metrics
	LogProbs []LogProb `json:"logprobs,omitempty"`
}

// EmbedResponse is returned by /api/embed.
type EmbedResponse struct {
	Model           string      `json:"model"`
	Embeddings      [][]float64 `json:"embeddings"`
	TotalDuration   int64       `json:"total_duration,omitempty"`
	LoadDuration    int64       `json:"load_duration,omitempty"`
	PromptEvalCount int         `json:"prompt_eval_count,omitempty"`
}

// StatusResponse is returned by model management endpoints that report progress.
type StatusResponse struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

// ListModelsResponse is returned by /api/tags.
type ListModelsResponse struct {
	Models []Model `json:"models"`
}

// ListRunningModelsResponse is returned by /api/ps.
type ListRunningModelsResponse struct {
	Models []Model `json:"models"`
}

// Model contains local or running model metadata.
type Model struct {
	Name          string       `json:"name"`
	Model         string       `json:"model"`
	ModifiedAt    time.Time    `json:"modified_at,omitempty"`
	Size          int64        `json:"size,omitempty"`
	Digest        string       `json:"digest,omitempty"`
	Details       ModelDetails `json:"details,omitempty"`
	ExpiresAt     time.Time    `json:"expires_at,omitempty"`
	SizeVRAM      int64        `json:"size_vram,omitempty"`
	ContextLength int          `json:"context_length,omitempty"`
}

// ModelDetails contains high-level model metadata.
type ModelDetails struct {
	ParentModel       string   `json:"parent_model,omitempty"`
	Format            string   `json:"format,omitempty"`
	Family            string   `json:"family,omitempty"`
	Families          []string `json:"families,omitempty"`
	ParameterSize     string   `json:"parameter_size,omitempty"`
	QuantizationLevel string   `json:"quantization_level,omitempty"`
}

// ShowResponse is returned by /api/show.
type ShowResponse struct {
	Parameters   string         `json:"parameters,omitempty"`
	License      any            `json:"license,omitempty"`
	ModifiedAt   time.Time      `json:"modified_at,omitempty"`
	Details      ModelDetails   `json:"details,omitempty"`
	Template     string         `json:"template,omitempty"`
	Capabilities []string       `json:"capabilities,omitempty"`
	ModelInfo    map[string]any `json:"model_info,omitempty"`
}

// VersionResponse is returned by /api/version.
type VersionResponse struct {
	Version string `json:"version"`
}
