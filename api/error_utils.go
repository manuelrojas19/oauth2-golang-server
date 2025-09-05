package api

// ErrorResponseBody creates a standardized error response body.
func ErrorResponseBody(err error, description ...string) ErrorResponse {
	response := ErrorResponse{
		Error: err.Error(),
	}
	if len(description) > 0 && description[0] != "" {
		response.ErrorDescription = description[0]
	} else if defaultDesc, ok := errorDescriptions[err]; ok {
		response.ErrorDescription = defaultDesc
	}
	return response
}
