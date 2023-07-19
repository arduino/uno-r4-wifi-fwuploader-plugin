# uno-r4-wifi-fwuploader-plugin

[![Release status](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/release-go-task.yml/badge.svg)](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/release-go-task.yml)

[![Check Go Dependencies status](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-go-dependencies-task.yml/badge.svg)](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-go-dependencies-task.yml)
[![Check Go status](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-go-task.yml/badge.svg)](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-go-task.yml)
[![Sync Labels status](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/sync-labels.yml/badge.svg)](https://github.com/arduino/fwuploader-plugin-helper/actions/workflows/sync-labels.yml)
[![Test Go status](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/test-go-task.yml/badge.svg)](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/test-go-task.yml)
[![Check Markdown status](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-markdown-task.yml/badge.svg)](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-markdown-task.yml)
[![Check License status](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-license.yml/badge.svg)](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-license.yml)
[![Check Taskfiles status](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-taskfiles.yml/badge.svg)](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-taskfiles.yml)
[![Check Workflows status](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-workflows-task.yml/badge.svg)](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/actions/workflows/check-workflows-task.yml)
[![Codecov](https://codecov.io/gh/arduino/uno-r4-wifi-fwuploader-plugin/branch/main/graph/badge.svg)](https://codecov.io/gh/arduino/uno-r4-wifi-fwuploader-plugin)

The `uno-r4-wifi-fwuploader-plugin` is a core component of the [arduino-fwuploader](https://github.com/arduino/arduino-fwuploader). The purpose of this plugin is to abstract all the
business logic needed to update firmware and certificates for the [uno r4 wifi](https://docs.arduino.cc/hardware/uno-r4-wifi) board.

## How to contribute

Contributions are welcome!

:sparkles: Thanks to all our [contributors](https://github.com/arduino/uno-r4-wifi-fwuploader-plugin/graphs/contributors)! :sparkles:

### Requirements

1. [Go](https://go.dev/) version 1.20 or later
1. [Task](https://taskfile.dev/) to help you run the most common tasks from the command line
1. The [uno r4 wifi](https://docs.arduino.cc/hardware/uno-r4-wifi) board to test the core parts.

## Development

When running only the plugin without the fwuploader,  the required tools are downloaded by the fwuploader. If you run only the plugin, you must provide them by hand.
Therefore be sure to place the `espflash` and `bossac` binaries in the correct folders like the following:

```bash
.
├── bossac
│   └── 1.9.1-arduino5
│       └── bossac
├── espflash
│   └── 2.0.0
│       └── espflash
└── uno-r4-wifi-fwuploader-plugin_linux_amd64
    └── bin
        └── uno-r4-wifi-fwuploader-plugin
```

**Commands**

- `uno-r4-wifi-fwuploader-plugin cert flash -p /dev/ttyACM0 ./certificate/testdata/portenta.pem`
- `uno-r4-wifi-fwuploader-plugin firmware get-version -p /dev/ttyACM0`
- `uno-r4-wifi-fwuploader-plugin firmware flash -p /dev/ttyACM0 ~/Documents/fw0.2.0.bin`


## Known problems:

### Espflash panic `UnknownModel`

On some arm64 Linux distros, version 2.0.0 of [espflash](https://github.com/esp-rs/espflash/) might panic with the following error:

```
Error:   × Main thread panicked.
  ├─▶ at espflash/src/interface.rs:70:33
  ╰─▶ called `Result::unwrap()` on an `Err` value: UnknownModel
  help: set the `RUST_BACKTRACE=1` environment variable to display a
        backtrace.
```

### The esp32 module does not go into download mode

On Linux, the uno r4 must be plugged into a **hub usb** to make the flash process work. Otherwise, it won’t be able to reboot in download mode.

```bash
$ arduino-fwuploader firmware flash -b arduino:renesas_uno:unor4wifi -a /dev/ttyACM0 -v --log-level debug

Done in 0.001 seconds
Write 46588 bytes to flash (12 pages)
[==============================] 100% (12/12 pages)
Done in 3.106 seconds

Waiting to flash the binary...
time=2023-07-18T14:50:10.492+02:00 level=INFO msg="getting firmware version"
time=2023-07-18T14:50:10.509+02:00 level=INFO msg="firmware version is > 0.1.0 using sketch"
time=2023-07-18T14:50:10.511+02:00 level=INFO msg="check if serial port has changed"
[2023-07-18T12:50:20Z INFO ] 🚀 A new version of espflash is available: v2.0.1
[2023-07-18T12:50:20Z INFO ] Serial port: '/dev/ttyACM0'
[2023-07-18T12:50:20Z INFO ] Connecting...
[2023-07-18T12:50:20Z INFO ] Unable to connect, retrying with extra delay...
[2023-07-18T12:50:21Z INFO ] Unable to connect, retrying with default delay...
[2023-07-18T12:50:21Z INFO ] Unable to connect, retrying with extra delay...
[2023-07-18T12:50:21Z INFO ] Unable to connect, retrying with default delay...
[2023-07-18T12:50:21Z INFO ] Unable to connect, retrying with extra delay...
[2023-07-18T12:50:21Z INFO ] Unable to connect, retrying with default delay...
[2023-07-18T12:50:21Z INFO ] Unable to connect, retrying with extra delay...
Error: espflash::connection_failed

  × Error while connecting to device
  ╰─▶ Failed to connect to the device
  help: Ensure that the device is connected and the reset and boot pins are
        not being held down

Error: exit status 1
ERRO[0021] couldn't update firmware: exit status 3
INFO[0021] Waiting 1 second before retrying...
INFO[0022] Uploading firmware (try 2 of 9)
time=2023-07-18T14:50:22.229+02:00 level=INFO msg=upload_command_sketch
time=2023-07-18T14:50:22.230+02:00 level=INFO msg="sending serial reset"
Error: reboot mode: upload commands sketch: setting DTR to OFF
...
```

### I flashed the certificates, but I am unable to reach the host

The **whole certificate chain** is needed to make it work. Using `-u` flags (ex: `-u www.arduino.cc:443`) won’t work because it
only downloads the root certificates. The solution is to use only the `-f` flag and provide a pem certificate containing the whole chain.

### My antivirus says that `espflash` is a threat

The binary is not signed [#348](https://github.com/esp-rs/espflash/issues/348), and some antiviruses might complain. If still doubtful, https://github.com/esp-rs/espflash is open source,
and it's possible to double-check the md5 hashes of the binary and the source code.
For more information, you can follow [this](https://forum.arduino.cc/t/radio-module-firmware-version-0-2-0-is-now-available/1147361/11) forum thread.

## Security

If you think you found a vulnerability or other security-related bug in the uno-r4-wifi-fwuploader-plugin, please read our [security
policy] and report the bug to our Security Team 🛡️ Thank you!

e-mail contact: security@arduino.cc

## License

uno-r4-wifi-fwuploader-plugin is licensed under the [AGPL 3.0](LICENSE.txt) license.

You can be released from the requirements of the above license by purchasing a commercial license. Buying such a license
is mandatory if you want to modify or otherwise use the software for commercial activities involving the Arduino
software without disclosing the source code of your own applications. To purchase a commercial license, send an email to
license@arduino.cc

