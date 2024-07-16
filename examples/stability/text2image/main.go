package main

// File: examples/stability/text2image/main.go

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/rmrfslashbin/ami/stability/generate"
)

func main() {
	apiKey := os.Getenv("STABILITY_API_KEY")
	if apiKey == "" {
		log.Fatal("STABILITY_API_KEY environment variable is not set")
	}

	client := generate.New(
		generate.WithAPIKey(apiKey),
		generate.WithBaseURL("https://api.stability.ai"),
	)

	// Example for Ultra endpoint
	ultraParams, err := generate.NewGenerateParams(
		generate.WithPrompt("A futuristic cityscape with flying cars"),
		generate.WithAspectRatio("16:9"),
		generate.WithSeed(12345),
		generate.WithOutputFormat("png"),
	)
	if err != nil {
		log.Printf("Error creating Ultra parameters: %v", err)
	} else {
		ultraResponse, err := client.GenerateUltra(*ultraParams)
		if err != nil {
			log.Printf("Ultra generation error: %v", err)
		} else {
			saveImage(ultraResponse, "ultra_generated_image.png")
		}
	}

	// Example for Core endpoint
	coreParams, err := generate.NewGenerateParams(
		generate.WithPrompt("A serene lake surrounded by mountains at sunset"),
		generate.WithAspectRatio("16:9"),
		generate.WithSeed(67890),
		generate.WithOutputFormat("png"),
	)
	if err != nil {
		log.Printf("Error creating Core parameters: %v", err)
	} else {
		coreResponse, err := client.GenerateCore(*coreParams)
		if err != nil {
			log.Printf("Core generation error: %v", err)
		} else {
			saveImage(coreResponse, "core_generated_image.png")
		}
	}

	// Example for SD3 endpoint
	sd3Params, err := generate.NewGenerateParams(
		generate.WithModel("sd3-large"),
		generate.WithPrompt("A mystical forest with glowing fireflies"),
		generate.WithAspectRatio("1:1"),
		generate.WithSeed(54321),
		generate.WithOutputFormat("png"),
	)
	if err != nil {
		log.Printf("Error creating SD3 parameters: %v", err)
	} else {
		sd3Response, err := client.GenerateSD3(*sd3Params)
		if err != nil {
			log.Printf("SD3 generation error: %v", err)
		} else {
			saveImage(sd3Response, "sd3_generated_image.png")
		}
	}
}

func saveImage(response *generate.GenerateResponse, filename string) {
	fmt.Printf("Generation completed with finish reason: %s\n", response.FinishReason)
	fmt.Printf("Seed used: %d\n", response.Seed)

	imageData, err := base64.StdEncoding.DecodeString(response.Image)
	if err != nil {
		log.Printf("Error decoding image: %v", err)
		return
	}

	err = os.WriteFile(filename, imageData, 0644)
	if err != nil {
		log.Printf("Error saving image: %v", err)
		return
	}

	fmt.Printf("Image saved as %s\n\n", filename)
}
