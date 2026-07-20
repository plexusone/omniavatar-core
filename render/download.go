package render

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// DownloadURL streams the content at url to dst using client, a helper for
// implementing Provider.Download from a completed job's video URL. A nil
// client uses http.DefaultClient. Cancellation is via ctx.
func DownloadURL(ctx context.Context, client *http.Client, url string, dst io.Writer) error {
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create download request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download video: %w", err)
	}
	defer func() {
		// Close error after a completed copy is unactionable.
		_ = resp.Body.Close() //nolint:errcheck // see above
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download video: unexpected status %s", resp.Status)
	}

	if _, err := io.Copy(dst, resp.Body); err != nil {
		return fmt.Errorf("write video: %w", err)
	}
	return nil
}
