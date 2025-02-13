package grpc

import (
	"context"
	"errors"

	"bikraj.movie_microservice.net/gen"
	"bikraj.movie_microservice.net/rating/internal/controller/rating"
	"bikraj.movie_microservice.net/rating/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedRatingServiceServer
	svc *rating.Controller
}

func New(svc *rating.Controller) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.RecordId == "" {
		return nil, status.Error(codes.InvalidArgument, "record ID cannot be empty")
	}

	if req.RecordType == "" {
		return nil, status.Error(codes.InvalidArgument, "record type cannot be empty")
	}
	v, err := h.svc.GetAgrregatedRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil && errors.Is(err, rating.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetAggregatedRatingResponse{RatingValue: v}, nil
}

func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty record id , record type or UserId")
	}
	err := h.svc.PutRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType), &model.Rating{UserID: model.UserID(req.UserId), RatingValue: model.RatingValue(req.RatingValue)})
	if err != nil {
		return nil, err
	}
	return &gen.PutRatingResponse{}, nil

}
