package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nilsbu/lastfm/pkg/command"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
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

func createStore(webObserver chan<- format.Formatter) (io.Store, error) {
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
		[]chan<- format.Formatter{webObserver, dumpChan()},
	)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func main() {

	webObserver := make(chan format.Formatter)
	d := display.NewTimedTerminal(webObserver, 1*time.Second)

	s, err := createStore(webObserver)

	if err != nil {
		fmt.Println(err)
		return
	}

	session, _ := unpack.LoadSessionInfo(s)
	pl := pipeline.New(session, s)

	err = command.Execute(os.Args, session, s, pl, d)
	if err != nil {
		fmt.Println(err)
	}
}
