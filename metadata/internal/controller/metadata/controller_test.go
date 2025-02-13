package metadata

import (
	"context"
	"errors"
	"testing"

	gen "bikraj.movie_microservice.net/gen/mock/metadata/repository"
	"bikraj.movie_microservice.net/metadata/internal/repository"
	"bikraj.movie_microservice.net/metadata/pkg"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestController(t *testing.T) {
	tests := []struct {
		name         string
		expRepoRes   *model.Metadata
		expRepoError error
		wantRes      *model.Metadata
		wantErr      error
	}{
		{
			name:         "not found",
			expRepoError: repository.ErrNotFound,
			wantErr:      repository.ErrNotFound,
		},
		{
			name:         "unexpected error",
			expRepoError: errors.New("unexpected error"),
			wantErr:      errors.New("unexpected error"),
		},
		{
			name:       "success",
			expRepoRes: &model.Metadata{},
			wantRes:    &model.Metadata{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := gen.NewMockmetadataRepository(ctrl)

			c := New(repoMock)
			ctx := context.Background()

			id := "id"
			repoMock.EXPECT().Get(ctx, id).Return(tt.expRepoRes, tt.expRepoError)

			res, err := c.Get(ctx, id)
			assert.Equal(t, tt.wantRes, res, tt.name)
			assert.Equal(t, tt.wantErr, err, tt.name)
		})
	}
}
