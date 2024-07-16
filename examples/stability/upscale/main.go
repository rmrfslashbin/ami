package main

// file: examples/stability/upscale/main.go

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rmrfslashbin/ami/stability/upscale"
)

func main() {
	apiKey := os.Getenv("STABILITY_API_KEY")
	if apiKey == "" {
		log.Fatal("STABILITY_API_KEY environment variable is not set")
	}

	skipConservative := os.Getenv("SKIP_CONSERVATIVE") != ""
	skipCreative := os.Getenv("SKIP_CREATIVE") != ""

	if skipConservative && skipCreative {
		log.Fatal("Both SKIP_CONSERVATIVE and SKIP_CREATIVE are set. Nothing to do.")
	}

	client := upscale.New(
		upscale.WithAPIKey(apiKey),
		upscale.WithBaseURL("https://api.stability.ai"),
	)

	// Get the path of the current file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get current file path")
	}

	// Calculate the project root directory (three levels up from the current file)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))

	// Construct the path to the sample image
	sampleImagePath := filepath.Join(projectRoot, "assets", "sample.jpg")

	// Verify that the sample image exists
	if _, err := os.Stat(sampleImagePath); os.IsNotExist(err) {
		log.Fatalf("Sample image not found at %s", sampleImagePath)
	}

	// Ensure the results directory exists
	resultsDir := filepath.Join(projectRoot, "assets", "results", "stability", "upscale")
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		log.Fatalf("Failed to create results directory: %v", err)
	}

	fmt.Printf("Using sample image: %s\n", sampleImagePath)
	fmt.Printf("Results will be saved in: %s\n", resultsDir)

	// Example for Conservative Upscale
	if !skipConservative {
		conservativeParams := &upscale.UpscaleParams{
			Image:  sampleImagePath,
			Width:  1024,
			Height: 1024,
			Prompt: "A high-resolution landscape image",
		}

		conservativeResponse, err := client.ConservativeUpscale(conservativeParams)
		if err != nil {
			log.Printf("Conservative upscale error: %v", err)
		} else {
			saveImage(conservativeResponse.Image, filepath.Join(resultsDir, "conservative_upscaled_image.png"))
		}
	} else {
		fmt.Println("Skipping Conservative Upscale as per SKIP_CONSERVATIVE environment variable")
	}

	// Example for Creative Upscale
	if !skipCreative {
		creativeParams := &upscale.UpscaleParams{
			Image:      sampleImagePath,
			Width:      2048,
			Height:     2048,
			Prompt:     "A beautiful landscape with enhanced details",
			Creativity: 0.35,
		}

		creativeResponse, err := client.CreativeUpscale(creativeParams)
		if err != nil {
			log.Printf("Creative upscale error: %v", err)
		} else {
			fmt.Printf("Creative upscale job started with ID: %s\n", creativeResponse.ID)

			// Poll for results
			maxAttempts := 30 // Maximum number of polling attempts
			for attempt := 0; attempt < maxAttempts; attempt++ {
				result, err := client.GetCreativeUpscaleResult(creativeResponse.ID)
				if err != nil {
					log.Printf("Error fetching creative upscale result: %v", err)
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}

				if result.Image != "" {
					fmt.Println("Creative upscale completed successfully")
					saveImage(result.Image, filepath.Join(resultsDir, "creative_upscaled_image.png"))
					break
				}

				if attempt == maxAttempts-1 {
					log.Println("Max polling attempts reached. The job may still be processing.")
				}

				fmt.Println("Creative upscale still processing...")
				time.Sleep(10 * time.Second) // Wait before polling again
			}
		}
	} else {
		fmt.Println("Skipping Creative Upscale as per SKIP_CREATIVE environment variable")
	}
}

func saveImage(base64Image, filename string) {
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		log.Printf("Error decoding image: %v", err)
		return
	}

	err = os.WriteFile(filename, imageData, 0644)
	if err != nil {
		log.Printf("Error saving image: %v", err)
		return
	}

	fmt.Printf("Image saved as %s\n", filename)
}
