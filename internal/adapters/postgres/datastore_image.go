package postgres

import (
	"go.uber.org/zap"

	"github.com/thomasmortensson/pxe-init/internal/adapters/models"
	"github.com/thomasmortensson/pxe-init/internal/domain/entities"
)

// listImages returns a list of all available image models
func (d *PgDatastore) listImages() ([]models.Image, error) {
	var images []models.Image
	result := d.client.Find(&images)
	if result.Error != nil {
		d.logger.Error(
			"Error encountered listing images",
			zap.Error(result.Error),
		)
		return nil, result.Error
	}
	return images, nil
}

// ListImages returns a list of all available entity models converted from domain
func (d *PgDatastore) ListImages() ([]*entities.Image, error) {
	modelImages, err := d.listImages()
	if err != nil {
		return nil, err
	}

	var entityImages = make([]*entities.Image, len(modelImages))
	for i, image := range modelImages {
		entityImages[i] = image.ToDomain()
	}

	return entityImages, nil
}

// getImageByName returns a single image model (or ErrRecordNotFound) representing the image with given name
func (d *PgDatastore) getImageByName(name string) (*models.Image, error) {
	var image models.Image
	result := d.client.Where("name = ?", name).First(&image)

	if d.IsNotFound(result.Error) {
		d.logger.Info(
			"Unable to find image with name",
			zap.String("name", name),
		)
		return nil, result.Error
	} else if result.Error != nil {
		d.logger.Error(
			"Error retrieving imabe by name",
			zap.Error(result.Error),
		)
		return nil, result.Error
	}
	return &image, nil
}

// GetImageByName returns a single image domain entity (or ErrRecordNotFound) representing the image with given name
func (d *PgDatastore) GetImageByName(name string) (*entities.Image, error) {
	modelImage, err := d.getImageByName(name)
	if err != nil {
		// Logging taken care of in low-level function
		return nil, err
	}

	return modelImage.ToDomain(), nil
}
