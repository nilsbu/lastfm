package main

import (
	"fmt"
	"time"

	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

func dumpChan() chan<- format.Formatter {
	obChan := make(chan format.Formatter)
	go func() {
		for range obChan {
		}
	}()

	return obChan
}

func createStore(observer chan format.Formatter) (io.Store, error) {
	key, err := unpack.LoadAPIKey(io.FileIO{})
	if err != nil {
		return nil, err
	}

	var webIOs []rsrc.IO
	for i := 0; i < 1; i++ {
		webIOs = append(webIOs, io.NewWebIO(key))
	}

	var fileIOs []rsrc.IO
	for i := 0; i < 10; i++ {
		fileIOs = append(fileIOs, io.FileIO{})
	}

	st, err := io.NewObservedStore(
		[][]rsrc.IO{webIOs, fileIOs},
		[]chan<- format.Formatter{observer, dumpChan()},
	)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func main() {
	webObserver := make(chan format.Formatter)
	// d := display.NewTimedTerminal(webObserver, 1*time.Second)

	s, err := createStore(webObserver)
	if err != nil {
		fmt.Println(err)
		return
	}

	session, err := unpack.LoadSessionInfo(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		if err := organize.BackupUpdateHistory(session.User, 30, s); err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(60 * 60 * time.Second)
	}
}
