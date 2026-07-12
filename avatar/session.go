package avatar

import (
	"context"
	"time"
)

// Session manages a lip-sync avatar that publishes video to a room.
//
// Implementations of this interface handle provider-specific API calls
// to create avatar sessions and manage their lifecycle.
//
// The typical flow is:
//  1. Create a Session with provider-specific configuration
//  2. Call Start() to initialize the avatar and connect it to the room
//  3. Call WaitForJoin() to wait for the avatar to be ready
//  4. Use AudioOutput() to stream TTS audio to the avatar
//  5. Call Close() when done to clean up resources
type Session interface {
	// Identity returns the participant identity of the avatar worker.
	// This is the identity that appears in the room and publishes video.
	Identity() string

	// Provider returns the provider name (e.g., "heygen", "tavus", "bithuman").
	Provider() string

	// Start initializes the avatar session with platform-specific options.
	//
	// The opts parameter is platform-specific. For LiveKit integration,
	// pass a *LiveKitStartOptions. Other platforms may define their own
	// options type.
	//
	// This method should:
	//  1. Create a session with the avatar provider's API
	//  2. Generate credentials for the avatar to join the room
	//  3. Configure the audio output to stream to the avatar
	//
	// The avatar will join the room asynchronously. Use WaitForJoin()
	// to wait for the avatar to be ready.
	Start(ctx context.Context, opts any) error

	// WaitForJoin blocks until the avatar participant joins the room
	// and publishes the expected tracks (typically video).
	//
	// Returns an error if the timeout is exceeded or the context is cancelled.
	WaitForJoin(ctx context.Context, timeout time.Duration) error

	// AudioOutput returns the audio destination for streaming TTS audio
	// to the avatar. Returns nil if the session is not started.
	AudioOutput() AudioDestination

	// Close disconnects the avatar and cleans up resources.
	//
	// This method should:
	//  1. End the session with the avatar provider
	//  2. Remove the avatar participant from the room
	//  3. Clean up any registered handlers
	Close(ctx context.Context) error

	// SetCallbacks registers event callbacks for the session.
	SetCallbacks(callbacks *SessionCallbacks)
}

// SessionCallbacks defines optional event callbacks for avatar sessions.
type SessionCallbacks struct {
	// OnAvatarJoined is called when the avatar participant joins the room.
	OnAvatarJoined func(identity string)

	// OnAvatarLeft is called when the avatar participant leaves the room.
	OnAvatarLeft func(identity string)

	// OnPlaybackStarted is called when the avatar starts speaking.
	OnPlaybackStarted func()

	// OnPlaybackFinished is called when the avatar finishes speaking.
	// The position is the playback position in seconds when stopped.
	// The interrupted flag indicates if playback was interrupted by
	// a ClearBuffer() call.
	OnPlaybackFinished func(position float64, interrupted bool)

	// OnError is called when an error occurs.
	// This is for non-fatal errors that don't cause the session to fail.
	OnError func(err error)
}

// Metrics contains avatar performance metrics.
type Metrics struct {
	// AvatarJoinLatency is the time from Start() to avatar joining the room.
	AvatarJoinLatency time.Duration

	// PlaybackLatency is the time from audio send to avatar speech start.
	// This is measured per utterance.
	PlaybackLatency time.Duration

	// Provider is the avatar provider name.
	Provider string

	// Timestamp is when the metrics were collected.
	Timestamp time.Time
}
