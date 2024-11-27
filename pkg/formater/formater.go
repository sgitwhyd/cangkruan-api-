package formater

type (
	response struct {
		Meta meta        `json:"meta"`
		Data interface{} `json:"data"`
	}

	meta struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Status  string `json:"status"`
	}
)

// APIResponse function returns a response object
func APIResponse(message string, code int, status string, data interface{}) response {

	// Construct the meta part of the response
	meta := meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	// Construct and return the full response
	jsonResponse := response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}