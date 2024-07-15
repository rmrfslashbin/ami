package main

// file: examples/claude/StreamingMessage/main.go

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/rmrfslashbin/ami/claude"
	"github.com/rmrfslashbin/ami/claude/messages"
)

func main() {
	apikey := os.Getenv("ANTHROPIC_API_KEY")
	if apikey == "" {
		fmt.Println("ANTHROPIC_API_KEY is required")
		os.Exit(1)
	}

	claudeClient, err := claude.New(claude.WithAPIKey(apikey))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client, err := messages.New(
		messages.WithClaude(claudeClient),
		messages.WithHaiku(),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client.SetStreaming(true)

	err = client.AddRoleUser("Tell me a short story about a brave knight.")
	if err != nil {
		log.Fatalf("Failed to add user message: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	results := client.Stream(ctx)

	var wg sync.WaitGroup
	doneChan := make(chan struct{})
	var streamErr error

	// Process events
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(doneChan)

		var fullResponse string
		for event := range results.Response {
			select {
			case <-ctx.Done():
				return
			default:
				switch e := event.(type) {
				case messages.MessageStartEvent:
					fmt.Println("Message started:", e.Message.ID)
				case messages.ContentBlockStartEvent:
					fmt.Printf("Content block started: %d\n", e.Index)
				case messages.ContentBlockDeltaEvent:
					fmt.Print(e.Delta.Text)
					fullResponse += e.Delta.Text
				case messages.MessageDeltaEvent:
					if e.Delta.StopReason != "" {
						fmt.Printf("\nMessage stopped. Reason: %s\n", e.Delta.StopReason)
					}
				case messages.MessageStopEvent:
					fmt.Println("\nMessage completed")
					return
				case messages.StreamingErrorEvent:
					fmt.Printf("Streaming error: %s\n", e.Error.Message)
				default:
					fmt.Printf("Unknown event type: %T\n", e)
				}
			}
		}
		fmt.Println("\nFull response:")
		fmt.Println(fullResponse)
	}()

	// Handle errors
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case err, ok := <-results.Error:
				if !ok {
					return
				}
				if err != nil {
					streamErr = err
					cancel() // Cancel the context to stop streaming
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for streaming to finish or context to be cancelled
	select {
	case <-doneChan:
	case <-ctx.Done():
	}

	// Cancel the context to ensure all goroutines exit
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()

	if streamErr != nil {
		log.Fatalf("Streaming error: %v", streamErr)
	}

	fmt.Println("Program finished successfully")
}
