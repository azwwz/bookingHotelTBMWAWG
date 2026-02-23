package render

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
)

var sessionManager *scs.SessionManager

func TestMain(m *testing.M) {
	setup()

	os.Exit(m.Run())
}

func setup() {
	app = &config.AppConfig{}

	app.InProduction = false

	// config appConfig log.logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog		

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction

	app.SessionManager = sessionManager

	NewTemplates(app)

}

type myWriter struct{}

func (m *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (m *myWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
func (m *myWriter) WriteHeader(b int) {

}
