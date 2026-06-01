# Ollama Go SDK

A small Go SDK for the Ollama API using an OpenAI-style resource API: typed params in, typed responses out.

## Install

```sh
go get github.com/webdock-io/ollama-go-sdk
```

## Local Ollama

```go
package main

import (
	"context"
	"fmt"
	"log"

	ollama "github.com/webdock-io/ollama-go-sdk"
)

func main() {
	ctx := context.Background()

	client, err := ollama.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Generate.New(ctx, ollama.GenerateNewParams{
		Model:  "gemma3",
		Prompt: "Why is the sky blue?",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Response)
}
```

## Ollama Cloud With Headers

Use `WithHeaders` for Bearer auth or any other custom request headers:

```go
client, err := ollama.NewCloud(
	ollama.WithHeaders(map[string]string{
		"Authorization": "Bearer " + os.Getenv("OLLAMA_API_KEY"),
	}),
)
```

Equivalent explicit configuration:

```go
client, err := ollama.NewClient(
	ollama.WithBaseURL(ollama.CloudBaseURL),
	ollama.WithHeaders(map[string]string{
		"Authorization": "Bearer " + os.Getenv("OLLAMA_API_KEY"),
	}),
)
```

## Chat

```go
res, err := client.Chat.New(ctx, ollama.ChatNewParams{
	Model: "gemma3",
	Messages: []ollama.Message{
		ollama.NewMessage(ollama.RoleSystem, "You are concise."),
		ollama.NewMessage(ollama.RoleUser, "Give me one fact about Saturn."),
	},
})
```

## Streaming

Streaming endpoints use `NewStreaming(ctx, params, fn)`. The SDK sets `stream: true`, parses Ollama's newline-delimited JSON chunks, and skips chunks with no result payload.

```go
err := client.Generate.NewStreaming(ctx, ollama.GenerateNewParams{
	Model:  "gemma3",
	Prompt: "Write a haiku about compilers.",
}, func(chunk ollama.GenerateResponse) error {
	fmt.Print(chunk.Response)
	return nil
})
```

## Other Endpoints

```go
models, err := client.Models.List(ctx)
running, err := client.Models.ListRunning(ctx)
details, err := client.Models.Show(ctx, ollama.ModelShowParams{Model: "gemma3"})
embeddings, err := client.Embeddings.New(ctx, ollama.EmbeddingNewParams{
	Model: "embeddinggemma",
	Input: "hello",
})
version, err := client.Version.Get(ctx)

err = client.Models.Copy(ctx, ollama.ModelCopyParams{
	Source:      "gemma3",
	Destination: "gemma3-backup",
})
err = client.Models.Delete(ctx, ollama.ModelDeleteParams{Model: "gemma3-backup"})

status, err := client.Models.Pull(ctx, ollama.ModelPullParams{Model: "gemma3"})
status, err = client.Models.Push(ctx, ollama.ModelPushParams{Model: "my-username/my-model"})
status, err = client.Models.Create(ctx, ollama.ModelCreateParams{
	Model:  "alpaca",
	From:   "gemma3",
	System: "You are Alpaca, a helpful AI assistant.",
})

_ = models
_ = running
_ = details
_ = embeddings
_ = version
_ = status
```

## Runtime Options

Pass Ollama runtime options in the params struct:

```go
res, err := client.Generate.New(ctx, ollama.GenerateNewParams{
	Model:  "gemma3",
	Prompt: "Explain DNS in one paragraph.",
	Options: ollama.Options{
		ollama.OptionTemperature: 0.2,
		ollama.OptionNumCtx:      4096,
	},
})
```
# ollama-go-sdk
# ollama-go-sdk
# ollama-go-sdk
