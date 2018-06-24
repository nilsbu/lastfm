package io

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestFailReader(t *testing.T) {
	r := FailIO{}
	data, err := r.Read(rsrc.APIKey())
	if err == nil {
		t.Error("expected error but none occurred")
	} else {
		if f, ok := err.(fail.Threat); ok {
			if f.Severity() != fail.Control {
				t.Error("severity must be 'Control':", err)
			}
		} else {
			t.Error("error must implement Threat but does not:", err)
		}
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
	} else {
		if f, ok := err.(fail.Threat); ok {
			if f.Severity() != fail.Control {
				t.Error("severity must be 'Control':", err)
			}
		} else {
			t.Error("error must implement Threat but does not:", err)
		}
	}
}
