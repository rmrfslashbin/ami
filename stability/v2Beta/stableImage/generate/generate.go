package generate

import (
	"encoding/json"
	"log/slog"
	"slices"

	"github.com/rmrfslashbin/ami/stability"
)

// path: stability/stability.go

//https://platform.stability.ai/docs/api-reference

/*
Method: POST
Headers:
- authorization: Bearer ${API_KEY}
- content-type: multipart/form-data
- accept: image/jpeg, image/png, image/gif, image/webp -- OR -- application/json to receive image as base64 string

Body: form-data
Required:
- prompt ([1 .. 10000] characters)

Optional:
- aspect_ratio (default 1:1; enum)
- mode: text-to-image || image-to-image (default text-to-image)
- negative_prompt (< 10000 chars) not valid for sd3-turbo
- model (default sd3; enum)
- seed (0 .. 4294967294)
- output_format (default png; enum)

Outputs:
- byte array of the generated image
- The resolution of the generated image will be 1 megapixel. The default resolution is 1024x1024.

Credits:
- SD3: Flat rate of 6.5 credits per successful generation of a 1MP image. You will not be charged for failed generations.
- SD3 Turbo: Flat rate of 4 credits per successful generation of a 1MP image. You will not be charged for failed generations.
*/

// MODULE_NAME is the module name
const MODULE_NAME = "stabilityV3"

// Option is a function that takes a pointer to a Config struct and sets a value.
type Option func(config *StabilityV3)

// Configuration structure.
type StabilityV3 struct {
	log            *slog.Logger
	stability      *stability.Stability
	prompt         *string
	aspectRatio    *string
	mode           string
	negativePrompt *string
	model          *string
	seed           *int
	outputFormat   *string
}

// New creates a new StabilityV3 instance.
func New(opts ...func(*StabilityV3)) (*StabilityV3, error) {
	config := &StabilityV3{}

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	if config.stability == nil {
		return nil, &ErrMissingStability{}
	}

	if config.log == nil {
		return nil, &ErrMissingLogger{}
	}

	return config, nil
}

// WithLogger sets the logger for the StabilityV3 instance.
func WithLogger(log *slog.Logger) Option {
	return func(config *StabilityV3) {
		moduleLogger := log.With(
			slog.Group("module_info",
				slog.String("module", MODULE_NAME),
			),
		)
		config.log = moduleLogger
	}
}

// WithStability sets the stability instance for the StabilityV3 instance.
func WithStability(stability *stability.Stability) Option {
	return func(config *StabilityV3) {
		config.stability = stability
	}
}

// WithPrompt sets the prompt for the StabilityV3 instance.
func WithPrompt(prompt string) Option {
	return func(config *StabilityV3) {
		config.prompt = &prompt
	}
}

// WithAspectRatio sets the aspect ratio for the StabilityV3 instance.
func WithAspectRatio(aspectRatio string) Option {
	return func(config *StabilityV3) {
		config.aspectRatio = &aspectRatio
	}
}

// WithNegativePrompt sets the negative prompt for the StabilityV3 instance.
func WithNegativePrompt(negativePrompt string) Option {
	return func(config *StabilityV3) {
		config.negativePrompt = &negativePrompt
	}
}

// WithModel sets the model for the StabilityV3 instance.
func WithModel(model string) Option {
	return func(config *StabilityV3) {
		config.model = &model
	}
}

// WithSeed sets the seed for the StabilityV3 instance.
func WithSeed(seed int) Option {
	return func(config *StabilityV3) {
		config.seed = &seed
	}
}

// WithOutputFormat sets the output format for the StabilityV3 instance.
func WithOutputFormat(outputFormat string) Option {
	return func(config *StabilityV3) {
		config.outputFormat = &outputFormat
	}
}

// SetPrompt sets the prompt for the StabilityV3 instance.
func (c *StabilityV3) SetPrompt(prompt string) {
	c.prompt = &prompt
}

// SetNegativePrompt sets the negative prompt for the StabilityV3 instance.
func (c *StabilityV3) SetNegativePrompt(negativePrompt string) {
	c.negativePrompt = &negativePrompt
}

// SetAspectRatio sets the aspect ratio for the StabilityV3 instance.
func (c *StabilityV3) SetAspectRatio(aspectRatio string) {
	c.aspectRatio = &aspectRatio
}

// SetModel sets the model for the StabilityV3 instance.
func (c *StabilityV3) SetModel(model string) {
	c.model = &model
}

// SetSeed sets the seed for the StabilityV3 instance.
func (c *StabilityV3) SetSeed(seed int) {
	c.seed = &seed
}

// SetStability sets the stability instance for the StabilityV3 instance.
func (c *StabilityV3) SetStability(stability *stability.Stability) {
	c.stability = stability
}

// Generate generates an image using the StabilityV3 instance.
func (c *StabilityV3) Generate() (*Response, error) {
	// validate prompt
	if c.prompt == nil {
		return nil, &ErrMissingPrompt{}
	}

	// check length of prompt
	if len(*c.prompt) < 1 || len(*c.prompt) > MAX_PROMPT_LENGTH {
		return nil, &ErrInvalidPromptLength{}
	}

	c.stability.AddFormPart("prompt", *c.prompt)

	// validate negative prompt
	if c.negativePrompt != nil {
		// check length of negative prompt
		if len(*c.negativePrompt) < 1 || len(*c.negativePrompt) > MAX_PROMPT_LENGTH {
			return nil, &ErrInvalidNegativePromptLength{}
		}
		c.stability.AddFormPart("negative_prompt", *c.negativePrompt)
	}

	// validate aspect ratio
	if c.aspectRatio == nil {
		aspectRatio := DEFAULT_ASPECT_RATIO
		c.aspectRatio = &aspectRatio
	} else {
		// validate aspect ratio
		if !slices.Contains(ASPECT_RATIOS, *c.aspectRatio) {
			return nil, &ErrInvalidAspectRatio{}
		}
	}
	c.stability.AddFormPart("aspect_ratio", *c.aspectRatio)

	// validate model
	if c.model == nil {
		model := DEFAULT_MODEL
		c.model = &model
	} else {
		// validate model
		if !slices.Contains(MODELS, *c.model) {
			return nil, &ErrInvalidModel{}
		}
	}
	c.stability.AddFormPart("model", *c.model)

	// validate seed
	if c.seed == nil {
		seed := DEFAULT_SEED
		c.seed = &seed
	} else {
		// validate seed
		if *c.seed < 0 || *c.seed > MAX_SEED {
			return nil, &ErrInvalidSeed{}
		}
	}
	c.stability.AddFormPart("seed", *c.seed)

	// validate output format
	if c.outputFormat == nil {
		outputFormat := DEFAULT_OUTPUT_FORMAT
		c.outputFormat = &outputFormat
	} else {
		// validate output format
		if !slices.Contains(OUTPUT_FORMATS, *c.outputFormat) {
			return nil, &ErrInvalidOutputFormat{}
		}
	}
	c.stability.AddFormPart("output_format", *c.outputFormat)
	c.stability.AddHeader("accept", "application/json")

	// set default values. Only text-to-image is supported
	c.mode = "text-to-image"
	c.stability.AddFormPart("mode", c.mode)

	// Execute the request
	res, err := c.stability.Do(&ENDPOINT, stability.METHOD_POST)
	if err != nil {
		return nil, err
	}

	// check if response is nil
	if res == nil {
		return nil, &ErrEmptyResponse{}
	}

	// create a response object
	response := &Response{}

	// Unmarshal the response body
	if err = json.Unmarshal(res.Body, response); err != nil {
		return nil, &ErrUnableToParseResponse{Err: err, Response: res.Body}
	}

	return response, nil
}
