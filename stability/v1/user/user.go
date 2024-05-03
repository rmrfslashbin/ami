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

func (c *StableUser) Me() (*stability.ResponseUser, error) {
	c.stability.AddHeader("accept", "application/json")

	res, err := c.stability.Do(&stability.ENDPOINT_USER_V1, stability.METHOD_GET)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, &ErrEmptyResponse{}
	}

	if res.Errors != nil {
		errors := &stability.ResponseUserError{}
		if err := json.Unmarshal(*res.Errors, errors); err != nil {
			return nil, err
		}
		return &stability.ResponseUser{Error: errors}, nil
	}

	userResponse := &stability.ResponseUser{}
	if err := json.Unmarshal(res.Body, userResponse); err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (c *StableUser) Balance(input *stability.BalanceInput) (*stability.ResponseUserBalance, error) {
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

	res, err := c.stability.Do(&stability.ENDPOINT_USER_BALANCE_V1, stability.METHOD_GET)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, &ErrEmptyResponse{}
	}

	if res.Errors != nil {
		errors := &stability.ResponseUserError{}
		if err := json.Unmarshal(*res.Errors, errors); err != nil {
			return nil, err
		}
		return &stability.ResponseUserBalance{Error: errors}, nil
	}

	userBalanceResponse := &stability.ResponseUserBalance{}
	if err := json.Unmarshal(res.Body, userBalanceResponse); err != nil {
		return nil, err
	}

	return userBalanceResponse, nil
}
