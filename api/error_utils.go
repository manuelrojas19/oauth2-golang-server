package api

// ErrorResponseBody creates a standardized error response body.
func ErrorResponseBody(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
	}
}
