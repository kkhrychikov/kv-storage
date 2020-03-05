package kvstorage

import "testing"

func TestStorageCommands(t *testing.T) {
	storage := NewStorage()

	err := storage.Insert("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	err = storage.Insert("foo", "foobar")
	if err == nil {
		t.Error("expected to be err with key already exists")
	}
	err = storage.Insert("", "foobar")
	if err == nil {
		t.Error("expected to be err with empty key")
	}

	res, err := storage.Select("foo")
	if err != nil {
		t.Error(err)
	}
	if res != "bar" {
		t.Error("expected value to be bar, not ", res)
	}
	_, err = storage.Select("")
	if err == nil {
		t.Error("expected to be err with empty key")
	}

	err = storage.Update("foo", "foobar")
	if err != nil {
		t.Error(err)
	}
	err = storage.Update("bar", "foobar")
	if err == nil {
		t.Error("expected to be err with key not exists")
	}
	err = storage.Update("", "foobar")
	if err == nil {
		t.Error("expected to be err with empty key")
	}

	res, err = storage.Select("foo")
	if err != nil {
		t.Error(err)
	}
	if res != "foobar" {
		t.Error("expected value to be foobar, not ", res)
	}

	err = storage.Delete("foo")
	if err != nil {
		t.Error(err)
	}
	err = storage.Delete("bar")
	if err == nil {
		t.Error("expected to be err with key not exists")
	}
	err = storage.Delete("")
	if err == nil {
		t.Error("expected to be err with empty key")
	}

	_, err = storage.Select("foo")
	if err == nil {
		t.Error("expected to be err with no key")
	}
}
