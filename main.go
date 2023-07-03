package main

import (
	"fmt"
	"time"

	"github.com/arduino/arduino-cli/executils"
	helper "github.com/arduino/fwuploader-plugin-helper"
	"github.com/arduino/go-paths-helper"
	"github.com/sstallion/go-hid"
	semver "go.bug.st/relaxed-semver"
)

const (
	pluginName    = "uno-r4-wifi-fwuploader"
	pluginVersion = "1.0.0"
)

func main() {
	helper.RunPlugin(&unoR4WifiPlugin{})
}

type unoR4WifiPlugin struct{}

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

	d, err := openFirstHID()
	if err != nil {
		return err
	}

	if err := reboot(d); err != nil {
		return err
	}

	if err := hid.Exit(); err != nil {
		return err
	}

	// Wait a bit before flashing the firmware to allow the board to become available again.
	time.Sleep(3 * time.Second)

	cmd, err := executils.NewProcess([]string{}, "espflash", "flash", firmwarePath.String(), "-p", portAddress)
	if err != nil {
		return err
	}
	cmd.RedirectStderrTo(feedback.Err())
	cmd.RedirectStdoutTo(feedback.Out())
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Fprintf(feedback.Out(), "\nUpload completed! You can now detach the board.\n")
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

	d, err := openFirstHID()
	if err != nil {
		return err
	}

	if err := reboot(d); err != nil {
		return err
	}

	if err := hid.Exit(); err != nil {
		return err
	}

	// Wait a bit before flashing the certificate to allow the board to become available again.
	time.Sleep(3 * time.Second)

	cmd, err := executils.NewProcess([]string{}, "espflash", "write-bin", "-p", portAddress, "-b", "921600", "0x3C0000", certificatePath.String())
	if err != nil {
		return err
	}
	cmd.RedirectStderrTo(feedback.Err())
	cmd.RedirectStdoutTo(feedback.Out())
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Fprintf(feedback.Out(), "\nUpload completed! You can now detach the board.\n")
	return nil
}

// GetFirmwareVersion retrieve the firmware version installed on the board
func (p *unoR4WifiPlugin) GetFirmwareVersion(portAddress string, feedback *helper.PluginFeedback) (*semver.RelaxedVersion, error) {
	d, err := openHID(portAddress) // if port is empty, fallback to vid+pid
	if err != nil {
		return nil, err
	}
	defer hid.Exit()

	buf := make([]byte, 65)
	if _, err := d.GetFeatureReport(buf); err != nil {
		return nil, err
	}
	return semver.ParseRelaxed(fmt.Sprintf("%d.%d.%d", buf[1], buf[2], buf[3])), nil
}
