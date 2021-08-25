package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImage_ToDomain(t *testing.T) {
	entityImage := testImage.ToDomain()
	assert.NotNil(t, entityImage)
	assert.Equal(t, testImage.Name, entityImage.Name)
	assert.Equal(t, "http://127.0.0.1/assets/coreos/vmlinuz", entityImage.GetKernel(testServer))
	assert.Equal(t, "http://127.0.0.1/assets/coreos/initrd.img", entityImage.GetInitrd(testServer))
	assert.Equal(t, "http://127.0.0.1/assets/coreos/rootfs.img", entityImage.GetRootfs(testServer))
}
