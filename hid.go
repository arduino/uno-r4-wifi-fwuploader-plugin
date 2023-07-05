package main

import "github.com/karalabe/hid"

const (
	vid uint16 = 0x2341
	pid uint16 = 0x1002
)

func openHID() (*hid.Device, error) {
	return hid.Open(vid, pid)
}

func reboot(d *hid.Device) error {
	b := make([]byte, 65)
	b[0] = 0
	b[1] = 0xAA
	if _, err := d.SendFeatureReport(b); err != nil {
		return err
	}
	return nil
}
