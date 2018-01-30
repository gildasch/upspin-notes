package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"upspin.io/client"
	"upspin.io/config"
	"upspin.io/rpc/dirserver"
	"upspin.io/rpc/storeserver"
	_ "upspin.io/transports"
	"upspin.io/upspin"
)

type Accesser interface {
	Get(path string) (io.Reader, error)
}

type TranscoderServer interface {
	DirServer() upspin.DirServer
	StoreServer(accesser Accesser) upspin.StoreServer
}

func ServeUpspin(cfg upspin.Config, accesser Accesser, server TranscoderServer) {
	http.Handle("/api/Store/", storeserver.New(cfg, server.StoreServer(accesser), ""))
	http.Handle("/api/Dir/", dirserver.New(cfg, server.DirServer(), ""))

	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", nil)
}

type UpspinAccesser struct {
	upspin.Client
}

func NewUpspinAccesser(cfg upspin.Config) *UpspinAccesser {
	return &UpspinAccesser{client.New(cfg)}
}

func (ua *UpspinAccesser) Get(path string) (io.Reader, error) {
	return ua.Open(upspin.PathName(path))
}

// Server provides DirServer and StoreServer implementations
// that serve transcoded videos
type Server struct {
	// Set by New.
	server upspin.Config
	// Set by Dial.
	user upspin.Config
}

func (s *Server) DirServer() upspin.DirServer {
	return dirServer{s}
}

type dirServer struct {
	*Server
}

func (s dirServer) Dial(cfg upspin.Config, e upspin.Endpoint) (upspin.Service, error) {
	dialed := *s.Server
	dialed.user = cfg
	return dirServer{&dialed}, nil
}

func (s dirServer) Lookup(name upspin.PathName) (*upspin.DirEntry, error) {
	return &upspin.DirEntry{
		Name:    name,
		Packing: upspin.PlainPack,
		Time:    upspin.Now(),
		Blocks: []upspin.DirBlock{
			upspin.DirBlock{
				Location: upspin.Location{
					Endpoint: upspin.Endpoint{
						Transport: upspin.Remote,
						NetAddr:   ":8080",
					},
					Reference: upspin.Reference(name)},
				Offset: 0, // how to find it?
				Size:   0, // how to find it?
			},
		},
		Writer: s.server.UserName(), // TODO: Is there a better answer?
	}, nil
}

func (s *Server) StoreServer(accesser Accesser) upspin.StoreServer {
	return storeServer{s, accesser}
}

type storeServer struct {
	*Server
	accesser Accesser
}

func (s storeServer) Dial(cfg upspin.Config, e upspin.Endpoint) (upspin.Service, error) {
	dialed := *s.Server
	dialed.user = cfg
	return storeServer{&dialed, s.accesser}, nil
}

func (s storeServer) Get(ref upspin.Reference) ([]byte, *upspin.Refdata, []upspin.Location, error) {
	r, err := s.accesser.Get(string(ref))
	if err != nil {
		return nil, nil, nil, err
	}
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, nil, err
	}
	return bytes,
		&upspin.Refdata{
			Reference: ref,
			Volatile:  true, // or Duration
		},
		nil, nil
}

func main() {
	cfg, err := config.FromFile("config")
	if err != nil {
		panic(err)
	}

	accesser := NewUpspinAccesser(cfg)
	server := &Server{}

	ServeUpspin(cfg, accesser, server)
}
