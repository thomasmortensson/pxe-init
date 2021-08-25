package postgres

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/thomasmortensson/pxe-init/internal/adapters/models"
)

type ConnProvider struct {
	ConnString string
}

// NewConnProvider constructs a new provider for the database specified in the config
func NewConnProvider(
	host string,
	port int,
	user,
	password,
	databaseName,
	sslMode string,
) *ConnProvider {
	return &ConnProvider{
		ConnString: fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			host,
			port,
			user,
			password,
			databaseName,
			sslMode,
		),
	}
}

type PgDatastore struct {
	logger *zap.Logger
	client *gorm.DB
}

// NewPgDatastore initializes a new datastore object implemented using the postgres driver.
func NewPgDatastore(logger *zap.Logger, connProvider *ConnProvider) (*PgDatastore, error) {
	client, err := gorm.Open(postgres.Open(connProvider.ConnString), &gorm.Config{
		// By default GORM outputs in non structured mode. Prefer to check errors internally
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Error(
			"Failed to setup DB client",
			zap.Error(err),
		)
		return nil, err
	}

	return &PgDatastore{
		logger: logger,
		client: client,
	}, nil
}

// AutoMigrate takes the existing database models and migrates to match the specified models for this process
func (d *PgDatastore) AutoMigrate() error {
	err := d.client.AutoMigrate(&models.Image{})
	if err != nil {
		d.logger.Error(
			"Failed to migrate image table",
			zap.Error(err),
		)
		return err
	}

	err = d.client.AutoMigrate(&models.Machine{})
	if err != nil {
		d.logger.Error(
			"Failed to migrate image table",
			zap.Error(err),
		)
		return err
	}
	return nil
}

// IsNotFound inspects a returned gorm results.Error structure and determines if a row has not been found
func (d *PgDatastore) IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
