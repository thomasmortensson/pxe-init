package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testMachineName = "test vm"
	testMAC         = "08-00-27-c3-62-83"
	testUUID        = "2f1214fe-59ba-9f42-a5c5-1af6f124aaf7"
)

func TestNewMachine(t *testing.T) {
	testMachine := NewMachine(
		testMachineName,
		testMAC,
		testUUID,
	)

	assert.NotNil(t, testMachine)
	assert.Equal(t, testMachineName, testMachine.Name)
	assert.Equal(t, testMAC, testMachine.MAC)
	assert.Equal(t, testUUID, testMachine.UUID)
	assert.Nil(t, testMachine.Image)
}

func TestMachine_SetImage(t *testing.T) {
	testMachine := NewMachine(
		testMachineName,
		testMAC,
		testUUID,
	)
	assert.NotNil(t, testMachine)
	assert.Nil(t, testMachine.Image)

	image := NewImage(testName, testKernel, testInitrd, testRootfs)

	testMachine.SetImage(image)
	assert.NotNil(t, testMachine.Image)
}
