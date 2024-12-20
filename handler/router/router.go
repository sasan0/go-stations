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
	mux.Handle("/healthz", middleware.IdentifyDevice(middleware.Recovery(handler.NewHealthzHandler())))
	mux.Handle("/todos", middleware.IdentifyDevice(middleware.Recovery(handler.NewTODOHandler(service.NewTODOService(todoDB)))))
	mux.Handle("/do-panic", middleware.IdentifyDevice(middleware.Recovery(handler.NewPanicHandler())))
	return mux
}
