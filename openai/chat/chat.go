package chat

// file: openai/chat/chat.go

import (
	"context"

	"github.com/rmrfslashbin/ami/openai"
	gopenai "github.com/sashabaranov/go-openai"
)

type Chat struct {
	openai   *openai.OpenAI
	messages []gopenai.ChatCompletionMessage
}

func New(o *openai.OpenAI) *Chat {
	return &Chat{
		openai:   o,
		messages: []gopenai.ChatCompletionMessage{},
	}
}

func (c *Chat) AddMessage(role, content string) {
	c.messages = append(c.messages, gopenai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
}

func (c *Chat) Send(ctx context.Context) (gopenai.ChatCompletionResponse, error) {
	req := gopenai.ChatCompletionRequest{
		Model:    c.openai.GetModel(),
		Messages: c.messages,
	}

	return c.openai.GetClient().CreateChatCompletion(ctx, req)
}

// Add more methods as needed
