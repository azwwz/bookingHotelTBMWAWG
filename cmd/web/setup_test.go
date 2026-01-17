package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
)

func TestMain(m *testing.M) {
	app = &config.AppConfig{InProduction: false}
	// sessionManager = scs.New()
	// app.SessionManager = sessionManager
	os.Exit(m.Run())
}

type MyHandler struct{}

func (m *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
