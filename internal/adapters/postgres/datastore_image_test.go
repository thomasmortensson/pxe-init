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
	testImageName = "coreos"
	testKernel    = "/assets/coreos/fedora-coreos-34.20210725.3.0-live-kernel-x86_64"
	testInitrd    = "/assets/coreos/fedora-coreos-34.20210725.3.0-live-initramfs.x86_64.img"
	testRootfs    = "/assets/coreos/fedora-coreos-34.20210725.3.0-live-rootfs.x86_64.img"
)

var (
	testImageEntity = entities.NewImage(testImageName, testKernel, testInitrd, testRootfs)
)

func TestPgDatastore_ListImages(t *testing.T) {
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

	expectedSelectSQL := `SELECT * FROM "images" WHERE "images"."deleted_at" IS NULL`

	mock.ExpectQuery(expectedSelectSQL).WithArgs().WillReturnRows(sqlmock.NewRows([]string{"name", "kernel", "initrd", "rootfs"}))

	images, err := datastore.ListImages()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(images))

	mock.ExpectQuery(expectedSelectSQL).WithArgs().WillReturnRows(
		sqlmock.NewRows([]string{
			"name",
			"kernel",
			"initrd",
			"rootfs",
		}).AddRow(
			testImageName,
			testKernel,
			testInitrd,
			testRootfs,
		),
	)
	images, err = datastore.ListImages()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(images))

	assert.Equal(t, testImageEntity, images[0])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPgDatastore_GetImageByName(t *testing.T) {
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

	expectedSelectSQL := `SELECT * FROM "images" WHERE name = $1 AND "images"."deleted_at" IS NULL ORDER BY "images"."id" LIMIT 1`
	mock.ExpectQuery(expectedSelectSQL).WithArgs(testImageName).WillReturnRows(
		sqlmock.NewRows([]string{
			"name",
			"kernel",
			"initrd",
			"rootfs",
		}).AddRow(
			testImageName,
			testKernel,
			testInitrd,
			testRootfs,
		),
	)

	image, err := datastore.GetImageByName(testImageName)
	assert.Nil(t, err)

	assert.Equal(t, testImageEntity, image)

	assert.NoError(t, mock.ExpectationsWereMet())
}
