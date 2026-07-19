// Package omniavatar provides unified, provider-agnostic interfaces for AI avatars.
//
// This is the core package that defines interfaces without any provider dependencies.
// For the batteries-included package with all providers, use github.com/plexusone/omniavatar.
//
// # Architecture
//
// The omniavatar-core module provides two surfaces:
//
//   - live: real-time streaming avatar sessions (Provider, Session,
//     AudioDestination) for conversational agents
//   - render: asynchronous batch avatar video generation (Provider, Job,
//     AudioUploader) for offline pipelines such as presentation videos
//   - registry: shared configuration and factory types for provider
//     registration
//
// # Usage
//
// Import omniavatar-core for the interfaces only:
//
//	import (
//	    "github.com/plexusone/omniavatar-core/live"
//	    "github.com/plexusone/omniavatar-core/render"
//	)
//
// Import omniavatar for the batteries-included experience:
//
//	import (
//	    "github.com/plexusone/omniavatar"
//	    _ "github.com/plexusone/omniavatar/providers/all"
//	)
//
//	liveProvider, err := omniavatar.GetLiveProvider("heygen",
//	    omniavatar.WithAPIKey(os.Getenv("LIVEAVATAR_API_KEY")),
//	    omniavatar.WithExtension("avatar_id", avatarID))
//
//	renderProvider, err := omniavatar.GetRenderProvider("heygen",
//	    omniavatar.WithAPIKey(os.Getenv("HEYGEN_API_KEY")))
package omniavatar
