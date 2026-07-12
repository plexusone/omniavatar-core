// Package omniavatar provides a unified, provider-agnostic interface for real-time AI avatars.
//
// This is the core package that defines interfaces without any provider dependencies.
// For the batteries-included package with all providers, use github.com/plexusone/omniavatar.
//
// # Architecture
//
// The omniavatar-core package provides:
//
//   - Provider interface: Creates avatar sessions
//   - Session interface: Manages avatar lifecycle and audio streaming
//   - AudioDestination interface: Streams audio to avatars
//   - Registry types: Factory types for provider registration
//
// # Usage
//
// Import omniavatar-core for the interfaces only:
//
//	import "github.com/plexusone/omniavatar-core/avatar"
//
// Import omniavatar for batteries-included experience:
//
//	import (
//	    "github.com/plexusone/omniavatar"
//	    _ "github.com/plexusone/omniavatar/providers/all"
//	)
//
//	provider, err := omniavatar.GetAvatarProvider("heygen",
//	    omniavatar.WithAPIKey(os.Getenv("HEYGEN_API_KEY")),
//	    omniavatar.WithExtension("avatar_id", avatarID))
package omniavatar
