package main

import "github.com/sstallion/go-hid"

const (
	vid uint16 = 0x2341
	pid uint16 = 0x1002
)

func openFirstHID() (*hid.Device, error) {
	return openHID("")
}

func openHID(port string) (*hid.Device, error) {
	if err := hid.Init(); err != nil {
		return nil, err
	}
	if port == "" {
		return hid.OpenFirst(vid, pid)
	}
	return hid.OpenPath(port)
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
