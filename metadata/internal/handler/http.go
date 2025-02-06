package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"bikraj.movie_microservice.net/metadata/internal/controller/metadata"
)

type Handler struct {
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := req.Context()

	m, err := h.ctrl.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, metadata.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			log.Printf("Repository get error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Reponse encoder error: %v", err)
	}
}
