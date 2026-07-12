package avatar

// Provider creates avatar sessions.
//
// Implementations of this interface handle provider-specific configuration
// and session creation. Each provider (HeyGen, Tavus, bitHuman, etc.) has
// its own Provider implementation.
//
// Example usage:
//
//	provider, err := heygen.NewProvider(heygen.Config{
//	    APIKey:   os.Getenv("HEYGEN_API_KEY"),
//	    AvatarID: os.Getenv("HEYGEN_AVATAR_ID"),
//	})
//	if err != nil {
//	    return err
//	}
//
//	session, err := provider.CreateSession(avatar.SessionConfig{
//	    AudioConfig: avatar.DefaultAudioConfig(),
//	})
type Provider interface {
	// Name returns the provider name (e.g., "heygen", "tavus", "bithuman").
	Name() string

	// CreateSession creates a new avatar session with the given config.
	// The session is not started until Session.Start() is called.
	CreateSession(cfg SessionConfig) (Session, error)
}

// SessionConfig contains configuration for creating an avatar session.
type SessionConfig struct {
	// AudioConfig specifies the audio format requirements.
	// Default: 24kHz mono PCM16
	AudioConfig AudioConfig

	// Extensions holds provider-specific configuration.
	// Keys and values depend on the provider.
	Extensions map[string]any
}
