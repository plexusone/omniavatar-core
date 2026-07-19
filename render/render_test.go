package render

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

// fakeProvider returns scripted statuses in order, then repeats the last.
type fakeProvider struct {
	statuses []JobStatus
	calls    int
	err      error
}

func (f *fakeProvider) Name() string { return "fake" }

func (f *fakeProvider) Generate(_ context.Context, _ GenerateRequest) (*Job, error) {
	return &Job{ID: "job-1", Provider: "fake"}, nil
}

func (f *fakeProvider) Status(_ context.Context, _ string) (*JobStatus, error) {
	if f.err != nil {
		return nil, f.err
	}
	i := f.calls
	if i >= len(f.statuses) {
		i = len(f.statuses) - 1
	}
	f.calls++
	status := f.statuses[i]
	return &status, nil
}

func (f *fakeProvider) Download(_ context.Context, _ string, _ io.Writer) error {
	return nil
}

func TestJobStateTerminal(t *testing.T) {
	tests := []struct {
		state JobState
		want  bool
	}{
		{JobStatePending, false},
		{JobStateProcessing, false},
		{JobStateCompleted, true},
		{JobStateFailed, true},
	}
	for _, tt := range tests {
		if got := tt.state.Terminal(); got != tt.want {
			t.Errorf("JobState(%q).Terminal() = %v, want %v", tt.state, got, tt.want)
		}
	}
}

func TestGenerateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     GenerateRequest
		wantErr bool
	}{
		{"audio ok", GenerateRequest{AvatarID: "a", AudioURL: "https://x/a.mp3"}, false},
		{"script ok", GenerateRequest{AvatarID: "a", Script: "hello"}, false},
		{"missing avatar", GenerateRequest{AudioURL: "https://x/a.mp3"}, true},
		{"missing input", GenerateRequest{AvatarID: "a"}, true},
		{"both inputs", GenerateRequest{AvatarID: "a", AudioURL: "https://x/a.mp3", Script: "hi"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrInvalidRequest) {
				t.Fatalf("Validate() error = %v, want errors.Is ErrInvalidRequest", err)
			}
		})
	}
}

func TestWaitCompletes(t *testing.T) {
	p := &fakeProvider{statuses: []JobStatus{
		{ID: "job-1", State: JobStatePending},
		{ID: "job-1", State: JobStateProcessing},
		{ID: "job-1", State: JobStateCompleted, VideoURL: "https://x/v.mp4"},
	}}

	status, err := Wait(context.Background(), p, "job-1", time.Millisecond)
	if err != nil {
		t.Fatalf("Wait() error = %v", err)
	}
	if status.State != JobStateCompleted {
		t.Errorf("Wait() state = %q, want %q", status.State, JobStateCompleted)
	}
	if status.VideoURL != "https://x/v.mp4" {
		t.Errorf("Wait() VideoURL = %q, want set", status.VideoURL)
	}
	if p.calls != 3 {
		t.Errorf("Status called %d times, want 3", p.calls)
	}
}

func TestWaitFails(t *testing.T) {
	p := &fakeProvider{statuses: []JobStatus{
		{ID: "job-1", State: JobStateFailed, ErrorCode: "E1", ErrorMsg: "boom"},
	}}

	status, err := Wait(context.Background(), p, "job-1", time.Millisecond)
	if !errors.Is(err, ErrJobFailed) {
		t.Fatalf("Wait() error = %v, want errors.Is ErrJobFailed", err)
	}
	if status == nil || status.ErrorCode != "E1" {
		t.Errorf("Wait() status = %+v, want failed status with ErrorCode", status)
	}
}

func TestWaitContextCancelled(t *testing.T) {
	p := &fakeProvider{statuses: []JobStatus{
		{ID: "job-1", State: JobStateProcessing},
	}}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := Wait(ctx, p, "job-1", time.Hour)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Wait() error = %v, want context.DeadlineExceeded", err)
	}
}

func TestWaitStatusError(t *testing.T) {
	wantErr := errors.New("api down")
	p := &fakeProvider{err: wantErr}

	_, err := Wait(context.Background(), p, "job-1", time.Millisecond)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Wait() error = %v, want %v", err, wantErr)
	}
}

func TestProviderErrorFormat(t *testing.T) {
	underlying := errors.New("boom")
	err := NewProviderError("heygen", "generate", underlying)
	want := "render/heygen: generate: boom"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
	if !errors.Is(err, underlying) {
		t.Error("errors.Is(err, underlying) = false, want true")
	}
}
