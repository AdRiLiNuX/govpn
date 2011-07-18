package main

import (
	"exec"
	"log"
	"os"
	"syscall"
	"unsafe"
)

func (tt *tuntap) open() {
	deviceFile := "/dev/net/tun"
	fd, err := os.OpenFile(deviceFile, os.O_RDWR, 0)
	if err != nil {
		log.Fatalf("Note: Cannot open TUN/TAP dev %s: %v", deviceFile, err)
	}
	tt.fd = fd

	ifr := make([]byte, 18)
	ifr[17] = 0x10 // IFF_NO_PI
	ifr[16] = 0x02 // IFF_TAP

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(tt.fd.Fd()), uintptr(0x400454ca), // TUNSETIFF
		uintptr(unsafe.Pointer(&ifr[0])))
	if errno != 0 {
		log.Fatalf("Cannot ioctl TUNSETIFF: %v", os.Errno(errno))
	}

	tt.actualName = string(ifr)
	log.Printf("TUN/TAP device %s opened.", tt.actualName)
}

func (tt *tuntap) ifconfig() {
	cmd := exec.Command("/sbin/ifconfig", tt.actualName,
		tt.address.IP.String(), "netmask",
		tt.netmask.IP.String(), "mtu", "1500")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Linux ifconfig failed: %v.", err)
	}
}
