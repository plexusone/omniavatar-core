# Architecture

`omniavatar-core` is the interfaces-only half of a two-module design,
following the same pattern as
[OmniVoice](https://github.com/plexusone/omnivoice-core).

## Core vs. Batteries-Included

| Module | Contents |
|--------|----------|
| **omniavatar-core** (this module) | Interfaces only (`live`, `render`, `registry`) — no provider dependencies |
| [omniavatar](https://github.com/plexusone/omniavatar) | Provider implementations (HeyGen, Tavus, bitHuman) with auto-registration; pulls in provider SDKs and LiveKit |

Depend on `omniavatar-core` when you consume avatars behind an interface
and want a minimal dependency footprint; depend on `omniavatar` when you
construct providers.

```
omniavatar-core/              # Core interfaces (no provider deps)
├── live/                     # Real-time sessions
│   ├── provider.go           # live.Provider
│   ├── session.go            # live.Session + callbacks
│   ├── audio.go              # live.AudioDestination
│   └── errors.go
├── render/                   # Batch video generation
│   ├── provider.go           # render.Provider, render.AudioUploader
│   ├── request.go            # render.GenerateRequest
│   ├── job.go                # render.Job, JobStatus, render.Wait
│   └── errors.go
└── registry/                 # Shared config + factory types
```

## The Two Surfaces

The packages are named after the *mode*, not the vendor concept:

- **`live`** — session-oriented: `Start`, `WaitForJoin`, stream PCM
  frames, `Close`. The avatar joins a room and speaks in real time.
- **`render`** — job-oriented: `Generate`, `Status` / `Wait`,
  `Download`. Audio in, MP4 out, minutes later.

`avatar` is deliberately *not* a package name: with two surfaces,
`live.Session` and `render.Job` are self-describing at every call site,
while `avatar.X` would not say which mode is involved.

## Registry

`registry` provides the shared `ProviderConfig` and the two factory
types the batteries-included module wires up:

```go
type LiveProviderFactory   func(config ProviderConfig) (live.Provider, error)
type RenderProviderFactory func(config ProviderConfig) (render.Provider, error)
```

Providers self-register both surfaces via `init()` in the `omniavatar`
module, with thin/thick priority so alternative implementations can
coexist. See the
[omniavatar Architecture guide](https://plexusone.dev/omniavatar/guides/architecture/)
for the full registry and capability-interface patterns.
