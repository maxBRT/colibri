package posts

import (
	"encoding/json"
	"net/http"
	"strings"

	"www.github.com/maxbrt/colibri/internal/database"
)

type Handler struct {
	db *database.Queries
}

func NewHandler(db *database.Queries) *Handler {
	return &Handler{db: db}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	sourceIDs := q["sources"]
	var posts []database.Post

	if len(sourceIDs) > 0 {
		for i, s := range sourceIDs {
			sourceIDs[i] = strings.ToLower(s)
		}

		p, err := h.db.ListPostsForSource(r.Context(), sourceIDs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		posts = p
	} else {
		p, err := h.db.ListPosts(r.Context())
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
}
