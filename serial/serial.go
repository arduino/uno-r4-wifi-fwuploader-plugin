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

package serial

import (
	"fmt"
	"time"

	"go.bug.st/serial"
)

// Command represent a command sent through the serial port. This is used to distinguish
// a command that will trigger a specific function of the `commands.ino` sketch.
type Command string

const (
	// RebootCommand puts the board in ESP mode.
	RebootCommand Command = "r\n\r"
	// VersionCommand gets the semver firmware version.
	VersionCommand Command = "v\n\r"
)

// Open used to open the given serial port at 9600 BaudRate
func Open(portAddress string) (serial.Port, error) {
	return serial.Open(portAddress, &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	})
}

// SendCommandAndClose send a command and immediately close the serial port afterwards.
func SendCommandAndClose(port serial.Port, msg Command) error {
	if _, err := port.Write([]byte(string(msg))); err != nil {
		return fmt.Errorf("write to serial port: %v", err)
	}
	if err := port.Close(); err != nil {
		return fmt.Errorf("closing serial port: %v", err)
	}

	return nil
}

// AllPorts returns the list of alla available serial ports.
func AllPorts() (AvailablePorts, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, fmt.Errorf("listing serial ports: %v", err)
	}
	res := map[string]bool{}
	for _, port := range ports {
		res[port] = true
	}
	return res, nil
}

// AvailablePorts represent all the available serial ports
type AvailablePorts map[string]bool

// NewPort every 250ms checks if a new serial port is detected, for a maximum of 10 seconds.
// If a new serial port is detect it's added to the AvailablePorts map and returned as a string.
func (last *AvailablePorts) NewPort() (string, bool, error) {
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		now, err := AllPorts()
		if err != nil {
			return "", false, err
		}

		hasNewPorts := false
		for p := range now {
			if !(*last)[p] {
				hasNewPorts = true
				break
			}
		}

		if hasNewPorts {
			// on OS X, if the port is opened too quickly after it is detected,
			// a "Resource busy" error occurs, add a delay to workaround.
			// This apply to other platforms as well.
			time.Sleep(time.Second)

			// Some boards have a glitch in the bootloader: some user experienced
			// the USB serial port appearing and disappearing rapidly before
			// settling.
			// This check ensure that the port is stable after one second.
			check, err := AllPorts()
			if err != nil {
				return "", false, err
			}
			for p := range check {
				if !(*last)[p] {
					return p, true, nil // Found it!
				}
			}
		}

		*last = now
		time.Sleep(250 * time.Millisecond)
	}

	return "", false, nil
}
