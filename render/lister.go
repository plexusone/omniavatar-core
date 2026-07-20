package render

import "context"

// AvatarInfo describes an avatar available from a render provider,
// discovered via AvatarLister.
type AvatarInfo struct {
	// ID is directly usable as GenerateRequest.AvatarID with the same
	// provider.
	ID string

	// Name is the human-readable avatar name.
	Name string

	// Gender is the avatar gender, when the provider reports it.
	Gender string
}

// AvatarLister is an optional capability for render providers that can
// enumerate the avatars available to the account. Callers feature-detect
// it, like AudioUploader:
//
//	if l, ok := provider.(render.AvatarLister); ok {
//	    avatars, err := l.ListAvatars(ctx, "")
//	}
//
// The returned AvatarInfo.ID values are directly usable as
// GenerateRequest.AvatarID with the same provider — useful because a
// provider's avatar-listing endpoint may return different identifiers
// than its generation endpoint accepts (HeyGen's v3 avatar groups vs. v2
// avatar IDs, for example).
type AvatarLister interface {
	// ListAvatars returns avatars whose ID or name matches search
	// (case-insensitive substring; an empty search returns all).
	ListAvatars(ctx context.Context, search string) ([]AvatarInfo, error)
}
