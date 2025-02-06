package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	movie "bikraj.movie_microservice.net/movie/internal/controller"
	"bikraj.movie_microservice.net/movie/internal/gateway"
)

type Handler struct {
	ctrl *movie.Controller
}

func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	details, err := h.ctrl.Get(r.Context(), id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Repository get error: %v\n", err)
		return
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
