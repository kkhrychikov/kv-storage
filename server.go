package kvstorage

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	httpServer *http.Server
	storage    Storage
}

// NewServer creates new HTTP server with user storage
// Use rwTimeoutSec to specify read and write timeout
func NewServer(port int, rwTimeoutSec int, storage Storage) (*Server, error) {
	if port < 0 || 65535 < port {
		return nil, fmt.Errorf("Port %v is not valid", port)
	}
	addr := ":" + strconv.Itoa(port)
	server := Server{storage: storage}
	router := &router{&server}
	server.httpServer = &http.Server{
		Addr:         addr,
		Handler:      router,
		WriteTimeout: time.Second * time.Duration(rwTimeoutSec),
		ReadTimeout:  time.Second * time.Duration(rwTimeoutSec)}
	return &server, nil
}

// Start launches server in separate goroutine
func (server *Server) Start() {
	go func() {
		err := server.httpServer.ListenAndServe()
		log.Printf("Error start server on %s : %v", server.httpServer.Addr, err)
	}()
}

// Stop gracefully shuts down the server without interrupting any active connections
func (server *Server) Stop() {
	go func() {
		err := server.httpServer.ListenAndServe()
		log.Printf("Error start server on %s : %v", server.httpServer.Addr, err)
	}()
}

// ForceStop immediately closes all active connections.
// For a graceful shutdown, use Stop.
func (server *Server) ForceStop() {
	go func() {
		err := server.httpServer.ListenAndServe()
		log.Printf("Error start server on %s : %v", server.httpServer.Addr, err)
	}()
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
	r.ParseForm()
	key := r.Form.Get("key")
	val := r.Form.Get("val")
	var res string
	var err error
	switch r.URL.Path {
	case "/insert":
		err = router.server.storage.Insert(key, val)
	case "/update":
		err = router.server.storage.Update(key, val)
	case "/select":
		res, err = router.server.storage.Select(key)
	case "/delete":
		err = router.server.storage.Delete(key)
	default:
		err = errors.New(`Unknown command path. Use "insert","update", "select" or "delete"`)
	}
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	if res != "" {
		w.Write([]byte(res))
	}
}
