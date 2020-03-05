package kvstorage

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	target     string
	timeoutSec int
}

// GetServerAddr returns server hostname
func (client *Client) GetServerAddr() string {
	return client.target
}

// NewClient gets server addres IP:PORT and request timeout in seconds
func NewClient(serverAddr string, timeoutSec int) *Client {
	return &Client{target: serverAddr, timeoutSec: timeoutSec}
}

// Insert sends key value to storage
// If key already in storage error will be returned
func (client *Client) Insert(key, val string) error {
	u := url.URL{
		Scheme: "http",
		Host:   client.target,
		Path:   "insert"}
	form := url.Values{}
	form.Add("key", key)
	form.Add("val", val)
	httpClient := http.Client{Timeout: time.Duration(client.timeoutSec) * time.Second}
	resp, err := httpClient.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
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

// Update sends key value to storage
// If key not in storage error will be returned
func (client *Client) Update(key, val string) error {
	u := url.URL{
		Scheme: "http",
		Host:   client.target,
		Path:   "update"}
	form := url.Values{}
	form.Add("key", key)
	form.Add("val", val)
	httpClient := http.Client{Timeout: time.Duration(client.timeoutSec) * time.Second}
	resp, err := httpClient.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
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

// Delete sends key value to storage
// If key not in storage error will be returned
func (client *Client) Delete(key string) error {
	u := url.URL{
		Scheme: "http",
		Host:   client.target,
		Path:   "delete"}
	form := url.Values{}
	form.Add("key", key)
	httpClient := http.Client{Timeout: time.Duration(client.timeoutSec) * time.Second}
	resp, err := httpClient.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
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

// Select requests value from storage
// If key not in storage error will be returned
func (client *Client) Select(key string) (string, error) {
	u := url.URL{
		Scheme: "http",
		Host:   client.target,
		Path:   "select"}
	form := url.Values{}
	form.Add("key", key)
	httpClient := http.Client{Timeout: time.Duration(client.timeoutSec) * time.Second}
	resp, err := httpClient.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Select failed: Can't get message from server " + client.target)
	}
	if resp.StatusCode != 200 {
		return "", errors.New("Select failed: " + string(b))
	}
	return string(b), nil
}
