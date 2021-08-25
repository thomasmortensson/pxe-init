package models

import (
	"gorm.io/gorm"

	"github.com/thomasmortensson/pxe-init/internal/domain/entities"
)

type Machine struct {
	gorm.Model
	Name    string `gorm:"unique"`
	MAC     string `gorm:"unique"`
	UUID    string `gorm:"unique"`
	ImageID int
	Image   Image `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// ToDomain converts a postgres model to a domain entity machine
func (m *Machine) ToDomain() *entities.Machine {
	machine := entities.NewMachine(
		m.Name,
		m.MAC,
		m.UUID,
	)
	machine.SetImage(m.LinkedImage())
	return machine
}

// LinkedImage converts a postgres image model linked via foreign key relationship
// from a machine to the corresponding domain entity image
func (m *Machine) LinkedImage() *entities.Image {
	return m.Image.ToDomain()
}
