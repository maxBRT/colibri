package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"www.github.com/maxbrt/colibri/internal/categories"
	"www.github.com/maxbrt/colibri/internal/database"
	"www.github.com/maxbrt/colibri/internal/posts"
	"www.github.com/maxbrt/colibri/internal/sources"
)

type Server struct {
	Router *chi.Mux
}

func NewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers(db *database.Queries) {
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
		MaxAge:         300,
	}))

	s.Router.Use(middleware.Heartbeat("/health"))
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.AllowContentType("application/json"))
	s.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			next.ServeHTTP(w, r)
		})
	})
	s.Router.Use(httprate.LimitByIP(100, time.Minute))

	postsHandler := posts.NewHandler(db)
	sourcesHandler := sources.NewHandler(db)
	categoriesHandler := categories.NewHandler(db)

	s.Router.Route("/v1", func(r chi.Router) {
		r.Get("/categories", categoriesHandler.List)
		r.Get("/sources", sourcesHandler.List)
		r.Get("/posts", postsHandler.List)
	})
}
