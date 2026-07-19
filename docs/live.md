# Live Interfaces

Package `live` defines the real-time streaming avatar surface: create a
session, connect it to a room, and stream TTS audio for lip-sync
playback. This is the counterpart to the [render](render.md) (batch)
surface.

Full godoc: [pkg.go.dev/.../live](https://pkg.go.dev/github.com/plexusone/omniavatar-core/live).

## Provider

Creates avatar sessions with provider-specific configuration.

```go
type Provider interface {
    Name() string
    CreateSession(cfg SessionConfig) (Session, error)
}
```

## Session

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

The `opts` passed to `Start` are platform-specific. For LiveKit, the
[omniavatar](https://github.com/plexusone/omniavatar) package provides
`LiveKitStartOptions` and token generation.

## AudioDestination

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

## Session Callbacks

```go
session.SetCallbacks(&live.SessionCallbacks{
    OnAvatarJoined:     func(identity string) { /* ... */ },
    OnPlaybackStarted:  func() { /* ... */ },
    OnPlaybackFinished: func(position float64, interrupted bool) { /* ... */ },
    OnError:            func(err error) { /* ... */ },
})
```

## Session Lifecycle

```
1. Get Provider    → omniavatar.GetLiveProvider("heygen", opts...)
2. Create Session  → provider.CreateSession(cfg)
3. Start           → session.Start(ctx, startOptions)
4. Wait for Join   → session.WaitForJoin(ctx, timeout)
5. Stream Audio    → session.AudioOutput().CaptureFrame(ctx, pcm)
6. Close           → session.Close(ctx)
```

## Audio Format

Default audio configuration (`live.DefaultAudioConfig()`):

| Parameter | Value |
|-----------|-------|
| Sample Rate | 24000 Hz |
| Channels | 1 (mono) |
| Encoding | PCM16 (linear16) |
