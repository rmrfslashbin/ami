package stability

// Universal constants

// MODULE_NAME is the name of the module
const MODULE_NAME = "stability"

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
