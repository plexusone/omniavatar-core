package render

import (
	"path"
	"strings"
)

// AudioContentType guesses a MIME type from an audio filename extension,
// for use by AudioUploader implementations. Unknown extensions return
// "application/octet-stream".
func AudioContentType(filename string) string {
	switch strings.ToLower(path.Ext(filename)) {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".m4a", ".aac":
		return "audio/aac"
	case ".flac":
		return "audio/flac"
	default:
		return "application/octet-stream"
	}
}
