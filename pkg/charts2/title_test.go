package charts2

import "testing"

func TestKeyTitle(t *testing.T) {
	kt := KeyTitle("a")

	if kt.String() != "a" {
		t.Errorf("String() was expected to be 'a' but is '%v'", kt.String())
	}

	if kt.Key() != "a" {
		t.Errorf("Key() was expected to be 'a' but is '%v'", kt.Key())
	}

	if kt.Artist() != "" {
		t.Errorf("Artist() was expected to be '' but is '%v'", kt.Artist())
	}

	if kt.Song() != "" {
		t.Errorf("Song() was expected to be '' but is '%v'", kt.Song())
	}
}

func TestArtistTitle(t *testing.T) {
	at := ArtistTitle("a")

	if at.String() != "a" {
		t.Errorf("String() was expected to be 'a' but is '%v'", at.String())
	}

	if at.Key() != "a" {
		t.Errorf("Key() was expected to be 'a' but is '%v'", at.Key())
	}

	if at.Artist() != "a" {
		t.Errorf("Artist() was expected to be 'a' but is '%v'", at.Artist())
	}

	if at.Song() != "" {
		t.Errorf("Song() was expected to be '' but is '%v'", at.Song())
	}
}
