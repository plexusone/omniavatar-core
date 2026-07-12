# OmniAvatar Core

[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/omniavatar-core
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/omniavatar-core
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/omniavatar-core/blob/main/LICENSE

Core interfaces for real-time AI avatars. This package provides the interface definitions without any provider dependencies.

For a batteries-included package with all providers, see [omniavatar](https://github.com/plexusone/omniavatar).

## Architecture

```
omniavatar-core/              # Core interfaces (no provider deps)
├── avatar/
│   ├── provider.go           # Provider interface
│   ├── session.go            # Session interface + callbacks
│   ├── audio.go              # AudioDestination interface
│   └── errors.go             # Error types
└── registry/
    └── registry.go           # Factory types + options

omniavatar/                   # Provider implementations
├── registry.go               # Global registry
├── providers/
│   ├── heygen/               # HeyGen LiveAvatar
│   ├── tavus/                # Tavus Conversational Video
│   ├── bithuman/             # bitHuman Real-time Avatars
│   └── all/                  # Convenience import
└── go.mod
```

## Interfaces

### Provider

Creates avatar sessions with provider-specific configuration.

```go
type Provider interface {
    Name() string
    CreateSession(cfg SessionConfig) (Session, error)
}
```

### Session

Manages the avatar lifecycle: start, audio streaming, and cleanup.

```go
type Session interface {
    Identity() string
    Provider() string
    Start(ctx context.Context, opts any) error
    WaitForJoin(ctx context.Context, timeout time.Duration) error
    AudioOutput() AudioDestination
    Close(ctx context.Context) error
    SetCallbacks(callbacks *SessionCallbacks)
}
```

### AudioDestination

Streams TTS audio to the avatar for lip-sync playback.

```go
type AudioDestination interface {
    CaptureFrame(ctx context.Context, frame []byte) error
    Flush(ctx context.Context) error
    ClearBuffer(ctx context.Context) error
    SampleRate() int
    Channels() int
    Close() error
}
```

## Usage

Import only the interfaces:

```go
import "github.com/plexusone/omniavatar-core/avatar"

func processAvatar(session avatar.Session) error {
    audioOut := session.AudioOutput()
    return audioOut.CaptureFrame(ctx, pcmData)
}
```

Import with all providers (batteries-included):

```go
import (
    "github.com/plexusone/omniavatar"
    _ "github.com/plexusone/omniavatar/providers/all"
)

provider, err := omniavatar.GetAvatarProvider("heygen",
    omniavatar.WithAPIKey(os.Getenv("HEYGEN_API_KEY")),
    omniavatar.WithExtension("avatar_id", avatarID),
    omniavatar.WithExtension("sandbox", true),
)
```

## Provider Registry Pattern

Providers register with a priority level:

| Priority | Constant | Description |
|----------|----------|-------------|
| 0 | `PriorityThin` | Minimal implementations |
| 10 | `PriorityThick` | Full SDK implementations |

Higher priority providers override lower priority registrations for the same name.

### Registration via init()

Provider packages use `init()` to auto-register when imported:

```go
// In omniavatar/providers/heygen/register.go
func init() {
    omniavatar.RegisterAvatarProvider("heygen", newProvider, omniavatar.PriorityThick)
}
```

## Supported Providers

| Provider | Description | Latency |
|----------|-------------|---------|
| **HeyGen** | LiveAvatar LITE mode | ~500ms |
| **Tavus** | Conversational Video | ~300ms |
| **bitHuman** | Real-time Avatars | ~200ms |

## Session Lifecycle

```
1. Get Provider    → omniavatar.GetAvatarProvider("heygen", opts...)
2. Create Session  → provider.CreateSession(cfg)
3. Start           → session.Start(ctx, startOptions)
4. Wait for Join   → session.WaitForJoin(ctx, timeout)
5. Stream Audio    → session.AudioOutput().CaptureFrame(ctx, pcm)
6. Close           → session.Close(ctx)
```

## Session Callbacks

```go
session.SetCallbacks(&avatar.SessionCallbacks{
    OnAvatarJoined: func(identity string) {
        log.Printf("Avatar joined: %s", identity)
    },
    OnPlaybackStarted: func() {
        log.Print("Avatar started speaking")
    },
    OnPlaybackFinished: func(position float64, interrupted bool) {
        log.Printf("Avatar finished speaking at %.2fs (interrupted: %v)", position, interrupted)
    },
    OnError: func(err error) {
        log.Printf("Avatar error: %v", err)
    },
})
```

## Audio Format

Default audio configuration for avatar providers:

| Parameter | Value |
|-----------|-------|
| Sample Rate | 24000 Hz |
| Channels | 1 (mono) |
| Encoding | PCM16 (linear16) |

## Resources

- [HeyGen LiveAvatar](https://liveavatar.com/)
- [Tavus](https://www.tavus.io/)
- [bitHuman](https://www.bithuman.io/)
