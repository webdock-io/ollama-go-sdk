package ollama

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func testResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestWithHTTPClientUsesCustomClient(t *testing.T) {
	calls := 0
	customClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if req.Method != http.MethodPost {
				t.Fatalf("method = %s, want %s", req.Method, http.MethodPost)
			}
			if got, want := req.URL.String(), "https://ollama.test/api/generate"; got != want {
				t.Fatalf("url = %s, want %s", got, want)
			}
			if got, want := req.Header.Get("Accept"), "application/json"; got != want {
				t.Fatalf("accept = %s, want %s", got, want)
			}
			if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
				t.Fatalf("content-type = %s, want %s", got, want)
			}
			return testResponse(http.StatusOK, `{"model":"gemma3","response":"ok","done":true}`), nil
		}),
	}

	client, err := NewClient(
		WithBaseURL("https://ollama.test/api"),
		WithHTTPClient(customClient),
	)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Generate.New(context.Background(), GenerateNewParams{
		Model:  "gemma3",
		Prompt: "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := res.Response, "ok"; got != want {
		t.Fatalf("response = %q, want %q", got, want)
	}
	if calls != 1 {
		t.Fatalf("custom client calls = %d, want 1", calls)
	}
}

func TestWithHTTPClientUsesCustomClientForStreaming(t *testing.T) {
	calls := 0
	customClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if req.Method != http.MethodPost {
				t.Fatalf("method = %s, want %s", req.Method, http.MethodPost)
			}
			if got, want := req.URL.String(), "https://ollama.test/api/generate"; got != want {
				t.Fatalf("url = %s, want %s", got, want)
			}
			if got, want := req.Header.Get("Accept"), "application/x-ndjson"; got != want {
				t.Fatalf("accept = %s, want %s", got, want)
			}
			return testResponse(http.StatusOK, `{"model":"gemma3","response":"chunk","done":false}
{"model":"gemma3","done":true}
`), nil
		}),
	}
	client, err := NewClient(
		WithBaseURL("https://ollama.test/api"),
		WithHTTPClient(customClient),
	)
	if err != nil {
		t.Fatal(err)
	}

	var chunks []string
	err = client.Generate.NewStreaming(context.Background(), GenerateNewParams{
		Model:  "gemma3",
		Prompt: "hello",
	}, func(chunk GenerateResponse) error {
		chunks = append(chunks, chunk.Response)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("custom client calls = %d, want 1", calls)
	}
	if got, want := strings.Join(chunks, ""), "chunk"; got != want {
		t.Fatalf("chunks = %q, want %q", got, want)
	}
}

func TestWithHTTPClientNilUsesDefaultClient(t *testing.T) {
	client, err := NewClient(WithHTTPClient(nil))
	if err != nil {
		t.Fatal(err)
	}
	if client.httpClient != http.DefaultClient {
		t.Fatal("nil HTTP client should use http.DefaultClient")
	}
}
