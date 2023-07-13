// This file is part of uno-r4-wifi-fwuploader-plugin.
//
// Copyright (c) 2023 Arduino LLC.  All right reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
		return fmt.Errorf("send HID command: %v", err)
	}

	return nil
}
