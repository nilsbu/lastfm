package main

import (
	"fmt"
	"time"

	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

func createStore() (io.Store, error) {
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

	st, err := io.NewStore(
		[][]rsrc.IO{webIOs, fileIOs},
	)
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

	session, err := unpack.LoadSessionInfo(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		organize.BackupUpdateHistory(session.User, 30, s)

		time.Sleep(60 * 60 * time.Second)
	}
}
