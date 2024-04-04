package state

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"time"
)

type role string

var User role = "user"
var Assistant role = "assistant"

// Conversation is a conversation.
type Conversation struct {
	// Id is the unique object identifier.
	Id string `json:"id"`

	// Model is the model used in the conversation.
	Model string `json:"model"`

	// Created is the time the conversation was created.
	Created time.Time `json:"created"`

	// Updated is the time the conversation was updated.
	Updated time.Time `json:"updated"`

	// Messages is a list of messages in the conversation.
	Messages []*Message `json:"messages"`
}

// Message is a message.
type Message struct {
	// Role is the conversational role of the message.
	Role string `json:"role"`

	// Content is the content of the message.
	Content string `json:"content"`
}

// Option is a configuration option.
type Option func(config *State)

type State struct {
	fpqn         *string
	Conversation *Conversation
}

// New creates a new State configuration.
func New(opts ...func(*State)) (*State, error) {
	config := &State{}

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	if config.fpqn == nil {
		return nil, &ErrMissingFPQN{}
	}

	// Initialize the conversation
	config.Conversation = &Conversation{}

	err := config.Load()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// WithFPQN sets the FPQN configuration.
func WithFPQN(fpqn string) Option {
	return func(config *State) {
		config.fpqn = &fpqn
	}
}

func (s *State) Load() error {
	cleanFqpn, err := filepath.Abs(*s.fpqn)
	if err != nil {
		return err
	}
	s.fpqn = &cleanFqpn

	// Check if the file exists
	_, err = os.Stat(*s.fpqn)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		// Load the file
		file, err := os.Open(*s.fpqn)
		if err != nil {
			return &ErrOpeningFile{Err: err}
		}
		defer file.Close()

		// Decode the GOB data
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&s.Conversation)
		if err != nil {
			return &ErrLoadingGOB{Err: err}
		}
	}

	return nil
}

func (s *State) Save() error {
	// Create the file
	file, err := os.Create(*s.fpqn)
	if err != nil {
		return &ErrOpeningFile{Err: err}
	}
	defer file.Close()

	// Encode the GOB data
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(s.Conversation)
	if err != nil {
		return &ErrSavingGOB{Err: err}
	}

	return nil
}

func (c *Conversation) AddMessage(messageRole role, content string) {
	c.Messages = append(c.Messages, &Message{
		Role:    string(messageRole),
		Content: content,
	})
}

/*

 New -> fqpn
 if fqpn exists, open and load
 if fpqn does not exist, create new struct
*/
