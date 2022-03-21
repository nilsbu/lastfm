package main

import (
	"fmt"
	"os"

	"github.com/nilsbu/lastfm/pkg/command"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

func createStore() (io.Store, error) {
	key, err := unpack.LoadAPIKey(io.FileIO{})
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

	st, err := io.NewStore([][]rsrc.IO{webIOs, fileIOs})
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
	d := display.NewCSV("total.csv", ",") // TODO file name as param

	err = command.Execute(os.Args, session, s, d)
	if err != nil {
		fmt.Println(err)
	}
}
