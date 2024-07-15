# Ami Library Overview

Ami is a Go library designed to provide a client for accessing Anthropic's Claude Messages API. Here's a summary of its key components and features:

## Structure
- Main package: `github.com/rmrfslashbin/ami`
- Subpackages:
  - `claude`: Core functionality for interacting with Claude API
  - `claude/messages`: Implements the Messages API

## Key Features
1. Support for Claude 3 models:
   - Claude 3 Opus
   - Claude 3 Sonnet
   - Claude 3 Haiku
   - Claude 3.5 Sonnet

2. API Interactions:
   - Basic message sending
   - Streaming responses
   - Conversation management
   - Tool usage (though not supported in streaming mode)

3. Content Handling:
   - Text messages
   - Image handling (base64 encoded)
   - Tool results

4. Configuration Options:
   - API key management
   - Model selection
   - Max tokens setting
   - Temperature and Top-P controls
   - System prompts

5. Error Handling:
   - Custom error types for various scenarios

6. Conversation Persistence:
   - Save and load conversations to/from files

7. Streaming Support:
   - Event-based streaming of responses
   - Multiple event types (message start, content block, deltas, etc.)

## Usage Examples
The `examples` folder contains sample applications demonstrating:
- Basic message sending
- Streaming message handling

## Dependencies
- `github.com/gabriel-vasile/mimetype`: For MIME type detection
- `github.com/invopop/jsonschema`: For JSON schema handling
- `github.com/rs/xid`: For generating unique IDs
- `github.com/tmaxmax/go-sse`: For server-sent events handling

This library provides a comprehensive set of tools for interacting with Claude's API, supporting both basic and advanced use cases.