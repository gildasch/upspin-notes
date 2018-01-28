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

type Accesser interface {
	Get(path string) (io.Reader, error)
}

type Transcoder interface {
	Transcode(input io.Reader, output io.Writer) error
}

func ServeHTTP(accesser Accesser, transcoder Transcoder) {
	// e.g. GET /augie@upspin.io/Images/camstream.mp4
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		video, err := accesser.Get(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = transcoder.Transcode(video, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", nil)
}

type UpspinAccesser struct {
	upspin.Client
}

func NewUpspinAccesser() *UpspinAccesser {
	cfg, err := config.FromFile("config")
	if err != nil {
		panic(err)
	}

	return &UpspinAccesser{client.New(cfg)}
}

func (ua *UpspinAccesser) Get(path string) (io.Reader, error) {
	return ua.Open(upspin.PathName(path))
}

type FFMPEGTranscoder struct{}

func (ft *FFMPEGTranscoder) Transcode(input io.Reader, output io.Writer) error {
	// nothing is done here, TODO: implement real transcoder
	_, err := io.Copy(output, input)
	return err
}

func main() {
	accesser := NewUpspinAccesser()
	transcoder := &FFMPEGTranscoder{}

	ServeHTTP(accesser, transcoder)
}
