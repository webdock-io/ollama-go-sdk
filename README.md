# Ollama Go SDK

A lightweight Go SDK for the [Ollama API](https://docs.ollama.com/api/introduction).

The client uses an OpenAI-style shape: services on the client, typed params for requests, and typed responses back.

```go
res, err := client.Generate.New(ctx, ollama.GenerateNewParams{
	Model:  "gemma3",
	Prompt: "Why is the sky blue?",
})
```

## Install

```sh
go get github.com/webdock-io/ollama-go-sdk
```

## Quick Start

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
		Prompt: "Explain DNS in one paragraph.",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Response)
}
```

By default, `NewClient` uses the local Ollama API at `http://localhost:11434/api`.

## Cloud And Headers

Use `NewCloud` for Ollama Cloud, and pass auth or custom headers with `WithHeaders`.

```go
client, err := ollama.NewCloud(
	ollama.WithHeaders(map[string]string{
		"Authorization": "Bearer " + os.Getenv("OLLAMA_API_KEY"),
	}),
)
```

You can also set the base URL explicitly:

```go
client, err := ollama.NewClient(
	ollama.WithBaseURL("https://ollama.com/api"),
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
if err != nil {
	log.Fatal(err)
}

fmt.Println(res.Message.Content)
```

## Streaming

Streaming methods set `stream: true`, parse Ollama's newline-delimited JSON, and skip chunks with no result payload.

```go
err := client.Generate.NewStreaming(ctx, ollama.GenerateNewParams{
	Model:  "gemma3",
	Prompt: "Write a haiku about compilers.",
}, func(chunk ollama.GenerateResponse) error {
	fmt.Print(chunk.Response)
	return nil
})
if err != nil {
	log.Fatal(err)
}
```

Chat streaming works the same way:

```go
err := client.Chat.NewStreaming(ctx, ollama.ChatNewParams{
	Model: "gemma3",
	Messages: []ollama.Message{
		ollama.NewMessage(ollama.RoleUser, "Tell me a short story."),
	},
}, func(chunk ollama.ChatResponse) error {
	fmt.Print(chunk.Message.Content)
	return nil
})
```

## Runtime Options

Runtime options use enum-style keys instead of raw strings.

```go
res, err := client.Generate.New(ctx, ollama.GenerateNewParams{
	Model:  "gemma3",
	Prompt: "Explain DNS in one paragraph.",
	Options: ollama.Options{
		ollama.Temperature: 0.2,
		ollama.NumCtx:      4096,
		ollama.TopP:        0.9,
	},
})
```

## Models

```go
models, err := client.Models.List(ctx)
running, err := client.Models.ListRunning(ctx)

details, err := client.Models.Show(ctx, ollama.ModelShowParams{
	Model: "gemma3",
})

status, err := client.Models.Pull(ctx, ollama.ModelPullParams{
	Model: "gemma3",
})

err = client.Models.Copy(ctx, ollama.ModelCopyParams{
	Source:      "gemma3",
	Destination: "gemma3-backup",
})

err = client.Models.Delete(ctx, ollama.ModelDeleteParams{
	Model: "gemma3-backup",
})

_ = models
_ = running
_ = details
_ = status
```

## Embeddings

```go
res, err := client.Embeddings.New(ctx, ollama.EmbeddingNewParams{
	Model: "embeddinggemma",
	Input: []string{
		"hello",
		"world",
	},
})
if err != nil {
	log.Fatal(err)
}

fmt.Println(len(res.Embeddings))
```

## Version

```go
version, err := client.Version.Get(ctx)
```

## Local Examples

If you add a runnable example inside this repository, put it in a separate folder such as `examples/quickstart`.

Do not place a `package main` example file next to the SDK files, because Go requires one package per directory.
