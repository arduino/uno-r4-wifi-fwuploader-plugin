package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/arduino/arduino-cli/arduino/serialutils"
	"github.com/arduino/arduino-cli/executils"
	helper "github.com/arduino/fwuploader-plugin-helper"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/uno-r4-wifi-fwuploader-plugin/certificate"
	"github.com/arduino/uno-r4-wifi-fwuploader-plugin/serial"
	semver "go.bug.st/relaxed-semver"
)

const (
	pluginName = "uno-r4-wifi-fwuploader"
)

var (
	versionString = "0.0.0-git"
	commit        = ""
	date          = ""
)

//go:embed sketches/commands/build/arduino.renesas_uno.unor4wifi/commands.ino.bin
var commandSketchBinary embed.FS

func main() {
	espflashPath, err := helper.FindToolPath("espflash", semver.MustParse("2.0.0"))
	if err != nil {
		log.Fatalln("Couldn't find espflash@2.0.0 binary")
	}
	bossacPath, err := helper.FindToolPath("bossac", semver.MustParse("1.9.1-arduino5"))
	if err != nil {
		log.Fatalln("Couldn't find bossac@1.9.1-arduino5 binary")
	}

	helper.RunPlugin(&unoR4WifiPlugin{
		espflashBin: espflashPath.Join("espflash"),
		bossacBin:   bossacPath.Join("bossac"),
	})
}

type unoR4WifiPlugin struct {
	espflashBin *paths.Path
	bossacBin   *paths.Path
}

var _ helper.Plugin = (*unoR4WifiPlugin)(nil)

// GetPluginInfo returns information about the plugin
func (p *unoR4WifiPlugin) GetPluginInfo() *helper.PluginInfo {
	return &helper.PluginInfo{
		Name:    pluginName,
		Version: semver.MustParse(versionString),
	}
}

// UploadFirmware performs a firmware upload on the board
func (p *unoR4WifiPlugin) UploadFirmware(portAddress string, firmwarePath *paths.Path, feedback *helper.PluginFeedback) error {
	if portAddress == "" {
		return fmt.Errorf("invalid port address")
	}
	if firmwarePath == nil || firmwarePath.IsDir() || !firmwarePath.Exist() {
		return fmt.Errorf("invalid firmware path")
	}

	if err := p.reboot(&portAddress, feedback); err != nil {
		return fmt.Errorf("reboot mode: %v", err)
	}

	cmd, err := executils.NewProcess([]string{}, p.espflashBin.String(), "write-bin", "-p", portAddress, "-b", "115200", "0x0", firmwarePath.String())
	if err != nil {
		return err
	}
	cmd.RedirectStderrTo(feedback.Err())
	cmd.RedirectStdoutTo(feedback.Out())
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Fprintf(feedback.Out(), "\nUpload completed! To complete the update process please disconnect and then reconnect the board.\n")
	return nil
}

// UploadCertificate performs a certificate upload on the board.
func (p *unoR4WifiPlugin) UploadCertificate(portAddress string, certificatePath *paths.Path, feedback *helper.PluginFeedback) error {
	if portAddress == "" {
		return fmt.Errorf("invalid port address")
	}
	if certificatePath == nil || certificatePath.IsDir() || !certificatePath.Exist() {
		return fmt.Errorf("invalid certificate path")
	}
	fmt.Fprintf(feedback.Out(), "Uploading certificates to %s...\n", portAddress)

	if err := p.reboot(&portAddress, feedback); err != nil {
		return fmt.Errorf("reboot mode: %v", err)
	}

	crtBundle, err := certificate.PemToCrt(certificatePath)
	if err != nil {
		return fmt.Errorf("certificate: %v", err)
	}

	// The certificate must be in crt format and be multiple of 4, otherwise `espflash` won't work!
	// (https://github.com/esp-rs/espflash/issues/440)
	for (len(crtBundle) & 3) != 0 {
		crtBundle = append(crtBundle, 0xff)
	}

	crtFile, err := paths.WriteToTempFile(crtBundle, paths.TempDir(), "fw-uploader-uno-r4-wifi-cert")
	if err != nil {
		return err
	}
	defer crtFile.Remove()

	cmd, err := executils.NewProcess([]string{}, p.espflashBin.String(), "write-bin", "-p", portAddress, "-b", "921600", "0x3C0000", crtFile.String())
	if err != nil {
		return err
	}
	cmd.RedirectStderrTo(feedback.Err())
	cmd.RedirectStdoutTo(feedback.Out())
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Fprintf(feedback.Out(), "\nUpload completed! To complete the update process please disconnect and then reconnect the board.\n")
	return nil
}

// GetFirmwareVersion retrieve the firmware version installed on the board
func (p *unoR4WifiPlugin) GetFirmwareVersion(portAddress string, feedback *helper.PluginFeedback) (*semver.RelaxedVersion, error) {
	if err := p.uploadCommandsSketch(portAddress, feedback); err != nil {
		return nil, err
	}

	port, err := serial.Open(serial.Port(portAddress))
	if err != nil {
		return nil, err
	}
	defer port.Close()

	if _, err := port.Write([]byte(string(serial.VersionCommand))); err != nil {
		return nil, fmt.Errorf("write to serial port: %v", err)
	}

	var version string
	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		version = scanner.Text()
		break
	}

	return semver.ParseRelaxed(version), nil
}

func (p *unoR4WifiPlugin) reboot(portAddress *string, feedback *helper.PluginFeedback) error {
	// Will be used later to check if the OS changed the serial port.
	allSerialPorts, err := serial.AllPorts()
	if err != nil {
		return err
	}

	if err := p.uploadCommandsSketch(*portAddress, feedback); err != nil {
		return fmt.Errorf("upload commands sketch: %v", err)
	}

	port, err := serial.Open(serial.Port(*portAddress))
	if err != nil {
		return err
	}
	if err := serial.SendCommandAndClose(port, serial.RebootCommand); err != nil {
		return err
	}

	fmt.Fprintf(feedback.Out(), "Waiting to flash the binary...\n")

	time.Sleep(3 * time.Second)

	// On Windows, when a board is successfully rebooted in esp32 mode, it will change the serial port.
	// Every 250ms we're watching for new ports, if a new one is found we return that otherwise
	// we'll wait the the 10 seconds timeout expiration.
	newPort, changed, err := allSerialPorts.NewPort()
	if err != nil {
		return err
	}
	if changed {
		*portAddress = newPort
	}

	// Older firmware version (v0.1.0) do not support rebooting using the command sketch.
	// So we use HID to reboot. We're consciosly ignoring the error because for boards
	// running a firmware >= v0.2.0 will alaways throw an HID error as we're already in
	// esp32 mode.
	_ = rebootUsingHID()

	time.Sleep(3 * time.Second)

	// On Windows, when a board is successfully rebooted in esp32 mode, it will change the serial port.
	// Every 250ms we're watching for new ports, if a new one is found we return that otherwise
	// we'll wait the the 10 seconds timeout expiration.
	newPort, changed, err = allSerialPorts.NewPort()
	if err != nil {
		return err
	}
	if changed {
		*portAddress = newPort
	}

	return nil
}

func (p *unoR4WifiPlugin) uploadCommandsSketch(portAddress string, feedback *helper.PluginFeedback) error {
	rebootData, err := commandSketchBinary.ReadFile("sketches/commands/build/arduino.renesas_uno.unor4wifi/commands.ino.bin")
	if err != nil {
		return err
	}
	rebootFile, err := paths.WriteToTempFile(rebootData, paths.TempDir(), "fw-uploader-uno-r4-wifi")
	if err != nil {
		return err
	}
	defer rebootFile.Remove()

	if _, err = serialutils.Reset(portAddress, false, nil, false); err != nil {
		return err
	}
	cmd, err := executils.NewProcess(nil, p.bossacBin.String(), "--port="+portAddress, "-U", "-e", "-w", rebootFile.String(), "-R")
	if err != nil {
		return err
	}
	cmd.RedirectStderrTo(feedback.Err())
	cmd.RedirectStdoutTo(feedback.Out())
	if err := cmd.Run(); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	return nil
}
