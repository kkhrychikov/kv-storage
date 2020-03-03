package kvstorage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	target string
}

func (client *Client) GetServerAddr() string {
	return client.target
}

func NewClient(serverAddr string) *Client {
	return &Client{target: serverAddr}
}

func prepareKV(key string, val interface{}) (io.Reader, error) {
	kv := keyVal{key, val}
	data, err := json.Marshal(kv)
	if err != nil {
		return nil, fmt.Errorf("Error on preparing data: %v", err)
	}
	return bytes.NewReader(data), nil
}

func (client *Client) Insert(key string, val interface{}) error {
	r, err := prepareKV(key, val)
	if err != nil {
		return err
	}
	u := url.URL{
		Scheme: "http",
		Host:   client.target,
		Path:   "insert"}
	resp, err := http.Post(u.String(), "application/json", r)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Insert failed: Can't get message from server " + client.target)
		}
		return errors.New("Insert failed: " + string(b))
	}
	return nil
}

func (client *Client) Update(key string, val interface{}) error {
	r, err := prepareKV(key, val)
	if err != nil {
		return err
	}
	u := url.URL{
		Scheme: "http",
		Host:   client.target,
		Path:   "update"}
	resp, err := http.Post(u.String(), "application/json", r)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Update failed: Can't get message from server " + client.target)
		}
		return errors.New("Update failed: " + string(b))
	}
	return nil
}

func (client *Client) Delete(key string) error {
	r, err := prepareKV(key, nil)
	if err != nil {
		return err
	}
	u := url.URL{
		Scheme: "http",
		Host:   client.target,
		Path:   "delete"}
	resp, err := http.Post(u.String(), "application/json", r)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Delete failed: Can't get message from server " + client.target)
		}
		return errors.New("Delete failed: " + string(b))
	}
	return nil
}

func (client *Client) Select(key string) ([]byte, error) {
	r, err := prepareKV(key, nil)
	if err != nil {
		return nil, err
	}
	u := url.URL{
		Scheme: "http",
		Host:   client.target,
		Path:   "select"}
	resp, err := http.Post(u.String(), "application/json", r)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Select failed: Can't get message from server " + client.target)
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Select failed: " + string(b))
	}
	return b, nil
}
