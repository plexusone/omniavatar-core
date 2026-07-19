# OmniAvatar Core

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/omniavatar-core/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/omniavatar-core/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/omniavatar-core/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/omniavatar-core/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/omniavatar-core/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/omniavatar-core/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/omniavatar-core
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/omniavatar-core
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/omniavatar
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fomniavatar-core
 [loc-svg]: https://tokei.rs/b1/github/plexusone/omniavatar-core
 [repo-url]: https://github.com/plexusone/omniavatar-core
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/omniavatar-core/blob/main/LICENSE

Core interfaces for AI avatars. This package provides the interface definitions without any provider dependencies, across two surfaces:

- **`live`** — real-time streaming avatar sessions (rooms, PCM audio streaming for lip-sync) for conversational agents
- **`render`** — asynchronous batch avatar video generation (narration audio in, talking-head MP4 out) for offline pipelines

For a batteries-included package with all providers, see [omniavatar](https://github.com/plexusone/omniavatar).

## Architecture

```
omniavatar-core/              # Core interfaces (no provider deps)
├── live/                     # Real-time sessions
│   ├── provider.go           # Provider interface
│   ├── session.go            # Session interface + callbacks
│   ├── audio.go              # AudioDestination interface
│   └── errors.go             # Error types
├── render/                   # Batch video generation
│   ├── provider.go           # Provider + AudioUploader interfaces
│   ├── request.go            # GenerateRequest
│   ├── job.go                # Job, JobState, JobStatus, Wait
│   └── errors.go             # Error types
└── registry/
    └── registry.go           # Factory types + options

omniavatar/                   # Provider implementations
├── registry.go               # Global live + render registries
├── providers/
│   ├── heygen/               # HeyGen (live + render)
│   ├── tavus/                # Tavus (live + render)
│   ├── bithuman/             # bitHuman (live + render)
│   └── all/                  # Convenience import
└── go.mod
```

## Live Interfaces

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

## Render Interfaces

### Provider

Generates avatar videos asynchronously: submit, poll, download.

```go
type Provider interface {
    Name() string
    Generate(ctx context.Context, req GenerateRequest) (*Job, error)
    Status(ctx context.Context, jobID string) (*JobStatus, error)
    Download(ctx context.Context, jobID string, dst io.Writer) error
}
```

### AudioUploader (optional capability)

Providers that can host local audio files implement this; feature-detect it:

```go
if up, ok := provider.(render.AudioUploader); ok {
    audioURL, err = up.UploadAudio(ctx, "narration.mp3", f)
}
```

### GenerateRequest

```go
job, err := provider.Generate(ctx, render.GenerateRequest{
    AvatarID: avatarID,        // heygen avatar_id / tavus replica_id / bithuman agent_id
    AudioURL: narrationURL,    // drives lip-sync from existing audio (primary path)
})
```

### Wait

```go
status, err := render.Wait(ctx, provider, job.ID, 5*time.Second)
// status.State: pending → processing → completed | failed
```

## Usage

Import only the interfaces:

```go
import (
    "github.com/plexusone/omniavatar-core/live"
    "github.com/plexusone/omniavatar-core/render"
)

func processAvatar(session live.Session) error {
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

liveProvider, err := omniavatar.GetLiveProvider("heygen",
    omniavatar.WithAPIKey(os.Getenv("LIVEAVATAR_API_KEY")),
    omniavatar.WithExtension("avatar_id", avatarID),
    omniavatar.WithExtension("sandbox", true),
)

renderProvider, err := omniavatar.GetRenderProvider("heygen",
    omniavatar.WithAPIKey(os.Getenv("HEYGEN_API_KEY")),
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
    omniavatar.RegisterLiveProvider("heygen", NewProviderFromConfig, omniavatar.PriorityThick)
    omniavatar.RegisterRenderProvider("heygen", NewRenderProviderFromConfig, omniavatar.PriorityThick)
}
```

## Supported Providers

| Provider | Live | Render | Live Latency |
|----------|------|--------|--------------|
| **HeyGen** | LiveAvatar LITE mode | Video Generation v2 | ~500ms |
| **Tavus** | Conversational Video | Video Generation (replicas) | ~300ms |
| **bitHuman** | Real-time Avatars | Video Generation + audio upload | ~200ms |

## Session Lifecycle (live)

```
1. Get Provider    → omniavatar.GetLiveProvider("heygen", opts...)
2. Create Session  → provider.CreateSession(cfg)
3. Start           → session.Start(ctx, startOptions)
4. Wait for Join   → session.WaitForJoin(ctx, timeout)
5. Stream Audio    → session.AudioOutput().CaptureFrame(ctx, pcm)
6. Close           → session.Close(ctx)
```

## Job Lifecycle (render)

```
1. Get Provider    → omniavatar.GetRenderProvider("heygen", opts...)
2. Upload Audio    → provider.(render.AudioUploader).UploadAudio(...)  [optional]
3. Generate        → provider.Generate(ctx, render.GenerateRequest{...})
4. Wait            → render.Wait(ctx, provider, job.ID, interval)
5. Download        → provider.Download(ctx, job.ID, dst)
```

## Session Callbacks (live)

```go
session.SetCallbacks(&live.SessionCallbacks{
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

## Audio Format (live)

Default audio configuration for avatar providers:

| Parameter | Value |
|-----------|-------|
| Sample Rate | 24000 Hz |
| Channels | 1 (mono) |
| Encoding | PCM16 (linear16) |

## Resources

- [omniavatar](https://github.com/plexusone/omniavatar) - Provider implementations
- [HeyGen LiveAvatar](https://liveavatar.com/)
- [HeyGen API](https://docs.heygen.com/)
- [Tavus](https://www.tavus.io/)
- [bitHuman](https://www.bithuman.io/)
