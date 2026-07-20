# Render Interfaces

Package `render` defines the asynchronous (batch) avatar video generation
surface: submit narration audio or a script, poll for completion,
download a talking-head video. This is the counterpart to the
[live](live.md) (real-time) surface.

The design center is **audio-driven generation** — the avatar's lip-sync
comes from the exact audio you supply, so pauses, pronunciation, and
timing match your production audio.

Full godoc: [pkg.go.dev/.../render](https://pkg.go.dev/github.com/plexusone/omniavatar-core/render).

## Provider

Generates avatar videos asynchronously: submit, poll, download.

```go
type Provider interface {
    Name() string
    Generate(ctx context.Context, req GenerateRequest) (*Job, error)
    Status(ctx context.Context, jobID string) (*JobStatus, error)
    Download(ctx context.Context, jobID string, dst io.Writer) error
}
```

## AudioUploader (optional capability)

Providers that can host local audio files implement this; feature-detect
it rather than assuming it exists:

```go
if up, ok := provider.(render.AudioUploader); ok {
    audioURL, err = up.UploadAudio(ctx, "narration.mp3", f)
} else {
    // provider cannot host files; audioURL must be supplied externally
}
```

Modeling upload as a capability keeps the core `Provider` interface at the
lowest common denominator — all providers consume audio by URL, but only
some can host it. `render.ErrAudioUploadUnsupported` is available for
callers that need a typed error.

## AvatarLister (optional capability)

Providers that can enumerate the account's avatars implement this; the
returned `AvatarInfo.ID` values are directly usable as
`GenerateRequest.AvatarID`:

```go
if l, ok := provider.(render.AvatarLister); ok {
    avatars, err := l.ListAvatars(ctx, "abigail") // case-insensitive substring; "" = all
}
```

Useful because a provider's listing endpoint may return different
identifiers than its generation endpoint accepts (e.g. HeyGen v3 avatar
groups vs. v2 avatar IDs). Added in v0.3.0.

## GenerateRequest

Exactly one of `AudioURL` (primary) or `Script` (secondary, provider TTS)
must be set.

```go
job, err := provider.Generate(ctx, render.GenerateRequest{
    AvatarID: avatarID,     // heygen avatar_id / tavus replica_id / bithuman agent_id
    AudioURL: narrationURL, // drives lip-sync from existing audio
})
```

`GenerateRequest.Validate()` enforces these rules; provider-specific
options go in `Extensions`.

## Job and Status

```go
type JobState string // pending | processing | completed | failed

func (s JobState) Terminal() bool

type JobStatus struct {
    ID           string
    State        JobState
    RawStatus    string  // provider-native status, preserved for logging
    VideoURL     string  // set when completed
    ThumbnailURL string  // when the provider reports one
    Duration     float64 // seconds, when reported
    ErrorCode    string
    ErrorMsg     string
}
```

Unknown provider states map to `processing` (non-terminal), so pollers
keep waiting rather than aborting on an unrecognized vendor status.

## Wait

Polls until a terminal state, honoring context cancellation:

```go
status, err := render.Wait(ctx, provider, job.ID, 5*time.Second)
if errors.Is(err, render.ErrJobFailed) {
    // status is non-nil: inspect status.ErrorCode / status.ErrorMsg
}
```

## Adapter Helpers

Small stdlib-only helpers for building render adapters (used by the
SDK-hosted adapters, so they don't duplicate them):

- `render.AudioContentType(filename)` — audio MIME type from a filename,
  for `AudioUploader` implementations
- `render.DownloadURL(ctx, client, url, dst)` — stream a URL to an
  `io.Writer`, for `Provider.Download` implementations

## Job Lifecycle

```
1. Get Provider    → omniavatar.GetRenderProvider("heygen", opts...)
2. Upload Audio    → provider.(render.AudioUploader).UploadAudio(...)  [optional]
3. Generate        → provider.Generate(ctx, render.GenerateRequest{...})
4. Wait            → render.Wait(ctx, provider, job.ID, interval)
5. Download        → provider.Download(ctx, job.ID, dst)
```

## Errors

| Sentinel | Meaning |
|----------|---------|
| `ErrInvalidRequest` | Request failed validation (missing `AvatarID`, neither/both of `AudioURL`/`Script`) |
| `ErrAudioUploadUnsupported` | Provider cannot host audio files |
| `ErrJobNotFound` | Provider does not recognize the job ID |
| `ErrJobFailed` | Job reached `failed` |
| `ErrJobNotCompleted` | `Download` called before successful completion |

Provider errors are wrapped in `render.ProviderError` (`render/heygen:
generate: ...`).
