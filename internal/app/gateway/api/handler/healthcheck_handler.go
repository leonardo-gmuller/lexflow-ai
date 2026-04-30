package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) healthcheckSetup(router chi.Router) {
	router.Get("/healthcheck", func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		response := map[string]string{"status": "ok"}

		if err := json.NewEncoder(rw).Encode(response); err != nil {
			http.Error(rw, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
