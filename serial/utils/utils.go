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

// Those function are token from https://github.com/arduino/arduino-cli/blob/master/arduino/serialutils/serialutils.go
// that's because we don't have the `tr` here and importing the serialutils from the cli will lead to a panic
package utils

import (
	"fmt"
	"runtime"
	"time"

	"go.bug.st/serial"
)

// TouchSerialPortAt1200bps open and close the serial port at 1200 bps. This
// is used on many Arduino boards as a signal to put the board in "bootloader"
// mode.
func TouchSerialPortAt1200bps(port string) error {
	// Open port
	p, err := serial.Open(port, &serial.Mode{BaudRate: 1200})
	if err != nil {
		return fmt.Errorf("opening port at 1200bps")
	}

	if runtime.GOOS != "windows" {
		// This is not required on Windows
		// TODO: Investigate if it can be removed for other OS too

		// Set DTR to false
		if err = p.SetDTR(false); err != nil {
			p.Close()
			return fmt.Errorf("setting DTR to OFF")
		}
	}

	// Close serial port
	p.Close()

	// Scanning for available ports seems to open the port or
	// otherwise assert DTR, which would cancel the WDT reset if
	// it happens within 250 ms. So we wait until the reset should
	// have already occurred before going on.
	time.Sleep(500 * time.Millisecond)

	return nil
}
