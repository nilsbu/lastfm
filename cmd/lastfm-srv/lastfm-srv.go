package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nilsbu/lastfm/pkg/command"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/refresh"
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

func handleRequest(
	session *unpack.SessionInfo,
	s io.Store,
	pl pipeline.Pipeline,
	w http.ResponseWriter,
	r *http.Request) {

	// Added this to allow CORS for development purposes
	w.Header().Set("Access-Control-Allow-Origin", "http://nboon.de:3000")

	fmt.Println("Request:", r.URL.Path)

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// TODO: this needs to work better
	if r.URL.Path == "/favicon.ico" {
		return
	}

	a := strings.Split(r.URL.Path, "/")[1:]
	var d display.Display
	if a[0] == "json" {
		d = display.NewJSON(w)
		a = a[1:]
	} else {
		d = display.NewWeb(w)
	}

	args := []string{"lastfm-srv"}
	args = append(args, a...)

	for k, vs := range r.URL.Query() {
		args = append(args, fmt.Sprintf("-%v=%v", k, vs[0]))
	}

	err := command.Execute(args, session, s, pl, d)
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	fmt.Println("Starting server...")

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Listening on port", port)

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

	pl, trigger := refresh.WrapRefresh(s, pipeline.New(session, s), session)
	go refresh.PeriodicRefresh(trigger, 0, 1, 0) // Every night at 01:00 (UTC)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		handleRequest(session, s, pl, rw, r)
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
