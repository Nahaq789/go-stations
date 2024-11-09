package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	todoService := service.NewTODOService(
		todoDB,
	)
	// register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.NewHealthzHandler().ServeHTTP)
	mux.HandleFunc("/todos", handler.NewTODOHandler(todoService).ServeHTTP)
	return mux
}
