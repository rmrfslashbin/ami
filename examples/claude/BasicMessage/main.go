package main

// file: examples/claude/BasicMessage/main.go

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/ami/claude"
	"github.com/rmrfslashbin/ami/claude/messages"
)

func main() {
	// Get env vars

	// ANTHROPIC_API_KEY
	apikey := os.Getenv("ANTHROPIC_API_KEY")

	if apikey == "" {
		fmt.Println("ANTHROPIC_API_KEY is required")
		os.Exit(1)
	}

	caludeClient, err := claude.New(claude.WithAPIKey(apikey))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client, err := messages.New(
		messages.WithClaude(caludeClient),
		messages.WithHaiku(),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := client.AddRoleUser("Tell me some fun things to do in New Your City"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resp, err := client.Send()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	spew.Dump(resp)

}
