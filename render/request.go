package render

import "fmt"

// GenerateRequest describes an avatar video generation job.
//
// Exactly one of AudioURL or Script must be set. AudioURL is the design
// center: it drives lip-sync from existing narration audio so the avatar
// matches the authoritative audio track exactly. Script uses provider
// TTS instead and is a secondary path.
type GenerateRequest struct {
	// AvatarID identifies the presenter with the provider.
	// HeyGen: avatar_id; Tavus: replica_id; bitHuman: agent_id.
	AvatarID string

	// AudioURL is a fetchable URL to narration audio (.mp3/.wav).
	// Providers implementing AudioUploader can host local files.
	AudioURL string

	// Script is text for provider TTS. Providers that require a voice
	// read it from Extensions (e.g., "voice_id").
	Script string

	// Width and Height are the requested output dimensions in pixels.
	// Optional; provider defaults apply when zero.
	Width int

	// Height is the requested output height in pixels.
	Height int

	// Background requests a background treatment. Optional and
	// best-effort; support varies by provider.
	Background *Background

	// Title is a human-readable job/video name. Optional.
	Title string

	// Extensions holds provider-specific options
	// (e.g., "voice_id", "avatar_style", "test", "fast").
	Extensions map[string]any
}

// Background describes the requested video background.
type Background struct {
	// Type is "color", "image", or "video".
	Type string

	// Value is a hex color or URL depending on Type.
	Value string
}

// Validate checks that the request is well-formed. It returns an error
// wrapping ErrInvalidRequest if AvatarID is empty, or if neither or both
// of AudioURL and Script are set.
func (r GenerateRequest) Validate() error {
	if r.AvatarID == "" {
		return fmt.Errorf("%w: AvatarID is required", ErrInvalidRequest)
	}
	if r.AudioURL == "" && r.Script == "" {
		return fmt.Errorf("%w: one of AudioURL or Script is required", ErrInvalidRequest)
	}
	if r.AudioURL != "" && r.Script != "" {
		return fmt.Errorf("%w: AudioURL and Script are mutually exclusive", ErrInvalidRequest)
	}
	return nil
}

// GetString retrieves a string extension value with a default.
func (r GenerateRequest) GetString(key, defaultValue string) string {
	if v, ok := r.Extensions[key].(string); ok {
		return v
	}
	return defaultValue
}

// GetBool retrieves a bool extension value with a default.
func (r GenerateRequest) GetBool(key string, defaultValue bool) bool {
	if v, ok := r.Extensions[key].(bool); ok {
		return v
	}
	return defaultValue
}
