package render

import "errors"

// Sentinel errors for render operations.
var (
	// ErrAudioUploadUnsupported indicates that the provider cannot host
	// audio files; callers must supply GenerateRequest.AudioURL.
	ErrAudioUploadUnsupported = errors.New("render: audio upload unsupported")

	// ErrInvalidRequest indicates that the request failed validation
	// before submission (e.g., neither or both of AudioURL/Script set,
	// missing AvatarID).
	ErrInvalidRequest = errors.New("render: invalid request")

	// ErrJobNotFound indicates that the provider does not recognize the
	// job ID.
	ErrJobNotFound = errors.New("render: job not found")

	// ErrJobFailed indicates that the job reached JobStateFailed.
	ErrJobFailed = errors.New("render: job failed")

	// ErrJobNotCompleted indicates that Download was called before the
	// job completed successfully.
	ErrJobNotCompleted = errors.New("render: job not completed")

	// ErrProviderUnavailable indicates that the provider API is
	// unreachable or returned an unexpected error.
	ErrProviderUnavailable = errors.New("render: provider unavailable")

	// ErrProviderAuthFailed indicates that authentication with the
	// provider failed (invalid API key, etc.).
	ErrProviderAuthFailed = errors.New("render: provider authentication failed")

	// ErrInvalidConfig indicates that the provider configuration is
	// invalid or incomplete.
	ErrInvalidConfig = errors.New("render: invalid configuration")
)

// ProviderError wraps an error from a render provider with additional context.
type ProviderError struct {
	Provider string // Provider name (e.g., "heygen", "tavus")
	Op       string // Operation that failed
	Err      error  // Underlying error
}

// Error implements the error interface.
func (e *ProviderError) Error() string {
	if e.Err != nil {
		return "render/" + e.Provider + ": " + e.Op + ": " + e.Err.Error()
	}
	return "render/" + e.Provider + ": " + e.Op
}

// Unwrap returns the underlying error.
func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError creates a new ProviderError.
func NewProviderError(provider, op string, err error) *ProviderError {
	return &ProviderError{
		Provider: provider,
		Op:       op,
		Err:      err,
	}
}
