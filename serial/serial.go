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

type PortDetector struct {
	ports              map[string]bool
	deadLine           time.Duration
	detectionFrequency time.Duration
}

type PortDetectorOption func(*PortDetector)

func WithDeadline(deadline time.Duration) PortDetectorOption {
	return func(h *PortDetector) {
		h.deadLine = deadline
	}
}

func WithDetectionFrequency(frequency time.Duration) PortDetectorOption {
	return func(h *PortDetector) {
		h.detectionFrequency = frequency
	}
}

func NewPortDetector(opts ...PortDetectorOption) (*PortDetector, error) {
	const (
		deadline  = 10 * time.Second
		frequency = 250 * time.Millisecond
	)

	res, err := getAllAvailablePorts()
	if err != nil {
		return nil, err
	}
	d := &PortDetector{ports: res, deadLine: deadline, detectionFrequency: frequency}

	for _, opt := range opts {
		opt(d)
	}

	return d, nil
}

func (d *PortDetector) DetectNewPorts() ([]string, bool, error) {
	deadline := time.Now().Add(d.deadLine)
	for time.Now().Before(deadline) {
		now, err := getAllAvailablePorts()
		if err != nil {
			return nil, false, err
		}

		hasNewPorts := false
		for p := range now {
			if !(d.ports)[p] {
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
			check, err := getAllAvailablePorts()
			if err != nil {
				return nil, false, err
			}

			newPorts := []string{}
			newPortFound := false
			for p := range check {
				if !(d.ports)[p] {
					newPortFound = true
					newPorts = append(newPorts, p)
				}
			}
			if newPortFound {
				return newPorts, true, nil // Found it!
			}
		}

		d.ports = now
		time.Sleep(d.detectionFrequency)
	}

	return nil, false, nil
}

func (d *PortDetector) Reset() error {
	res, err := getAllAvailablePorts()
	if err != nil {
		return err
	}
	d.ports = res
	return nil
}

func getAllAvailablePorts() (map[string]bool, error) {
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
