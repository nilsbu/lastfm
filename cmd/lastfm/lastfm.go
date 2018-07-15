package main

import (
	"fmt"
	"os"

	"github.com/nilsbu/lastfm/pkg/command"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

func isCacheRunning(port int) bool {
	io := io.CacheIO{Port: port}
	_, err := io.Read(rsrc.SessionInfo())
	if err == nil {
		return true
	} else if f, ok := err.(fail.Threat); !ok || f.Severity() > fail.Control {
		return false
	}
	return true
}

func createStore() (store.Store, error) {
	key, err := unpack.LoadAPIKey(io.FileIO{})
	if err != nil {
		return nil, err
	}

	ios := [][]rsrc.IO{}

	var webIOs []rsrc.IO
	for i := 0; i < 32; i++ {
		webIOs = append(webIOs, io.NewWebIO(key))
	}
	ios = append(ios, webIOs)

	var fileIOs []rsrc.IO
	for i := 0; i < 10; i++ {
		fileIOs = append(fileIOs, io.FileIO{})
	}
	ios = append(ios, fileIOs)

	if isCacheRunning(14003) {
		var cacheIOs []rsrc.IO
		for i := 0; i < 2; i++ {
			cacheIOs = append(cacheIOs, &io.CacheIO{Port: 14003})
		}
		ios = append(ios, cacheIOs)
	}

	st, err := store.NewCache(ios)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func main() {
	s, err := createStore()

	if err != nil {
		fmt.Println(err)
		return
	}

	session, _ := unpack.LoadSessionInfo(s)
	d := display.NewTerminal()

	err = command.Execute(os.Args, session, s, d)
	if err != nil {
		d.Display(&format.Error{Err: err})
	}
}
