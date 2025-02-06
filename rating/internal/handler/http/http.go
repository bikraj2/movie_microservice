package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"bikraj.movie_microservice.net/rating/internal/controller/rating"
	"bikraj.movie_microservice.net/rating/pkg/model"
)

// Handler defines a rating service controller
type Handler struct {
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	recordID := model.RecordID(req.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := model.RecordType(req.FormValue("type"))

	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		v, err := h.ctrl.GetAgrregatedRaring(req.Context(), recordID, recordType)
		if err != nil && errors.Is(err, rating.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("response error code: %v", err)
		}
	case http.MethodPut:
		userID := model.UserID(req.FormValue("userId"))
		v, err := strconv.ParseFloat(req.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := h.ctrl.PutRating(req.Context(), recordID, recordType, &model.Rating{RatingValue: model.RatingValue(v), UserID: userID}); err != nil {
			log.Printf("Repository put erro: %^v ", err)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
