package usecases

import (
	"go.uber.org/zap"

	"github.com/thomasmortensson/pxe-init/internal/domain/entities"
	"github.com/thomasmortensson/pxe-init/internal/domain/repositories"
)

// ListImagesUseCase
type ListImagesUseCase interface {
	Execute() ([]*entities.Image, error)
}

// ListImages
type ListImages struct {
	logger *zap.Logger
	db     repositories.Datastore
}

// NewListImages
func NewListImages(logger *zap.Logger, datastore repositories.Datastore) *ListImages {
	return &ListImages{
		logger: logger,
		db:     datastore,
	}
}

// Execute the ListImages usecase
func (u *ListImages) Execute() ([]*entities.Image, error) {
	images, err := u.db.ListImages()
	if err != nil {
		u.logger.Error(
			"Failed to list images",
			zap.Error(err),
		)
		return nil, err
	}
	return images, nil
}
