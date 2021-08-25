package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/thomasmortensson/pxe-init/internal/domain/entities"
	repositoryMocks "github.com/thomasmortensson/pxe-init/mocks/domain/repositories"
)

const (
	testMac           = "08-00-27-c3-62-83"
	testForwardServer = "127.0.0.1"
)

func TestPxeChainloadExecute(t *testing.T) {
	db := new(repositoryMocks.Datastore)

	logger, err := zap.NewProduction()
	assert.Nil(t, err)

	testImage := entities.NewImage(testName, testKernel, testInitrd, testRootfs)

	usecase := NewPxeChainload(logger, db)

	db.On("IsNotFound", mock.Anything).Return(false)
	db.On("FindMachineImageByMAC", testMac).Return(testImage, nil)

	script, err := usecase.Execute(testMac, testForwardServer)
	assert.Nil(t, err)
	assert.Equal(t, testImage.GetIpxeScript(testForwardServer), script)
}
