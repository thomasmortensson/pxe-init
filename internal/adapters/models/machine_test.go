package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testImage = Image{
		Name:   "coreos",
		Kernel: "/assets/coreos/vmlinuz",
		Initrd: "/assets/coreos/initrd.img",
		Rootfs: "/assets/coreos/rootfs.img",
	}

	testMachine = Machine{
		Name:  "test vm",
		MAC:   "08-00-27-c3-62-83",
		UUID:  "2f1214fe-59ba-9f42-a5c5-1af6f124aaf7",
		Image: testImage,
	}

	testServer = "127.0.0.1"
)

func TestMachine_ToDomain(t *testing.T) {
	entityMachine := testMachine.ToDomain()
	assert.Equal(t, testMachine.Name, entityMachine.Name)
	assert.Equal(t, testMachine.MAC, entityMachine.MAC)
	assert.Equal(t, testMachine.UUID, entityMachine.UUID)

	assert.NotNil(t, entityMachine.Image)
}

func TestMachine_LinkedImage(t *testing.T) {
	entityImage := testMachine.LinkedImage()
	assert.NotNil(t, entityImage)
	assert.Equal(t, testImage.Name, entityImage.Name)
	assert.Equal(t, "http://127.0.0.1/assets/coreos/vmlinuz", entityImage.GetKernel(testServer))
	assert.Equal(t, "http://127.0.0.1/assets/coreos/initrd.img", entityImage.GetInitrd(testServer))
	assert.Equal(t, "http://127.0.0.1/assets/coreos/rootfs.img", entityImage.GetRootfs(testServer))
}
