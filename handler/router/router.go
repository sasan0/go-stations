package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	mux.Handle("/healthz", middleware.Recovery(middleware.Auth(middleware.IdentifyDevice(middleware.Logger(handler.NewHealthzHandler())))))
	mux.Handle("/todos", middleware.Recovery(middleware.Auth(middleware.IdentifyDevice(middleware.Logger(handler.NewTODOHandler(service.NewTODOService(todoDB)))))))
	mux.Handle("/do-panic", middleware.Recovery(middleware.Auth(middleware.IdentifyDevice(middleware.Logger(handler.NewPanicHandler())))))
	return mux
}
