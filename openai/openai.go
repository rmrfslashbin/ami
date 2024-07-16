package openai

// file: openai/openai.go

import (
	"context"
	"errors"
	"io"

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

// Image Creation functionality

type ImageSize string

const (
	ImageSize256x256   ImageSize = "256x256"
	ImageSize512x512   ImageSize = "512x512"
	ImageSize1024x1024 ImageSize = "1024x1024"
)

type ImageFormat string

const (
	ImageFormatURL ImageFormat = "url"
	ImageFormatB64 ImageFormat = "b64_json"
)

func (o *OpenAI) CreateImage(ctx context.Context, prompt string, n int, size ImageSize, format ImageFormat) ([]string, error) {
	req := gopenai.ImageRequest{
		Prompt:         prompt,
		N:              n,
		Size:           string(size),
		ResponseFormat: string(format),
	}

	resp, err := o.client.CreateImage(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("no images generated")
	}

	results := make([]string, len(resp.Data))
	for i, data := range resp.Data {
		if format == ImageFormatURL {
			results[i] = data.URL
		} else {
			results[i] = data.B64JSON
		}
	}

	return results, nil
}

// TTSVoice represents the available voices for text-to-speech
type TTSVoice = gopenai.SpeechVoice

// TTSModel represents the available models for text-to-speech
type TTSModel = gopenai.SpeechModel

// Voice constants
const (
	VoiceAlloy   = gopenai.VoiceAlloy
	VoiceEcho    = gopenai.VoiceEcho
	VoiceFable   = gopenai.VoiceFable
	VoiceOnyx    = gopenai.VoiceOnyx
	VoiceNova    = gopenai.VoiceNova
	VoiceShimmer = gopenai.VoiceShimmer
)

// Model constants
const (
	TTSModel1   = gopenai.TTSModel1
	TTSModel1HD = gopenai.TTSModel1HD
)

// TextToSpeech converts text to speech audio
func (o *OpenAI) TextToSpeech(ctx context.Context, input string, voice TTSVoice, model TTSModel) ([]byte, error) {
	req := gopenai.CreateSpeechRequest{
		Model:          model,
		Input:          input,
		Voice:          voice,
		ResponseFormat: gopenai.SpeechResponseFormatMp3,
	}

	resp, err := o.client.CreateSpeech(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	audio, err := io.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	if len(audio) == 0 {
		return nil, errors.New("no audio generated")
	}

	return audio, nil
}

// StreamCompletion streams the chat completion response
func (o *OpenAI) StreamCompletion(ctx context.Context, messages []gopenai.ChatCompletionMessage) (<-chan string, <-chan error) {
	stream, err := o.client.CreateChatCompletionStream(
		ctx,
		gopenai.ChatCompletionRequest{
			Model:    o.model,
			Messages: messages,
			Stream:   true,
		},
	)

	resultChan := make(chan string)
	errChan := make(chan error, 1)

	if err != nil {
		errChan <- err
		close(resultChan)
		close(errChan)
		return resultChan, errChan
	}

	go func() {
		defer stream.Close()
		defer close(resultChan)
		defer close(errChan)

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				errChan <- err
				return
			}
			for _, choice := range response.Choices {
				resultChan <- choice.Delta.Content
			}
		}
	}()

	return resultChan, errChan
}

// StreamTextToSpeech streams the text-to-speech audio
func (o *OpenAI) StreamTextToSpeech(ctx context.Context, input string, voice TTSVoice, model TTSModel) (<-chan []byte, <-chan error) {
	req := gopenai.CreateSpeechRequest{
		Model:          model,
		Input:          input,
		Voice:          voice,
		ResponseFormat: gopenai.SpeechResponseFormatMp3,
	}

	resp, err := o.client.CreateSpeech(ctx, req)

	audioChan := make(chan []byte)
	errChan := make(chan error, 1)

	if err != nil {
		errChan <- err
		close(audioChan)
		close(errChan)
		return audioChan, errChan
	}

	go func() {
		defer resp.Close()
		defer close(audioChan)
		defer close(errChan)

		buffer := make([]byte, 1024)
		for {
			n, err := resp.Read(buffer)
			if err == io.EOF {
				return
			}
			if err != nil {
				errChan <- err
				return
			}
			audioChan <- buffer[:n]
		}
	}()

	return audioChan, errChan
}
