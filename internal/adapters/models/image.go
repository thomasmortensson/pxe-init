package models

import (
	"gorm.io/gorm"

	"github.com/thomasmortensson/pxe-init/internal/domain/entities"
)

type Image struct {
	gorm.Model
	Name   string `gorm:"unique"`
	Kernel string
	Initrd string
	Rootfs string
}

// ToDomain converts a postgres model to a domain entity image
func (img *Image) ToDomain() *entities.Image {
	return entities.NewImage(
		img.Name,
		img.Kernel,
		img.Initrd,
		img.Rootfs,
	)
}
