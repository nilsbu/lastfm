package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/nilsbu/lastfm/pkg/command"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
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

	st, err := io.NewObserved(
		[][]rsrc.IO{webIOs, fileIOs},
		[]chan<- format.Formatter{webObserver, dumpChan()},
	)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func handleRequest(
	session *unpack.SessionInfo,
	s io.Store,
	w http.ResponseWriter,
	r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// TODO: this needs to work better
	if r.URL.Path == "/favicon.ico" {
		return
	}

	args := []string{"lastfm-srv"}
	args = append(args, strings.Split(r.URL.Path, "/")[1:]...)

	for k, vs := range r.URL.Query() {
		args = append(args, fmt.Sprintf("-%v=%v", k, vs[0]))
	}

	err := command.Execute(args, session, s, display.NewWeb(w))
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	s, err := createStore(dumpChan())

	if err != nil {
		fmt.Println(err)
		return
	}

	session, err := unpack.LoadSessionInfo(s)

	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO Reuse early stages of the charts
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		handleRequest(session, s, rw, r)
	})

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
