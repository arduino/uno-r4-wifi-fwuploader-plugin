module github.com/arduino/uno-r4-wifi-fwuploader-plugin

go 1.20

// using our forked version to fix HID on windows
replace github.com/karalabe/hid => github.com/bcmi-labs/hid v0.0.0-20230703110227-931f677e7f17

require (
	github.com/arduino/arduino-cli v0.0.0-20230706071323-df12786440c1
	github.com/arduino/fwuploader-plugin-helper v0.0.0-20230706160803-17fd16a4651c
	github.com/arduino/go-paths-helper v1.9.1
	github.com/karalabe/hid v0.0.0-00010101000000-000000000000
	go.bug.st/relaxed-semver v0.11.0
	go.bug.st/serial v1.5.0
)

require (
	github.com/creack/goselect v0.1.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/leonelquinteros/gotext v1.5.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
