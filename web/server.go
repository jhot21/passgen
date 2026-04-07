package web

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewServer returns the root HTTP handler with API routes and static file serving.
func NewServer() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/generate", handleGenerate)
	r.Post("/api/generate", handleGenerate)

	r.Handle("/*", staticHandler())

	return r
}

func staticHandler() http.Handler {
	sub, err := fs.Sub(staticFiles, "app/dist")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(sub))
}
