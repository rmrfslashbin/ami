package main

// file: examples/openai/text-to-speech/main.go

import (
	"context"
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
	resultsDir := filepath.Join(projectRoot, "assets", "results", "openai", "text-to-speech")

	// Ensure the results directory exists
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		fmt.Printf("Failed to create results directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results will be saved in: %s\n\n", resultsDir)

	// Text to be converted to speech
	text := "Hello, this is a test of the OpenAI Text-to-Speech API. Isn't it amazing?"

	// Generate speech
	audio, err := client.TextToSpeech(context.Background(), text, openai.VoiceNova, openai.TTSModel1HD)
	if err != nil {
		fmt.Printf("Error generating speech: %v\n", err)
		os.Exit(1)
	}

	// Save the audio to a file
	timestamp := time.Now().Format("20060102-150405")
	filename = filepath.Join(resultsDir, fmt.Sprintf("tts_output_%s.mp3", timestamp))
	err = os.WriteFile(filename, audio, 0644)
	if err != nil {
		fmt.Printf("Error saving audio file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Text-to-Speech audio saved as %s\n", filename)
	fmt.Printf("Text converted: %s\n", text)
}
