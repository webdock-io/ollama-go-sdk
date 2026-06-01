package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func doStream[T any](c *Client, ctx context.Context, endpoint string, body any, fn func(T) error) error {
	if fn == nil {
		return errors.New("ollama: stream callback cannot be nil")
	}

	req, err := c.newRequest(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/x-ndjson")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ollama: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return decodeAPIError(resp)
	}

	decoder := json.NewDecoder(resp.Body)

	for {
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("ollama: decode stream: %w", err)
		}
		if isEmptyRawMessage(raw) {
			continue
		}

		var streamError struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(raw, &streamError); err == nil && streamError.Error != "" {
			return &APIError{
				StatusCode: resp.StatusCode,
				Status:     resp.Status,
				Message:    streamError.Error,
				Body:       raw,
			}
		}

		var chunk T
		if err := json.Unmarshal(raw, &chunk); err != nil {
			return fmt.Errorf("ollama: decode stream chunk: %w", err)
		}
		if !hasStreamResult(chunk) {
			continue
		}

		if err := fn(chunk); err != nil {
			return err
		}
	}
}

func isEmptyRawMessage(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) || bytes.Equal(trimmed, []byte("{}"))
}

func hasStreamResult[T any](chunk T) bool {
	switch v := any(chunk).(type) {
	case GenerateResponse:
		return v.Response != "" || v.Thinking != "" || len(v.LogProbs) > 0
	case ChatResponse:
		return messageHasStreamResult(v.Message) || len(v.LogProbs) > 0
	case StatusResponse:
		return v.Status != "" || v.Digest != "" || v.Total != 0 || v.Completed != 0
	default:
		return true
	}
}

func messageHasStreamResult(message Message) bool {
	return message.Content != "" ||
		message.Thinking != "" ||
		len(message.ToolCalls) > 0 ||
		len(message.Images) > 0
}
