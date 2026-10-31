package main

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer, "Akshaj")

	got := buffer.String()
	want := "Hello, Akshaj"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
