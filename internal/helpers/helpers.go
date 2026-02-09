package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
)

var app *config.AppConfig

func NewHelper(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("client error with error code ", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
