package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/udodinho/bookings/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// pass; do nothing
	default:
		t.Errorf("routes() returned %T, want *chi.Mux", v)
	}
}
