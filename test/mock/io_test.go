package mock

import (
	"errors"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// Resource that has no path.
type noPath string

func (n noPath) URL(apiKey rsrc.Key) (string, error) {
	return string(n), nil
}

func (n noPath) Path() (string, error) {
	return "", errors.New("resource has no path")
}

func TestFileIO(t *testing.T) {
	apiKeyPath, _ := rsrc.APIKey().Path()

	cases := []struct {
		files     map[string][]byte
		rs        rsrc.Resource
		writeData []byte
		result    []byte
		writeOK   bool
		readOK    bool
	}{
		{ // no data
			map[string][]byte{},
			rsrc.APIKey(),
			nil, nil,
			false, false,
		},
		{ // read prepared data
			map[string][]byte{apiKeyPath: []byte("xxd")},
			rsrc.APIKey(),
			nil,
			[]byte("xxd"),
			true, true,
		},
		{ // read what was written
			map[string][]byte{apiKeyPath: nil},
			rsrc.APIKey(),
			[]byte("xxd"),
			[]byte("xxd"),
			true, true,
		},
		{ // write fails (key not contained in files)
			map[string][]byte{},
			rsrc.APIKey(),
			[]byte("xxd"),
			nil,
			false, false,
		},
		{ // read from nil is not possible
			map[string][]byte{apiKeyPath: nil},
			rsrc.APIKey(),
			nil,
			nil,
			false, false,
		},
		{ // resolve of file path fails
			map[string][]byte{},
			noPath(""),
			[]byte(""),
			[]byte(""),
			false, false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			r, w := FileIO(c.files)

			if c.writeData != nil {
				err := w.Write(c.writeData, c.rs)
				if err != nil && c.writeOK {
					t.Error("unexpected error during write:", err)
				} else if err == nil && !c.writeOK {
					t.Error("write should have failed but did not")
				}
			}

			data, err := r.Read(c.rs)
			close(r)
			close(w)

			if err != nil && c.readOK {
				t.Error("unexpected error during read:", err)
			} else if err == nil && !c.readOK {
				t.Error("read should have failed but did not")
			}
			if err == nil {
				if string(data) != string(c.result) {
					t.Errorf("result does not match:\nresult:   %v\nexpected: %v",
						string(data), string(c.result))
				}
			}
		})
	}
}

func TestDownloader(t *testing.T) {
	userInfo, _ := rsrc.UserInfo("abc")
	userInfoURL, _ := userInfo.URL(APIKey)

	cases := []struct {
		files  map[string][]byte
		rs     rsrc.Resource
		result []byte
		readOK bool
	}{
		{ // API key has no URL
			map[string][]byte{},
			rsrc.APIKey(),
			nil,
			false,
		},
		{
			map[string][]byte{userInfoURL: []byte("xxx")},
			userInfo,
			[]byte("xxx"),
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			r := Downloader(c.files)

			data, err := r.Read(c.rs)
			close(r)

			if err != nil && c.readOK {
				t.Error("unexpected error during read:", err)
			} else if err == nil && !c.readOK {
				t.Error("read should have failed but did not")
			}
			if err == nil {
				if string(data) != string(c.result) {
					t.Errorf("result does not match:\nresult:   %v\nexpected: %v",
						string(data), string(c.result))
				}
			}
		})
	}
}

func TestAsyncFileIO(t *testing.T) {
	apiKeyPath, _ := rsrc.APIKey().Path()

	cases := []struct {
		files     map[string][]byte
		rs        rsrc.Resource
		writeData []byte
		result    []byte
		writeOK   bool
		readOK    bool
	}{
		{ // read what was written
			map[string][]byte{apiKeyPath: nil},
			rsrc.APIKey(),
			[]byte("xxd"),
			[]byte("xxd"),
			true, true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			r, w := AsyncFileIO(c.files)

			if c.writeData != nil {
				err := <-w.Write(c.writeData, c.rs)
				if err != nil && c.writeOK {
					t.Error("unexpected error during write:", err)
				} else if err == nil && !c.writeOK {
					t.Error("write should have failed but did not")
				}
			}

			res := <-r.Read(c.rs)
			data, err := res.Data, res.Err
			close(r)
			close(w)

			if err != nil && c.readOK {
				t.Error("unexpected error during read:", err)
			} else if err == nil && !c.readOK {
				t.Error("read should have failed but did not")
			}
			if err == nil {
				if string(data) != string(c.result) {
					t.Errorf("result does not match:\nresult:   %v\nexpected: %v",
						string(data), string(c.result))
				}
			}
		})
	}
}

func TestAsyncDownloader(t *testing.T) {
	userInfo, _ := rsrc.UserInfo("abc")
	userInfoURL, _ := userInfo.URL(APIKey)

	cases := []struct {
		files  map[string][]byte
		rs     rsrc.Resource
		result []byte
		readOK bool
	}{
		{
			map[string][]byte{userInfoURL: []byte("xxx")},
			userInfo,
			[]byte("xxx"),
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			r := AsyncDownloader(c.files)

			res := <-r.Read(c.rs)
			data, err := res.Data, res.Err
			close(r)

			if err != nil && c.readOK {
				t.Error("unexpected error during read:", err)
			} else if err == nil && !c.readOK {
				t.Error("read should have failed but did not")
			}
			if err == nil {
				if string(data) != string(c.result) {
					t.Errorf("result does not match:\nresult:   %v\nexpected: %v",
						string(data), string(c.result))
				}
			}
		})
	}
}
