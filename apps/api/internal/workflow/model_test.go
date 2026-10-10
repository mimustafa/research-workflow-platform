package workflow

import (
	"testing"
	"time"
)

func TestNewCreatesQueuedWorkflow(t *testing.T) {
	now := time.Date(2026, 10, 10, 21, 0, 0, 0, time.FixedZone("CEST", 2*60*60))

	result, err := New("Summarize recent research", now)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if result.ID == "" {
		t.Fatal("New() returned an empty ID")
	}
	if result.Status != StatusQueued {
		t.Fatalf("status = %q, want %q", result.Status, StatusQueued)
	}
	if !result.CreatedAt.Equal(now.UTC()) {
		t.Fatalf("created_at = %v, want %v", result.CreatedAt, now.UTC())
	}
}

func TestNewRejectsEmptyBrief(t *testing.T) {
	_, err := New("   ", time.Now())
	if err == nil {
		t.Fatal("New() error = nil, want validation error")
	}
}
