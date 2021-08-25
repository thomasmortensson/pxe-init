package entities

import (
	"fmt"
)

type Image struct {
	Name   string
	kernel string
	initrd string
	rootfs string
}

// NewImage generates a new Image entity structure
func NewImage(name, kernel, initrd, rootfs string) *Image {
	return &Image{
		Name:   name,
		kernel: kernel,
		initrd: initrd,
		rootfs: rootfs,
	}
}

// GetKernel returns the bootable kernel parameter for available asset on server
func (img *Image) GetKernel(server string) string {
	return fmt.Sprintf("http://%v%v", server, img.kernel)
}

// GetInitrd returns the initial ramdisk parameter initrd for available asset on server
func (img *Image) GetInitrd(server string) string {
	return fmt.Sprintf("http://%v%v", server, img.initrd)
}

// GetInitrd returns the initial ramdisk parameter initrd for available asset on server
func (img *Image) GetRootfs(server string) string {
	return fmt.Sprintf("http://%v%v", server, img.rootfs)
}

// GetIpxeScript templates a pre-defined ipxe script with the available server and image parameters
func (img *Image) GetIpxeScript(server string) string {
	// TODO this is obviously built for coreos.
	// This would be extended using enums defining image type with specific scripts defined for each OS

	// This could also be improved by using text/template, I've run out of time.
	return "#!ipxe\n" +
		fmt.Sprintf("kernel %v ", img.GetKernel(server)) +
		fmt.Sprintf("coreos.live.rootfs_url=%v ", img.GetRootfs(server)) +
		"coreos.first_boot=1 coreos.autologin\n" +
		fmt.Sprintf("initrd %v", img.GetInitrd(server)) +
		"\nboot"
}
