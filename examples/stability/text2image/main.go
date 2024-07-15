// File: examples/stability_generate/main.go

package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/rmrfslashbin/ami/stability"
)

func main() {
	apiKey := os.Getenv("STABILITY_API_KEY")
	if apiKey == "" {
		log.Fatal("STABILITY_API_KEY environment variable is not set")
	}

	client := stability.NewClient(
		stability.WithAPIKey(apiKey),
		stability.WithBaseURL("https://api.stability.ai"),
	)

	// Example for Ultra endpoint
	ultraParams, err := stability.NewGenerateParams(
		stability.WithPrompt("A futuristic cityscape with flying cars"),
		stability.WithAspectRatio("16:9"),
		stability.WithSeed(12345),
		stability.WithOutputFormat("png"),
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
	coreParams, err := stability.NewGenerateParams(
		stability.WithPrompt("A serene lake surrounded by mountains at sunset"),
		stability.WithAspectRatio("16:9"),
		stability.WithSeed(67890),
		stability.WithOutputFormat("png"),
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
	sd3Params, err := stability.NewGenerateParams(
		stability.WithModel("sd3-large"),
		stability.WithPrompt("A mystical forest with glowing fireflies"),
		stability.WithAspectRatio("1:1"),
		stability.WithSeed(54321),
		stability.WithOutputFormat("png"),
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

func saveImage(response *stability.GenerateResponse, filename string) {
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
