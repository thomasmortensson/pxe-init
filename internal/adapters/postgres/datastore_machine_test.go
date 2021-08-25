package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/thomasmortensson/pxe-init/internal/domain/entities"
)

const (
	testMachineName = "coreos"
	testMAC         = "08-00-27-c3-62-83"
	testUUID        = "2f1214fe-59ba-9f42-a5c5-1af6f124aaf7"

	selectFromMachinesMAC = `SELECT * FROM "machines" WHERE mac = $1 AND "machines"."deleted_at" IS NULL ORDER BY "machines"."id" LIMIT 1`
	selectFromImagesID    = `SELECT * FROM "images" WHERE "images"."id" = $1 AND "images"."deleted_at" IS NULL`
)

var (
	testMachineEntity = entities.NewMachine(
		testMachineName,
		testMAC,
		testUUID,
	)
)

func TestPgDatastore_FindMachineByMAC(t *testing.T) {
	var mock sqlmock.Sqlmock
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	datastore := &PgDatastore{
		logger: logger,
		client: gdb,
	}

	mock.ExpectQuery(selectFromMachinesMAC).WithArgs().WillReturnRows(
		sqlmock.NewRows([]string{
			"name",
			"mac",
			"uuid",
			"image_id",
		}).AddRow(
			testMachineName,
			testMAC,
			testUUID,
			1,
		),
	)

	mock.ExpectQuery(selectFromImagesID).WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"name",
			"kernel",
			"initrd",
			"rootfs",
		}).AddRow(
			1,
			testImageName,
			testKernel,
			testInitrd,
			testRootfs,
		),
	)

	machine, err := datastore.FindMachineByMAC(testImageName)
	assert.Nil(t, err)

	testMachineEntity.SetImage(testImageEntity)
	assert.Equal(t, testMachineEntity, machine)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPgDatastore_FindMachineImageByMAC(t *testing.T) {
	var mock sqlmock.Sqlmock
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	datastore := &PgDatastore{
		logger: logger,
		client: gdb,
	}

	mock.ExpectQuery(selectFromMachinesMAC).WithArgs().WillReturnRows(
		sqlmock.NewRows([]string{
			"name",
			"mac",
			"uuid",
			"image_id",
		}).AddRow(
			testMachineName,
			testMAC,
			testUUID,
			1,
		),
	)

	mock.ExpectQuery(selectFromImagesID).WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"name",
			"kernel",
			"initrd",
			"rootfs",
		}).AddRow(
			1,
			testImageName,
			testKernel,
			testInitrd,
			testRootfs,
		),
	)

	image, err := datastore.FindMachineImageByMAC(testImageName)
	assert.Nil(t, err)

	assert.Equal(t, testImageEntity, image)

	assert.NoError(t, mock.ExpectationsWereMet())
}
