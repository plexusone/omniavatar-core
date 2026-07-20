# Architecture

`omniavatar-core` is the interfaces-only core of a multi-module design,
following the same pattern as
[OmniVoice](https://github.com/plexusone/omnivoice-core).

## Where adapters live

| Module | Contents |
|--------|----------|
| **omniavatar-core** (this module) | Interfaces (`live`, `render`, `registry`) + small stdlib-only render helpers — no provider dependencies |
| `heygen-go/omniavatar`, `tavus-go/omniavatar`, `bithuman-go/omniavatar` | **Render** adapters, hosted in each provider SDK repo, depending only on `omniavatar-core`, so provider-specific knowledge stays with the SDK |
| [omniavatar](https://github.com/plexusone/omniavatar) | Batteries-included: the **live** (LiveKit-coupled) adapters, the registries, and `providers/all` which registers every provider |

Render adapters live in the provider SDK repos (the PlexusOne convention);
live adapters live in the batteries package because their LiveKit
integration does. Depend on `omniavatar-core` when you consume avatars
behind an interface; depend on `omniavatar` (or a specific
`<sdk>/omniavatar` render adapter) when you construct providers.

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

The batteries `omniavatar` package wires these up: each
`providers/<name>` registers its local **live** adapter and the SDK-hosted
**render** adapter (constructor-based, `<sdk>/omniavatar.NewRenderProviderFromConfig`)
via `init()`, with thin/thick priority so alternative implementations can
coexist. See the
[omniavatar Architecture guide](https://plexusone.dev/omniavatar/guides/architecture/)
for the full registry and capability-interface patterns.
