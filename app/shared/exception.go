package shared

type ErrorTag string

const (
	MIN_LENGTH_EX      ErrorTag = "MIN_LENGTH"
	MAX_LENGTH_EX      ErrorTag = "MIN_LENGTH"
	ALREADY_CREATED_EX ErrorTag = "MIN_LENGTH"
	INTERNAL_EX        ErrorTag = "INTERNAL"
	DEPENDENCY_EX      ErrorTag = "DEPENDENCY"
	UNKNOWN_EX         ErrorTag = "UNKNOWN"
	APPLICATION_EX     ErrorTag = "APPLICATION_EX"
	UNAUTHORIZED_EX    ErrorTag = "UNAUTHORIZED_EX"
	NOT_FOUND_EX       ErrorTag = "NOT_FOUND_EX"
)

type Exception struct {
	Tag    ErrorTag
	Field  string
	Reason string
}

// Create generic exception.
func DefaultException(tag ErrorTag, reason string) *Exception {
	return &Exception{Tag: tag, Reason: reason}
}

// Create form exception.
func FormException(tag ErrorTag, field string) *Exception {
	return &Exception{Tag: tag, Field: field}
}

// Create internal error exception.
func InternalErrorException() *Exception {
	return &Exception{Tag: INTERNAL_EX}
}
