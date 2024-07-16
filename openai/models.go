package openai

// file: openai/models.go

import (
	"github.com/sashabaranov/go-openai"
)

var ModelsList = map[string]string{
	"gpt-4":             openai.GPT4,
	"gpt-4-32k":         openai.GPT432K,
	"gpt-3.5-turbo":     openai.GPT3Dot5Turbo,
	"gpt-3.5-turbo-16k": openai.GPT3Dot5Turbo16K,
}

// Add more types as needed, based on the go-openai library
