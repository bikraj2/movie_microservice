package grpc

import (
	"context"
	"errors"

	"bikraj.movie_microservice.net/gen"
	"bikraj.movie_microservice.net/metadata/internal/controller/metadata"
	model "bikraj.movie_microservice.net/metadata/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedMetadataServiceServer
	svc *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{svc: ctrl}
}

func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataReqeust) (*gen.GetMetadataResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty movie_id")
	}

	m, err := h.svc.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, metadata.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetMetadataResponse{Metadata: model.MetadataToProto(m)}, nil
}

func (h *Handler) PutMetadata(ctx context.Context, req *gen.PutMetadataRequest) (*gen.PutMetadataResponse, error) {
	if req == nil || req.Metadata == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or metadata")
	}
	if err := h.svc.Put(ctx, model.MetadataFromProto(req.Metadata)); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.PutMetadataResponse{}, nil
}
