package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"upspin.io/client"
	"upspin.io/config"
	_ "upspin.io/transports"
	"upspin.io/upspin"
)

func main() {
	cfg, err := config.FromFile("config")
	if err != nil {
		panic(err)
	}

	cl := client.New(cfg)

	// e.g. GET /augie@upspin.io/Images/camstream.mp4
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		f, err := cl.Open(upspin.PathName(path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		_, err = io.Copy(w, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
