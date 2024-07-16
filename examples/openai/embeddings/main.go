package main

// file: examples/openai/embeddings/main.go

import (
	"context"
	"fmt"
	"os"

	"github.com/rmrfslashbin/ami/openai"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is required")
		os.Exit(1)
	}

	client, err := openai.New(
		openai.WithAPIKey(apiKey),
	)
	if err != nil {
		fmt.Printf("Error creating OpenAI client: %v\n", err)
		os.Exit(1)
	}

	// Single embedding
	input := "Hello, world!"
	embedding, err := client.CreateEmbedding(context.Background(), input)
	if err != nil {
		fmt.Printf("Error creating embedding: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Embedding for '%s':\n", input)
	fmt.Printf("Dimensions: %d\n", len(embedding))
	fmt.Printf("First 5 values: %v\n\n", embedding[:5])

	// Batch embeddings
	inputs := []string{"Hello, world!", "OpenAI is amazing", "Embeddings are useful"}
	embeddings, err := client.CreateEmbeddings(context.Background(), inputs)
	if err != nil {
		fmt.Printf("Error creating batch embeddings: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Batch Embeddings:")
	for i, emb := range embeddings {
		fmt.Printf("Embedding for '%s':\n", inputs[i])
		fmt.Printf("Dimensions: %d\n", len(emb))
		fmt.Printf("First 5 values: %v\n\n", emb[:5])
	}
}
