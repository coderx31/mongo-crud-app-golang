package customerror

// for custom error
type CustomError struct {
	Error      error `json:"error"`
	StatusCode int16 `json:"status_code"`
}
