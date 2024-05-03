package stability

// Universal constants

// MODULE_NAME is the name of the module
const MODULE_NAME = "stability"

// MAX_PROMPT_LENGTH is the maximum length of a prompt
const MAX_PROMPT_LENGTH = 10000

// MAX_SEED is the maximum seed value
const MAX_SEED = 4294967294

// ASPECT_RATIOS is a list of valid aspect ratios
var ASPECT_RATIOS = []string{"16:9", "1:1", "21:9", "2:3", "3:2", "4:5", "5:4", "9:16", "9:21"}

/*
// STYLE_PRESETS is a list of valid style presets
var STYLE_PRESETS = []string{
	"3d-model", "analog-film", "anime",
	"cinematic", "comic-book", "digital-art",
	"enhance", "fantasy-art", "isometric",
	"line-art", "low-poly", "modeling-compound",
	"neon-punk", "origami", "photographic",
	"pixel-art", "tile-texture",
}
*/

type HttpMethod string

var METHOD_POST HttpMethod = "POST"
var METHOD_GET HttpMethod = "GET"

// STATUS_CODES is a map of status codes to their descriptions
var STATUS_CODES = map[int]string{
	200: "success",
	400: "bad request; see error field for details",
	413: "request entity too large",
	403: "request flagged by moderation system",
	429: "rate limit exceeded; max 150 requests per 10 seconds",
	500: "internal server error",
}

// ENDPOINT_ROOT is the root endpoint for all requests
var ENDPOINT_ROOT = "https://api.stability.ai"

// ENDPOINT_USER_V1 is the endpoint for V1 user requests
var ENDPOINT_USER_V1 = ENDPOINT_ROOT + "/v1/user/account"

// ENDPOINT_USER_BALANCE_V1 is the endpoint for V1 user balance requests
var ENDPOINT_USER_BALANCE_V1 = ENDPOINT_ROOT + "/v1/user/balance"

// V3 constants

// ENDPOINT_V3 is the endpoint for V3 requests
var ENDPOINT_V3 = ENDPOINT_ROOT + "/v2beta/stable-image/generate/sd3"

// MODELS_V3 is a list of valid models for V3 endpoints
var MODELS_V3 = []string{"sd3", "sd3-turbo"}

// OUTPUT_FORMATS_V3 is a list of valid output formats for V3 endpoints
var OUTPUT_FORMATS_V3 = []string{"jpeg", "png"}
