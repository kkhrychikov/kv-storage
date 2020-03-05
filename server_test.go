package kvstorage

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestServer(t *testing.T) {
	_, err := NewServer(-1, -1, NewStorage())
	if err == nil {
		t.Error("expected error if port is not valid")
	}

	port := 1234
	server, err := NewServer(port, 10, NewStorage())
	if err != nil {
		t.Error(err)
	}

	server.Start()
	runtime.Gosched()

	addr := "127.0.0.1:" + strconv.Itoa(port)
	client := NewClient(addr, 10)
	err = client.Insert("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	res, err := client.Select("foo")
	if err != nil {
		t.Error(err)
	}
	if res != "bar" {
		t.Error("expected value to be bar, not ", res)
	}
	err = client.Update("foo", "foobar")
	if err != nil {
		t.Error(err)
	}
	res, err = client.Select("foo")
	if err != nil {
		t.Error(err)
	}
	if res != "foobar" {
		t.Error("expected value to be foobar, not ", res)
	}
	err = client.Delete("foo")
	if err != nil {
		t.Error(err)
	}
	_, err = client.Select("foo")
	if err == nil {
		t.Error("expected to be err with no key")
	}
	err = server.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestServerLoad(t *testing.T) {
	port := 1234
	storage := NewStorage()
	server, err := NewServer(port, 10, storage)
	if err != nil {
		t.Error(err)
	}

	server.Start()
	runtime.Gosched()

	addr := "127.0.0.1:" + strconv.Itoa(port)
	client := NewClient(addr, 10)
	wg := sync.WaitGroup{}
	pairsLimit := 1000
	wg.Add(pairsLimit)
	for i := 0; i < pairsLimit; i++ {
		num := strconv.Itoa(i)
		go func() {
			err = client.Insert("foo"+num, "bar"+num)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	server.Stop()
	els, err := storage.Dump()
	pairs := strings.Split(string(els), ",")
	if len(pairs) != pairsLimit {
		t.Errorf("storage contains not %d pairs, only %d", pairsLimit, len(pairs))
	}
	storage.Reset()
	msg, err := storage.Load(els)
	if err != nil || msg != fmt.Sprintf("Added %d/%d pairs", pairsLimit, pairsLimit) {
		t.Errorf("cant load prevous dump")
	}
}
