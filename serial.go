package main

import (
	"fmt"

	"go.bug.st/serial"
)

type serialCommand string

const (
	rebootCommand  serialCommand = "r\n\r"
	versionCommand serialCommand = "v\n\r"
)

type serialPort string

func openSerialPort(portAddress serialPort) (serial.Port, error) {
	return serial.Open(string(portAddress), &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	})
}

func sendSerialCommandAndClose(port serial.Port, msg serialCommand) error {
	if _, err := port.Write([]byte(string(msg))); err != nil {
		return fmt.Errorf("write to serial port: %v", err)
	}
	if err := port.Close(); err != nil {
		return fmt.Errorf("closing serial port: %v", err)
	}

	return nil
}
