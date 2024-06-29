package generate

import "github.com/rmrfslashbin/ami/stability"

// ENDPOINT is the endpoint for V3 requests
var ENDPOINT = stability.ENDPOINT_ROOT + "/v2beta/stable-image/generate/sd3"

// MODELS is a list of valid models for V3 endpoints
var MODELS = []string{"sd3-medium", "sd3-large", "sd3-large-turbo"}

// OUTPUT_FORMATS is a list of valid output formats for V3 endpoints
var OUTPUT_FORMATS = []string{"jpeg", "png", "webp"}

// MAX_PROMPT_LENGTH is the maximum length of a prompt
const MAX_PROMPT_LENGTH = 10000

// MAX_SEED is the maximum seed value
const MAX_SEED = 4294967294

// ASPECT_RATIOS is a list of valid aspect ratios
var ASPECT_RATIOS = []string{"16:9", "1:1", "21:9", "2:3", "3:2", "4:5", "5:4", "9:16", "9:21"}

const DEFAULT_ASPECT_RATIO = "1:1"
const DEFAULT_MODEL = "sd3-large"
const DEFAULT_SEED = 0
const DEFAULT_OUTPUT_FORMAT = "png"
