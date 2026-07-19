# OmniAvatar Core

**Provider-agnostic AI avatar interfaces for Go — no provider dependencies.**

`omniavatar-core` defines the interfaces for working with AI avatar
vendors, across two surfaces, without pulling in any provider SDK:

- **[`live`](live.md)** — real-time streaming avatar sessions (rooms, PCM
  audio streaming for lip-sync) for conversational agents
- **[`render`](render.md)** — asynchronous batch avatar video generation
  (narration audio in, talking-head MP4 out) for offline pipelines

For a batteries-included package with all providers wired up, use
[omniavatar](https://github.com/plexusone/omniavatar).

## When to Use Which

| Import | Use when |
|--------|----------|
| `omniavatar-core` | You accept avatars as an interface (library code, consumers) — zero provider dependencies |
| [`omniavatar`](https://github.com/plexusone/omniavatar) | You construct providers — pulls in provider SDKs and LiveKit |

## Installation

```bash
go get github.com/plexusone/omniavatar-core
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
)

renderProvider, err := omniavatar.GetRenderProvider("heygen",
    omniavatar.WithAPIKey(os.Getenv("HEYGEN_API_KEY")),
)
```

## Documentation

- [Live interfaces](live.md) — `Provider`, `Session`, `AudioDestination`
- [Render interfaces](render.md) — `Provider`, `AudioUploader`, `Wait`
- [Architecture](architecture.md) — the core / batteries-included split and registry pattern
- [API reference on pkg.go.dev](https://pkg.go.dev/github.com/plexusone/omniavatar-core) — full godoc
