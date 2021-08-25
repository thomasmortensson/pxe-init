package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/thomasmortensson/pxe-init/internal/domain/entities"
	repositoryMocks "github.com/thomasmortensson/pxe-init/mocks/domain/repositories"
)

const (
	testName   = "coreos"
	testKernel = "/assets/coreos/fedora-coreos-34.20210725.3.0-live-kernel-x86_64"
	testInitrd = "/assets/coreos/fedora-coreos-34.20210725.3.0-live-initramfs.x86_64.img"
	testRootfs = "/assets/coreos/fedora-coreos-34.20210725.3.0-live-rootfs.x86_64.img"
)

func TestListImagesExecute(t *testing.T) {
	db := new(repositoryMocks.Datastore)

	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	var testImages = make([]*entities.Image, 1)
	testImages[0] = entities.NewImage(testName, testKernel, testInitrd, testRootfs)

	usecase := NewListImages(logger, db)

	db.On("ListImages").Return(testImages, nil)

	images, err := usecase.Execute()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(images))
	assert.Equal(t, testImages[0], images[0])
}
