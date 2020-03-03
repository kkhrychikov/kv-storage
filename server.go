package kvstorage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	httpServer *http.Server
	storage    Storage
}

func NewServer(port int, storage Storage) (*Server, error) {
	if port < 0 || 65535 < port {
		return nil, fmt.Errorf("Port %v is not valid", port)
	}
	addr := ":" + strconv.Itoa(port)
	server := Server{storage: storage}
	router := &router{&server}
	server.httpServer = &http.Server{
		Addr:         addr,
		Handler:      router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60}
	return &server, nil
}

func (server *Server) Start() {
	err := server.httpServer.ListenAndServe()
	log.Printf("Error start server on %s : %v", server.httpServer.Addr, err)
}

type router struct {
	server *Server
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(400)
		w.Write([]byte(`This is only storage server. It doesn't support GET requests`))
		return
	}
	router.ProcessStorageRequest(w, r)
	return
}

func (router *router) ProcessStorageRequest(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	kv, err := router.server.storage.Unmarshal(data)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	var res interface{}
	switch r.URL.Path {
	case "/insert":
		res, err = router.server.storage.Insert(kv)
	case "/update":
		res, err = router.server.storage.Update(kv)
	case "/select":
		res, err = router.server.storage.Select(kv)
	case "/delete":
		res, err = router.server.storage.Delete(kv)
	default:
		err = errors.New(`Unknown command path. Use "insert","update", "select" or "delete"`)
	}
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	var b []byte
	if res != nil {
		b, err = router.server.storage.Marshal(res)
	}
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	if b != nil {
		w.Write(b)
	}
}
