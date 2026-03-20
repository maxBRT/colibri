package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"www.github.com/maxbrt/colibri/internal/database"
)

type Server struct {
	Router *chi.Mux
	Db     *database.Queries
}

func NewServer(db *database.Queries) *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	s.Db = db
	return s
}

func (s *Server) MountHandlers() {
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

	s.Router.Route("/v1", func(r chi.Router) {
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		})

		r.Get("/categories", func(w http.ResponseWriter, r *http.Request) {
			categories, err := s.Db.ListCategories(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			c, err := json.Marshal(categories)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(200)
			w.Write(c)
		})

		r.Get("/sources", func(w http.ResponseWriter, r *http.Request) {
			categories := r.URL.Query()["category"]
			var sources []database.Source

			if len(categories) > 0 {
				for i, c := range categories {
					categories[i] = strings.ToLower(c)
				}
				s, err := s.Db.ListSourcesByCategory(r.Context(), categories)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				sources = s
			} else {
				s, err := s.Db.ListSources(r.Context())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				sources = s
			}

			w.Header().Set("Content-Type", "application/json")
			c, err := json.Marshal(sources)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(200)
			w.Write(c)
		})

		r.Get("/posts", func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()

			sourceIDs := q["sources"]
			var posts []database.Post

			if len(sourceIDs) > 0 {
				for i, s := range sourceIDs {
					sourceIDs[i] = strings.ToLower(s)
				}

				p, err := s.Db.ListPostsForSource(r.Context(), sourceIDs)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				posts = p
			} else {
				p, err := s.Db.ListPosts(r.Context())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				posts = p
			}

			w.Header().Set("Content-Type", "application/json")
			c, err := json.Marshal(posts)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(200)
			w.Write(c)
		})

	})
}
