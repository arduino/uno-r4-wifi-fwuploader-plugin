package serial

import (
	"fmt"
	"time"

	"go.bug.st/serial"
)

type SerialCommand string

const (
	RebootCommand  SerialCommand = "r\n\r"
	VersionCommand SerialCommand = "v\n\r"
)

type SerialPort string

func Open(portAddress SerialPort) (serial.Port, error) {
	return serial.Open(string(portAddress), &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	})
}

func SendCommandAndClose(port serial.Port, msg SerialCommand) error {
	if _, err := port.Write([]byte(string(msg))); err != nil {
		return fmt.Errorf("write to serial port: %v", err)
	}
	if err := port.Close(); err != nil {
		return fmt.Errorf("closing serial port: %v", err)
	}

	return nil
}

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

type AvailablePorts map[string]bool

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
