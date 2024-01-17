package main

import (
	"testing"
	"time"

	"github.com/munanadi/pokedex/pokecache"
)

func TestAddGet(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Second)

	cache.Add("foo", []byte("hi"))

	fooValue, _ := cache.Get("foo")

	if string(fooValue) != "hi" {
		t.Errorf("expected foo to have 'hi' but got %s instead", fooValue)
	}
}
