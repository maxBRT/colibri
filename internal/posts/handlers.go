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
	var posts []*Post

	if len(sourceIDs) > 0 {
		for i, s := range sourceIDs {
			sourceIDs[i] = strings.ToLower(s)
		}

		dbPosts, err := h.db.ListPostsForSource(r.Context(), sourceIDs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, p := range dbPosts {
			posts = append(posts, PostFromModel(p))
		}
	} else {
		dbPosts, err := h.db.ListPosts(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, p := range dbPosts {
			posts = append(posts, PostFromModel(p))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(posts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}
