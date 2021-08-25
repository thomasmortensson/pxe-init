package repositories

import (
	"github.com/thomasmortensson/pxe-init/internal/domain/entities"
)

// TODO use mockery to generate mocks for Datastore interface for use in testing

// Datastore provides an abstract implementation of an abstract. This is implemented concretely as PgDatastore, a postgres backend
type Datastore interface {
	ListImages() ([]*entities.Image, error)
	GetImageByName(name string) (*entities.Image, error)
	FindMachineByMAC(mac string) (*entities.Machine, error)
	FindMachineImageByMAC(mac string) (*entities.Image, error)
	IsNotFound(err error) bool
}
