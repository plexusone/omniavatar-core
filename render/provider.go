// Package render provides provider-agnostic interfaces for asynchronous
// (batch) avatar video generation: submit narration audio or a script,
// poll for completion, and download a talking-head video.
//
// This is the offline counterpart to package live, which handles
// real-time streaming avatar sessions.
//
// The typical flow is:
//
//	job, err := provider.Generate(ctx, render.GenerateRequest{
//	    AvatarID: avatarID,
//	    AudioURL: narrationURL,
//	})
//	status, err := render.Wait(ctx, provider, job.ID, 5*time.Second)
//	err = provider.Download(ctx, job.ID, outFile)
package render

import (
	"context"
	"io"
)

// Provider generates avatar videos asynchronously.
//
// Implementations of this interface handle provider-specific API calls
// to submit generation jobs, poll their status, and download results.
// Each provider (HeyGen, Tavus, bitHuman, etc.) has its own Provider
// implementation.
type Provider interface {
	// Name returns the provider name (e.g., "heygen", "tavus", "bithuman").
	Name() string

	// Generate submits a video generation job. It returns as soon as the
	// provider accepts the job; use Status or Wait to track completion.
	Generate(ctx context.Context, req GenerateRequest) (*Job, error)

	// Status returns the current status of a job.
	Status(ctx context.Context, jobID string) (*JobStatus, error)

	// Download streams the completed video to dst.
	// Returns an error wrapping ErrJobNotCompleted if the job has not
	// completed successfully.
	Download(ctx context.Context, jobID string, dst io.Writer) error
}

// AudioUploader is an optional capability for providers that can host
// local audio files. Callers should feature-detect:
//
//	if up, ok := provider.(render.AudioUploader); ok {
//	    url, err = up.UploadAudio(ctx, "narration.mp3", f)
//	}
//
// Providers without hosting support do not implement this interface;
// callers must supply a publicly fetchable GenerateRequest.AudioURL.
type AudioUploader interface {
	// UploadAudio uploads audio content and returns a URL usable as
	// GenerateRequest.AudioURL with the same provider.
	UploadAudio(ctx context.Context, filename string, r io.Reader) (string, error)
}
