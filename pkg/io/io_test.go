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
	err := r.Write([]byte("xyz"), rsrc.APIKey())
	if err == nil {
		t.Error("expected error but none occurred")
	}
}
