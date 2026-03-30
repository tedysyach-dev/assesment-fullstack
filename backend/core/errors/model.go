package errors

type AppError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
	Internal   error  `json:"-"`
	StackTrace string `json:"-"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
}

type ErrorResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
}
