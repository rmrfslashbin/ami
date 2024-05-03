package generate

import (
	"encoding/json"
	"log/slog"
	"slices"
	"strconv"

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

const MODULE_NAME = "stabilityV3"

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

func WithStability(stability *stability.Stability) Option {
	return func(config *StabilityV3) {
		config.stability = stability
	}
}

func WithPrompt(prompt string) Option {
	return func(config *StabilityV3) {
		config.prompt = &prompt
	}
}

func WithAspectRatio(aspectRatio string) Option {
	return func(config *StabilityV3) {
		config.aspectRatio = &aspectRatio
	}
}

func WithNegativePrompt(negativePrompt string) Option {
	return func(config *StabilityV3) {
		config.negativePrompt = &negativePrompt
	}
}

func WithModel(model string) Option {
	return func(config *StabilityV3) {
		config.model = &model
	}
}

func WithSeed(seed int) Option {
	return func(config *StabilityV3) {
		config.seed = &seed
	}
}

func WithOutputFormat(outputFormat string) Option {
	return func(config *StabilityV3) {
		config.outputFormat = &outputFormat
	}
}

func (c *StabilityV3) SetPrompt(prompt string) {
	c.prompt = &prompt
}

func (c *StabilityV3) SetNegativePrompt(negativePrompt string) {
	c.negativePrompt = &negativePrompt
}

func (c *StabilityV3) SetAspectRatio(aspectRatio string) {
	c.aspectRatio = &aspectRatio
}

func (c *StabilityV3) SetModel(model string) {
	c.model = &model
}

func (c *StabilityV3) SetSeed(seed int) {
	c.seed = &seed
}

func (c *StabilityV3) Generate() (*Response, error) {
	if c.prompt == nil {
		return nil, &ErrMissingPrompt{}
	}
	c.stability.AddFormPart("prompt", *c.prompt)

	if c.aspectRatio == nil {
		aspectRatio := "1:1"
		c.aspectRatio = &aspectRatio
	} else {
		// validate aspect ratio
		if !slices.Contains(ASPECT_RATIOS, *c.aspectRatio) {
			return nil, &ErrInvalidAspectRatio{}
		}
	}
	c.stability.AddFormPart("aspect_ratio", *c.aspectRatio)

	if c.model == nil {
		model := "sd3"
		c.model = &model
	} else {
		// validate model
		if !slices.Contains(MODELS, *c.model) {
			return nil, &ErrInvalidModel{}
		}
	}
	c.stability.AddFormPart("model", *c.model)

	if c.seed == nil {
		seed := 0
		c.seed = &seed
	} else {
		// validate seed
		if *c.seed < 0 || *c.seed > MAX_SEED {
			return nil, &ErrInvalidSeed{}
		}
	}
	c.stability.AddFormPart("seed", *c.seed)

	if c.outputFormat == nil {
		outputFormat := "png"
		c.outputFormat = &outputFormat
	} else {
		// validate output format
		if !slices.Contains(OUTPUT_FORMATS, *c.outputFormat) {
			return nil, &ErrInvalidOutputFormat{}
		}
	}
	c.stability.AddFormPart("output_format", *c.outputFormat)
	c.stability.AddHeader("accept", "image/*") // hard coded to accept image/* for now (default to png)

	// set default values. Only text-to-image is supported
	c.mode = "text-to-image"
	c.stability.AddFormPart("mode", c.mode)

	res, err := c.stability.Do(&ENDPOINT, stability.METHOD_POST)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, &ErrEmptyResponse{}
	}

	response := &Response{}
	if res.Errors != nil {
		responseErrors := &ResponseErrors{}
		if err := json.Unmarshal(*res.Errors, responseErrors); err != nil {
			return nil, err
		}
		response.Errors = responseErrors
		return response, nil
	}

	// check if Content-Type header is present
	if _, ok := res.Headers["Content-Type"]; !ok {
		return nil, &ErrFetchingReturnHeader{Header: "Content-Type"}
	} else {
		if res.Headers["Content-Type"][0] == "application/json" {
			return nil, &ErrJsonResponseNotSupported{}
		}
	}

	// check if Finish-Reason header is present
	if _, ok := res.Headers["Finish-Reason"]; !ok {
		return nil, &ErrFetchingReturnHeader{Header: "Finish-Reason"}
	}

	// check if Seed header is present
	if _, ok := res.Headers["Seed"]; !ok {
		return nil, &ErrFetchingReturnHeader{Header: "Seed"}
	}

	seed, _ := strconv.Atoi(res.Headers["Seed"][0])
	metadata := &ResponseMetadata{
		ContextType:  res.Headers["Content-Type"][0],
		FinishReason: res.Headers["Finish-Reason"][0],
		Seed:         seed,
	}
	response.Metadata = metadata
	response.Image = &res.Body

	return response, nil
}
