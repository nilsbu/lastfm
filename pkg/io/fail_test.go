package io

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestFailReader(t *testing.T) {
	r := FailIO{}
	data, err := r.Read(rsrc.APIKey())
	if err == nil {
		t.Error("expected error but none occurred")
	}
	if data != nil {
		t.Errorf("data should be nil but was '%v'", string(data))
	}
}

func TestFailWriter(t *testing.T) {
	r := FailIO{}
	if err := r.Write([]byte("xyz"), rsrc.APIKey()); err == nil {
		t.Error("expected error but none occurred")
	}
}

func TestFailRemover(t *testing.T) {
	rm := FailIO{}
	if err := rm.Remove(rsrc.APIKey()); err == nil {
		t.Error("expected error but none occurred")
	}
}
