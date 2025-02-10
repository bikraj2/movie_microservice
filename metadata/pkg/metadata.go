package model

import "bikraj.movie_microservice.net/gen"

type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
}

func MetadataToProto(m *Metadata) *gen.Metadata {

	return &gen.Metadata{
		Id:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}

func MetadataFromProto(m *gen.Metadata) *Metadata {

	return &Metadata{
		ID:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}
