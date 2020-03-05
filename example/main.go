package main

import (
	"fmt"
	"kvstorage"
	"log"
	"runtime"
	"strconv"
)

func main() {
	port := 1234
	server, err := kvstorage.NewServer(1234, 10, kvstorage.NewStorage())
	if err != nil {
		log.Fatalln(err)
	}
	server.Start()

	runtime.Gosched()
	client := kvstorage.NewClient("127.0.0.1:"+strconv.Itoa(port), 10)
	err = client.Insert("foo", "bar")
	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Select("foo")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
