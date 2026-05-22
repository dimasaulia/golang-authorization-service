package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/releasenotes/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type ReleaseNoteService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.ReleaseNote, error)
	FindByID(ctx context.Context, id int64) (*entities.ReleaseNote, error)
	Create(ctx context.Context, request dto.CreateReleaseNoteRequest) (*entities.ReleaseNote, error)
	Update(ctx context.Context, id int64, request dto.UpdateReleaseNoteRequest) (*entities.ReleaseNote, error)
	Delete(ctx context.Context, id int64, deletedBy *int64) error
}
