package main

import (
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/arduino/arduino-cli/arduino/serialutils"
	"github.com/arduino/arduino-cli/executils"
	helper "github.com/arduino/fwuploader-plugin-helper"
	"github.com/arduino/go-paths-helper"
	semver "go.bug.st/relaxed-semver"
)

const (
	pluginName    = "uno-r4-wifi-fwuploader"
	pluginVersion = "1.0.0"
)

//go:embed sketches/commands/build/arduino.renesas_uno.unor4wifi/commands.ino.bin
var rebootBinary embed.FS

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
		Version: semver.MustParse(pluginVersion),
	}
}

// UploadFirmware performs a firmware upload on the board
func (p *unoR4WifiPlugin) UploadFirmware(portAddress string, firmwarePath *paths.Path, feedback *helper.PluginFeedback) error {
	if portAddress == "" {
		fmt.Fprintln(feedback.Err(), "Port address not specified")
		return fmt.Errorf("invalid port address")
	}
	if firmwarePath == nil || firmwarePath.IsDir() {
		fmt.Fprintln(feedback.Err(), "Invalid firmware path")
		return fmt.Errorf("invalid firmware path")
	}

	if err := p.reboot(portAddress, feedback); err != nil {
		return fmt.Errorf("reboot mode: %v", err)
	}

	cmd, err := executils.NewProcess([]string{}, p.espflashBin.String(), "flash", firmwarePath.String(), "-p", portAddress)
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

// UploadCertificate performs a certificate upload on the board. The certificate must be in crt format
// and be multiple of 4, otherwise `espflash` won't work!
func (p *unoR4WifiPlugin) UploadCertificate(portAddress string, certificatePath *paths.Path, feedback *helper.PluginFeedback) error {
	if portAddress == "" {
		fmt.Fprintln(feedback.Err(), "Port address not specified")
		return fmt.Errorf("invalid port address")
	}
	if certificatePath == nil || certificatePath.IsDir() {
		fmt.Fprintln(feedback.Err(), "Invalid certificate path")
		return fmt.Errorf("invalid certificate path")
	}
	fmt.Fprintf(feedback.Out(), "Uploading certificates to %s...\n", portAddress)

	if err := p.reboot(portAddress, feedback); err != nil {
		return fmt.Errorf("reboot mode: %v", err)
	}

	cmd, err := executils.NewProcess([]string{}, p.espflashBin.String(), "write-bin", "-p", portAddress, "-b", "921600", "0x3C0000", certificatePath.String())
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
	d, err := openHID()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 65)
	if _, err := d.GetFeatureReport(buf); err != nil {
		return nil, err
	}
	return semver.ParseRelaxed(fmt.Sprintf("%d.%d.%d", buf[1], buf[2], buf[3])), nil
}

func (p *unoR4WifiPlugin) reboot(portAddress string, feedback *helper.PluginFeedback) error {
	rebootData, err := rebootBinary.ReadFile("sketches/commands/build/arduino.renesas_uno.unor4wifi/commands.ino.bin")
	if err != nil {
		return err
	}
	rebootFile, err := paths.WriteToTempFile(rebootData, paths.TempDir(), "fw-uploader-uno-r4-wifi")
	if err != nil {
		return err
	}
	defer rebootFile.Remove()

	newPortAddress, err := serialutils.Reset(portAddress, false, nil, false)
	if err != nil {
		return err
	}
	if newPortAddress != "" {
		portAddress = newPortAddress
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

	sendSerialCommandAndClose(serialPort(portAddress), rebootCommand)

	time.Sleep(3 * time.Second)

	// try to use HID to reboot in case firmware version is v0.1.0
	{
		d, err := openHID()
		if err != nil {
			return fmt.Errorf("open HID: %v", err)
		}

		if err := reboot(d); err != nil {
			return fmt.Errorf("reboot HID: %v", err)
		}

		time.Sleep(3 * time.Second)
	}

	return nil
}
