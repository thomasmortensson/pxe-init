package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testName   = "coreos"
	testKernel = "/assets/coreos/fedora-coreos-34.20210725.3.0-live-kernel-x86_64"
	testInitrd = "/assets/coreos/fedora-coreos-34.20210725.3.0-live-initramfs.x86_64.img"
	testRootfs = "/assets/coreos/fedora-coreos-34.20210725.3.0-live-rootfs.x86_64.img"

	testServer = "pxe-init.mortcloud.com"
)

var (
	testImage = NewImage(testName, testKernel, testInitrd, testRootfs)
)

func TestNewImage(t *testing.T) {
	image := NewImage(testName, testKernel, testInitrd, testRootfs)
	assert.Equal(t, testName, image.Name)
	assert.Equal(t, testKernel, image.kernel)
	assert.Equal(t, testInitrd, image.initrd)
	assert.Equal(t, testRootfs, image.rootfs)
}

func TestImage_GetKernel(t *testing.T) {
	kernelLoc := testImage.GetKernel(testServer)
	assert.Equal(t, "http://pxe-init.mortcloud.com/assets/coreos/fedora-coreos-34.20210725.3.0-live-kernel-x86_64", kernelLoc)
}

func TestImage_GetInitrd(t *testing.T) {
	initrdLoc := testImage.GetInitrd(testServer)
	assert.Equal(t, "http://pxe-init.mortcloud.com/assets/coreos/fedora-coreos-34.20210725.3.0-live-initramfs.x86_64.img", initrdLoc)
}

func TestImage_GetRootfs(t *testing.T) {
	rootfsLoc := testImage.GetRootfs(testServer)
	assert.Equal(t, "http://pxe-init.mortcloud.com/assets/coreos/fedora-coreos-34.20210725.3.0-live-rootfs.x86_64.img", rootfsLoc)
}
