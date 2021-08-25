package usecases

import (
	"go.uber.org/zap"

	"github.com/thomasmortensson/pxe-init/internal/domain/repositories"
)

// RegisterMachineImageUseCase
type RegisterMachineImageUseCase interface {
	Execute(string, string) error
}

// RegisterMachineImage
type RegisterMachineImage struct {
	logger *zap.Logger
	db     repositories.Datastore
}

// NewRegisterMachineImage
func NewRegisterMachineImage(logger *zap.Logger, datastore repositories.Datastore) *RegisterMachineImage {
	return &RegisterMachineImage{
		logger: logger,
		db:     datastore,
	}
}

// Execute the RegisterMachineImage usecase
func (u *RegisterMachineImage) Execute(mac, image string) error {
	_, err := u.db.FindMachineByMAC(mac)
	if u.db.IsNotFound(err) {
		// Machine not found, register machine
		u.logger.Info("First time seeing mac")
		return nil
	} else if err != nil {
		u.logger.Error(
			"Error getting machine details",
			zap.Error(err),
		)
		return err
	}

	// Machine found, setup boot configuration
	u.logger.Info(
		"Found machine. Registering image",
		zap.String("mac", mac),
		zap.String("image", image),
	)
	return nil
}
