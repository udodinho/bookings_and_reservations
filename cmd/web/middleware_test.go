package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myHandle myHandler

	h := NoSurf(&myHandle)

	switch v := h.(type) {
	case http.Handler:
	// pass
	default:
		t.Errorf("NoSurf returned %T, want http.Handler", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myHandle myHandler

	sessions := SessionLoad(&myHandle)

	switch v := sessions.(type) {
	case http.Handler:
	// pass
	default:
		t.Errorf("NoSurf returned %T, want http.Handler", v)
	}
}
