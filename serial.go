package main

import (
	"fmt"

	"go.bug.st/serial"
)

type serialCommand string

const (
	rebootCommand  serialCommand = "r"
	versionCommand serialCommand = "v"
)

type serialPort string

func sendSerialCommandAndClose(portAddress serialPort, msg serialCommand) error {
	port, err := serial.Open(string(portAddress), &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		return fmt.Errorf("open serial port: %v", err)
	}
	if _, err := port.Write([]byte(string(msg) + "\n\r")); err != nil {
		return fmt.Errorf("write to serial port: %v", err)
	}
	if err := port.Close(); err != nil {
		return fmt.Errorf("closing serial port: %v", err)
	}

	return nil
}
