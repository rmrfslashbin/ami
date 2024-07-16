package main

// file: examples/openai/chat/main.go

import (
	"context"
	"fmt"
	"os"

	"github.com/rmrfslashbin/ami/openai"
	"github.com/rmrfslashbin/ami/openai/chat"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is required")
		os.Exit(1)
	}

	client, err := openai.New(
		openai.WithAPIKey(apiKey),
		openai.WithModel(openai.ModelsList["gpt-3.5-turbo"]),
	)
	if err != nil {
		fmt.Printf("Error creating OpenAI client: %v\n", err)
		os.Exit(1)
	}

	chatClient := chat.New(client)

	chatClient.AddMessage("user", "Tell me a short story about a brave knight.")

	resp, err := chatClient.Send(context.Background())
	if err != nil {
		fmt.Printf("Error sending chat message: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Response:")
	fmt.Println(resp.Choices[0].Message.Content)
}
