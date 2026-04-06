package sources

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
	categories := r.URL.Query()["category"]
	var sources []*Source

	if len(categories) > 0 {
		for i, c := range categories {
			categories[i] = strings.ToLower(c)
		}
		rows, err := h.db.ListSourcesByCategoryWithLogo(r.Context(), categories)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, row := range rows {
			sources = append(sources, SourceFromWithLogoRow(row.ID, row.Name, row.Url, row.Category, row.LogoUrl))
		}
	} else {
		rows, err := h.db.ListSourcesWithLogo(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, row := range rows {
			sources = append(sources, SourceFromWithLogoRow(row.ID, row.Name, row.Url, row.Category, row.LogoUrl))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	c, err := json.Marshal(sources)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	_, err = w.Write(c)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}
