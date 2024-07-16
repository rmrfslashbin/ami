package openai

// file: openai/openai.go

import (
	"context"

	gopenai "github.com/sashabaranov/go-openai"
)

const MODULE_NAME = "openai"

type OpenAI struct {
	client *gopenai.Client
	model  string
}

type Option func(*OpenAI)

func New(opts ...Option) (*OpenAI, error) {
	o := &OpenAI{}

	for _, opt := range opts {
		opt(o)
	}

	if o.client == nil {
		return nil, &ErrMissingClient{}
	}

	if o.model == "" {
		o.model = gopenai.GPT3Dot5Turbo
	}

	return o, nil
}

func WithAPIKey(apiKey string) Option {
	return func(o *OpenAI) {
		o.client = gopenai.NewClient(apiKey)
	}
}

func WithModel(model string) Option {
	return func(o *OpenAI) {
		o.model = model
	}
}

func (o *OpenAI) GetClient() *gopenai.Client {
	return o.client
}

func (o *OpenAI) GetModel() string {
	return o.model
}

// Embeddings functionality

func (o *OpenAI) CreateEmbedding(ctx context.Context, input string) ([]float32, error) {
	req := gopenai.EmbeddingRequest{
		Input: []string{input},
		Model: gopenai.AdaEmbeddingV2,
	}

	resp, err := o.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) > 0 {
		return resp.Data[0].Embedding, nil
	}

	return []float32{}, nil
}

func (o *OpenAI) CreateEmbeddings(ctx context.Context, inputs []string) ([][]float32, error) {
	req := gopenai.EmbeddingRequest{
		Input: inputs,
		Model: gopenai.AdaEmbeddingV2,
	}

	resp, err := o.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	embeddings := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}
