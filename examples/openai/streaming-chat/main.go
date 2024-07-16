package main

// file: examples/openai/streaming-chat/main.go

import (
	"context"
	"fmt"
	"os"

	"github.com/rmrfslashbin/ami/openai"
	gopenai "github.com/sashabaranov/go-openai"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is required")
		os.Exit(1)
	}

	client, err := openai.New(
		openai.WithAPIKey(apiKey),
		openai.WithModel(gopenai.GPT3Dot5Turbo),
	)
	if err != nil {
		fmt.Printf("Error creating OpenAI client: %v\n", err)
		os.Exit(1)
	}

	messages := []gopenai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: "Tell me a short story about a brave knight, one sentence at a time.",
		},
	}

	streamChan, errChan := client.StreamCompletion(context.Background(), messages)

	fmt.Println("Streaming response:")
	for {
		select {
		case content, ok := <-streamChan:
			if !ok {
				return
			}
			fmt.Print(content)
		case err, ok := <-errChan:
			if !ok {
				return
			}
			if err != nil {
				fmt.Printf("\n\nError: %v\n", err)
				return
			}
		}
	}
}
