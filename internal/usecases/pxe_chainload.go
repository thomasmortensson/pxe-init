package usecases

import (
	"go.uber.org/zap"

	"github.com/thomasmortensson/pxe-init/internal/domain/repositories"
)

// PxeChainloadUseCase
type PxeChainloadUseCase interface {
	Execute(string, string) (string, error)
}

// PxeChainload
type PxeChainload struct {
	logger *zap.Logger
	db     repositories.Datastore
}

// NewPxeChainload
func NewPxeChainload(logger *zap.Logger, datastore repositories.Datastore) *PxeChainload {
	return &PxeChainload{
		logger: logger,
		db:     datastore,
	}
}

// Execute the PxeChainload usecase
func (u *PxeChainload) Execute(mac, forwardAddress string) (string, error) {
	image, err := u.db.FindMachineImageByMAC(mac)
	if u.db.IsNotFound(err) {
		u.logger.Info(
			"Not Found",
			zap.String("mac-address", mac),
		)
		return "", err
	} else if err != nil {
		u.logger.Error(
			"Error encountered generating PXE script",
			zap.Error(err),
		)
		return "", err
	}
	return image.GetIpxeScript(forwardAddress), nil
}
