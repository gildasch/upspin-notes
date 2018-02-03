package main

import (
	"encoding/json"
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
	cfg, err := config.FromFile("config.fosdem")
	if err != nil {
		panic(err)
	}

	cl := client.New(cfg)

	http.HandleFunc("/list/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/list/")
		des, err := cl.Glob(path + "/*")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(des)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

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

// func otherHandlers() {
// 		http.HandleFunc("/inspect/", func(w http.ResponseWriter, r *http.Request) {
// 		path := strings.TrimPrefix(r.URL.Path, "/inspect/")
// 		de, err := cl.Lookup(upspin.PathName(path), false)
// 		if err != nil {
// 			if uerrors.Is(uerrors.NotExist, err) {
// 				fmt.Fprintf(w, "%q does not exist", path)
// 				return
// 			}
// 			json.NewEncoder(w).Encode(err)
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 			return
// 		}
// 		json.NewEncoder(w).Encode(de)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	})

// 	http.HandleFunc("/list/", func(w http.ResponseWriter, r *http.Request) {
// 		path := strings.TrimPrefix(r.URL.Path, "/list/")
// 		des, err := cl.Glob(path + "/*")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 			return
// 		}

// 		json.NewEncoder(w).Encode(des)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	})

// }
