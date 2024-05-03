package user

import (
	"encoding/json"
	"log/slog"

	"github.com/rmrfslashbin/ami/stability"
)

const MODULE_NAME = "stableUser"

type Option func(config *StableUser)

// Configuration structure.
type StableUser struct {
	log       *slog.Logger
	stability *stability.Stability
}

func New(opts ...func(*StableUser)) (*StableUser, error) {
	config := &StableUser{}

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
	return func(config *StableUser) {
		moduleLogger := log.With(
			slog.Group("module_info",
				slog.String("module", MODULE_NAME),
			),
		)
		config.log = moduleLogger
	}
}

func WithStability(stability *stability.Stability) Option {
	return func(config *StableUser) {
		config.stability = stability
	}
}

func (c *StableUser) Me() (*Response, error) {
	c.stability.AddHeader("accept", "application/json")

	res, err := c.stability.Do(&ENDPOINT_USER, stability.METHOD_GET)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, &ErrEmptyResponse{}
	}

	if res.Errors != nil {
		errors := &ResponseUserError{}
		if err := json.Unmarshal(*res.Errors, errors); err != nil {
			return nil, err
		}
		return &Response{Error: errors}, nil
	}

	userResponse := &ResponseUser{}
	if err := json.Unmarshal(res.Body, userResponse); err != nil {
		return nil, err
	}

	return &Response{
		User: userResponse,
	}, nil
}

func (c *StableUser) Balance(input *BalanceInput) (*Response, error) {
	c.stability.AddHeader("accept", "application/json")

	if input != nil {
		if input.Organization != nil {
			c.stability.AddFormPart("Organization", *input.Organization)
		}

		if input.StabilityClientID != nil {
			c.stability.AddFormPart("Stability-Client-ID", *input.StabilityClientID)
		}

		if input.StabilityClientVersion != nil {
			c.stability.AddFormPart("Stability-Client-Version", *input.StabilityClientVersion)
		}
	}

	res, err := c.stability.Do(&ENDPOINT_USER_BALANCE, stability.METHOD_GET)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, &ErrEmptyResponse{}
	}

	if res.Errors != nil {
		errors := &ResponseUserError{}
		if err := json.Unmarshal(*res.Errors, errors); err != nil {
			return nil, err
		}
		return &Response{Error: errors}, nil
	}

	userBalanceResponse := &ResponseUserBalance{}
	if err := json.Unmarshal(res.Body, userBalanceResponse); err != nil {
		return nil, err
	}

	return &Response{
		Credits: userBalanceResponse,
	}, nil
}
