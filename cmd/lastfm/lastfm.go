package main

import (
	"fmt"
	"os"

	"github.com/nilsbu/lastfm/pkg/command"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
)

func createStore() (store.Store, error) {
	key, err := organize.LoadAPIKey(io.FileIO{})
	if err != nil {
		return nil, err
	}

	var webIOs []rsrc.IO
	for i := 0; i < 32; i++ {
		webIOs = append(webIOs, io.NewWebIO(key))
	}

	var fileIOs []rsrc.IO
	for i := 0; i < 10; i++ {
		fileIOs = append(fileIOs, io.FileIO{})
	}

	st, err := store.NewCache([][]rsrc.IO{webIOs, fileIOs})
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

	sid, err := organize.LoadSessionID(s)
	if err != nil {
		sid = ""
	}

	d := display.NewTerminal()

	err = command.Execute(os.Args, sid, s, d)
	if err != nil {
		d.Display(&format.Error{Err: err})
	}
}
