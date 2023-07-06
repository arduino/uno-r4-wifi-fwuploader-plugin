package main

import (
	"fmt"

	"github.com/karalabe/hid"
)

const (
	vid uint16 = 0x2341
	pid uint16 = 0x1002
)

func rebootUsingHID() error {
	d, err := hid.Open(vid, pid)
	if err != nil {
		return fmt.Errorf("open HID: %v", err)
	}

	b := make([]byte, 65)
	b[0] = 0
	b[1] = 0xAA
	if _, err := d.SendFeatureReport(b); err != nil {
		return err
	}

	return nil
}

