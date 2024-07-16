package main

// file: examples/openai/image-creation/main.go

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

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

	// Get the path of the current file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Failed to get current file path")
		os.Exit(1)
	}

	// Calculate the project root directory (three levels up from the current file)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))

	// Construct the path to the results directory
	resultsDir := filepath.Join(projectRoot, "assets", "results", "openai", "image-creation")

	// Ensure the results directory exists
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		fmt.Printf("Failed to create results directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results will be saved in: %s\n\n", resultsDir)

	// Create an image with URL response
	prompt := "A futuristic city with flying cars and neon lights"
	urlImages, err := client.CreateImage(context.Background(), prompt, 1, openai.ImageSize1024x1024, openai.ImageFormatURL)
	if err != nil {
		fmt.Printf("Error creating image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Image created for prompt: '%s'\n", prompt)
	fmt.Printf("Image URL: %s\n\n", urlImages[0])

	// Create an image with base64 response
	prompt = "A serene mountain landscape with a calm lake"
	b64Images, err := client.CreateImage(context.Background(), prompt, 1, openai.ImageSize512x512, openai.ImageFormatB64)
	if err != nil {
		fmt.Printf("Error creating image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Image created for prompt: '%s'\n", prompt)
	fmt.Println("Image data (first 100 characters of base64):")
	fmt.Printf("%s...\n\n", b64Images[0][:100])

	// Save the base64 image to a file
	imageData, err := base64.StdEncoding.DecodeString(b64Images[0])
	if err != nil {
		fmt.Printf("Error decoding base64 image: %v\n", err)
		os.Exit(1)
	}

	timestamp := time.Now().Format("20060102-150405")
	filename = filepath.Join(resultsDir, fmt.Sprintf("generated_image_%s.png", timestamp))
	err = os.WriteFile(filename, imageData, 0644)
	if err != nil {
		fmt.Printf("Error saving image file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Image saved as %s\n", filename)
}
