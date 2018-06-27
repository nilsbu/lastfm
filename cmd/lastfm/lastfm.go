package main

import (
	"fmt"
	"os"

	"github.com/nilsbu/lastfm/pkg/command"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/store"
)

func createStore() (store.Store, error) {
	key, err := organize.LoadAPIKey(io.FileIO{})
	if err != nil {
		return nil, err
	}

	var downloaders []io.Reader
	for i := 0; i < 16; i++ {
		downloaders = append(downloaders, io.Downloader(key))
	}

	var fileReaders []io.Reader
	for i := 0; i < 10; i++ {
		fileReaders = append(fileReaders, io.FileIO{})
	}

	var fileWriters []io.Writer
	for i := 0; i < 10; i++ {
		fileWriters = append(fileWriters, io.FileIO{})
	}

	st, err := store.New(
		[][]io.Reader{downloaders, fileReaders},
		[][]io.Writer{[]io.Writer{io.FailIO{}}, fileWriters})
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

	command.Execute(os.Args, sid, s)
}
