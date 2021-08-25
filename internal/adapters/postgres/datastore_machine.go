package postgres

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/thomasmortensson/pxe-init/internal/adapters/models"
	"github.com/thomasmortensson/pxe-init/internal/domain/entities"
)

// findMachineByMAC returns the first machine model where mac is matched
func (d *PgDatastore) findMachineByMAC(mac string) (*models.Machine, error) {
	var machine models.Machine
	result := d.client.Preload("Image").Where("mac = ?", mac).First(&machine)

	if d.IsNotFound(result.Error) {
		d.logger.Info(
			"Unable to find machine with mac",
			zap.String("mac-address", mac),
		)
		return nil, result.Error
	} else if result.Error != nil {
		d.logger.Error(
			"Error encountered finding machine in DB",
			zap.Error(result.Error),
		)
		return nil, result.Error
	}

	if result.RowsAffected == 1 {
		return &machine, nil
	}
	// Shouldn't get here as we set a limit on the query but better to be defensive
	return nil, fmt.Errorf("returned %v rows. Expecting 1", result.RowsAffected)
}

// FindMachineByMAC returns the first machine domain entity where mac is matched
func (d *PgDatastore) FindMachineByMAC(mac string) (*entities.Machine, error) {
	modelMachine, err := d.findMachineByMAC(mac)
	if err != nil {
		// Logging taken care of in low-level function
		return nil, err
	}
	return modelMachine.ToDomain(), nil
}

// FindMachineImageByMAC returns the linked image domain entity for first machine where mac is matched
func (d *PgDatastore) FindMachineImageByMAC(mac string) (*entities.Image, error) {
	modelMachine, err := d.findMachineByMAC(mac)
	if err != nil {
		// Logging taken care of in low-level function
		return nil, err
	}
	return modelMachine.LinkedImage(), nil
}
