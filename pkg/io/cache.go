package io

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type cacheServer struct {
	srv  *http.Server
	data *sync.Map
	// TODO add max size
}

func (c *cacheServer) read(rs string) ([]byte, error) {
	data, ok := c.data.Load(rs)
	if !ok {
		return nil, fail.WrapError(fail.Control, fmt.Errorf("cannot read '%v'", rs))
	}
	return data.([]byte), nil
}

func (c *cacheServer) write(data []byte, rs string) {
	c.data.Store(rs, data)
}

func (c *cacheServer) remove(rs string) {
	c.data.Delete(rs)
}

func (c *cacheServer) handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	rsrcs, ok := r.PostForm["rsrc"]
	if !ok && len(rsrcs) != 1 {
		w.Header().Set("err", "no resource locator provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rs := rsrcs[0]

	actions, ok := r.PostForm["action"]
	if !ok && len(actions) != 1 {
		w.Header().Set("err", "no action provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch actions[0] {
	case "read":
		if data, err := c.read(rs); err != nil {
			w.Header().Set("err", fmt.Sprintf("resource '%v' not found", rs))
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Write(data)
		}
	case "write":
		if data, ok := r.PostForm["data"]; !ok || len(data) == 0 {
			w.Header().Set("err", "no data provided")
			w.WriteHeader(http.StatusBadRequest)
		} else {
			c.write([]byte(data[0]), rs)
		}
	case "remove":
		c.remove(rs)
	}
}

func RunCacheServer(port int) {
	cache := &cacheServer{data: &sync.Map{}}
	cache.srv = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: http.HandlerFunc(cache.handle)}

	cache.srv.ListenAndServe()
}

type CacheIO struct {
	Port int
}

func (io *CacheIO) Read(loc rsrc.Locator) ([]byte, error) {
	resp, err := io.postForm(loc, "read", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			err = fail.WrapError(fail.Control, errors.New(resp.Header["Err"][0]))
		default:
			// does not occurr
		}
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	return data, err
}

func (io *CacheIO) Write(data []byte, loc rsrc.Locator) error {
	resp, err := io.postForm(loc, "write", data)
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

func (io *CacheIO) Remove(loc rsrc.Locator) error {
	resp, err := io.postForm(loc, "remove", nil)
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

func (io *CacheIO) postForm(
	loc rsrc.Locator,
	action string,
	data []byte,
) (*http.Response, error) {
	path, err := loc.Path()
	if err != nil {
		return nil, err
	}

	params := url.Values{"rsrc": {path}, "action": {action}}
	if data != nil {
		params["data"] = []string{string(data)}
	}

	url := fmt.Sprintf("http://localhost:%v", io.Port)

	resp, err := http.PostForm(url, params)
	if err != nil {
		err = fail.WrapError(fail.Critical, err)
	}

	return resp, err
}
