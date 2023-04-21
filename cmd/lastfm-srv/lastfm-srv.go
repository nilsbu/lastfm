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

	// fmt.Println(r.URL.Path)
	fmt.Fprintln(os.Stderr, r.URL.Path)

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

	pl := pipeline.New(session, s)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/charts/", http.StripPrefix("/charts/", fs))

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		handleRequest(session, s, pl, rw, r)
	})

	tlsCertPath := os.Getenv("TLS_CERT_PATH")
	certPath := tlsCertPath + "/cert.pem"
	keyPath := tlsCertPath + "/privkey.pem"

	if err := http.ListenAndServeTLS(":3000", certPath, keyPath, nil); err != nil {
		log.Fatal(err)
	}
}
