package main

// file: examples/openai/streaming-tts/main.go

import (
	"context"
	"fmt"
	"io"
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

	resultsDir := getResultsDir()
	fmt.Printf("Results will be saved in: %s\n\n", resultsDir)

	// Text to be converted to speech
	text := "This is a test of streaming text-to-speech. The audio should be clear and in the correct order when saved."

	audioData, err := processAudioStream(client, text)
	if err != nil {
		fmt.Printf("Error processing audio stream: %v\n", err)
		os.Exit(1)
	}

	filename := saveAudioFile(resultsDir, audioData)
	verifyAudioFile(filename)
}

func getResultsDir() string {
	// Get the path of the current file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Failed to get current file path")
		os.Exit(1)
	}

	// Calculate the project root directory (three levels up from the current file)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))

	// Construct the path to the results directory
	resultsDir := filepath.Join(projectRoot, "assets", "results", "openai", "streaming-tts")

	// Ensure the results directory exists
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		fmt.Printf("Failed to create results directory: %v\n", err)
		os.Exit(1)
	}

	return resultsDir
}

func processAudioStream(client *openai.OpenAI, text string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	audioChan, errChan := client.StreamTextToSpeech(ctx, text, openai.VoiceNova, openai.TTSModel1HD)

	var audioData []byte
	fmt.Println("Receiving audio stream...")

	for {
		select {
		case chunk, ok := <-audioChan:
			if !ok {
				return audioData, nil // Stream completed
			}
			audioData = append(audioData, chunk...)
		case err, ok := <-errChan:
			if !ok {
				return audioData, nil // Error channel closed
			}
			if err != nil {
				return nil, fmt.Errorf("stream error: %w", err)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func saveAudioFile(resultsDir string, audioData []byte) string {
	timestamp := time.Now().Format("20060102-150405")
	filename := filepath.Join(resultsDir, fmt.Sprintf("streaming_tts_output_%s.mp3", timestamp))
	err := os.WriteFile(filename, audioData, 0644)
	if err != nil {
		fmt.Printf("Error saving audio file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nStreaming complete. Audio saved as %s\n", filename)
	return filename
}

func verifyAudioFile(filename string) {
	savedFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening saved file: %v\n", err)
		os.Exit(1)
	}
	defer savedFile.Close()

	fileInfo, err := savedFile.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Saved file size: %d bytes\n", fileInfo.Size())

	// Read and verify the first few bytes of the file
	header := make([]byte, 4)
	_, err = io.ReadFull(savedFile, header)
	if err != nil {
		fmt.Printf("Error reading file header: %v\n", err)
		os.Exit(1)
	}

	if string(header) != "ID3\x04" && string(header) != "\xFF\xFB\x90\x64" {
		fmt.Println("Warning: Saved file does not appear to be a valid MP3 file")
	} else {
		fmt.Println("Saved file appears to be a valid MP3 file")
	}
}
