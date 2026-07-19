package render

import (
	"context"
	"fmt"
	"time"
)

// JobState is the normalized lifecycle state of a generation job.
//
// Providers map their native status strings onto these states; the
// native string is preserved in JobStatus.RawStatus. Unknown provider
// states map to JobStateProcessing (non-terminal, safe for pollers).
type JobState string

const (
	// JobStatePending indicates the job is queued but not yet started.
	JobStatePending JobState = "pending"

	// JobStateProcessing indicates the provider is generating the video.
	JobStateProcessing JobState = "processing"

	// JobStateCompleted indicates the video is ready for download.
	JobStateCompleted JobState = "completed"

	// JobStateFailed indicates generation failed.
	JobStateFailed JobState = "failed"
)

// Terminal reports whether the state is final.
func (s JobState) Terminal() bool {
	return s == JobStateCompleted || s == JobStateFailed
}

// Job identifies a submitted generation job.
type Job struct {
	// ID is the provider job/video identifier.
	ID string

	// Provider is the provider name (e.g., "heygen", "tavus", "bithuman").
	Provider string
}

// JobStatus is a point-in-time snapshot of a job.
type JobStatus struct {
	// ID is the provider job/video identifier.
	ID string

	// State is the normalized lifecycle state.
	State JobState

	// RawStatus is the provider-native status string, preserved for
	// logging and debugging.
	RawStatus string

	// VideoURL is the download URL, set when State is JobStateCompleted.
	// It may be a time-limited signed URL; use Provider.Download for a
	// fresh URL.
	VideoURL string

	// ThumbnailURL is a preview image URL, when the provider reports one.
	ThumbnailURL string

	// Duration is the video duration in seconds, when reported.
	Duration float64

	// ErrorCode is the provider error code, when State is JobStateFailed.
	ErrorCode string

	// ErrorMsg is the provider error message, when State is JobStateFailed.
	ErrorMsg string
}

// DefaultPollInterval is the polling interval used by Wait when the
// caller passes a non-positive interval.
const DefaultPollInterval = 3 * time.Second

// Wait polls p.Status until the job reaches a terminal state, the
// context is cancelled, or a Status call fails. An interval <= 0
// defaults to DefaultPollInterval.
//
// If the job reaches JobStateFailed, Wait returns the final status AND
// an error wrapping ErrJobFailed, so callers can inspect ErrorCode and
// ErrorMsg while still using errors.Is.
func Wait(ctx context.Context, p Provider, jobID string, interval time.Duration) (*JobStatus, error) {
	if interval <= 0 {
		interval = DefaultPollInterval
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		status, err := p.Status(ctx, jobID)
		if err != nil {
			return nil, err
		}
		switch status.State {
		case JobStateCompleted:
			return status, nil
		case JobStateFailed:
			return status, fmt.Errorf("%w: %s: %s (code %q)", ErrJobFailed, jobID, status.ErrorMsg, status.ErrorCode)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}
