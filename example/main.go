package main

import (
	"fmt"
	kvstorage "kv-storage"
	"log"
	"runtime"
)

func main() {
	storage := kvstorage.CreateStringsStorage()
	server, err := kvstorage.NewServer(1234, storage)
	if err != nil {
		log.Fatalln(err)
	}
	go server.Start()
	runtime.Gosched()
	client := kvstorage.NewClient("127.0.0.1:1234")
	err = client.Insert("first", "kirill")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success insert")
	}
	d, err := client.Select("first")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(d))
	}
	d, err = client.Select("second")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(d))
	}
	err = client.Update("first", "Alex")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(d))
	}
	d, err = client.Select("first")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(d))
	}
}
