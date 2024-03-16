# XPlane Serial GPS Connector

XPlane 12 does not appear to have the ability to send positions out as NMEA sentences over a serial port. This is a simple tool to provide this functionality.

This tool will locate a running X-Plane 11 or 12 on the network and send NMEA GGA and VTG sentences out over a serial port of your choice.

## Installation

Until I get around to doing a release, you will have to compile it yourself. This is a [fyne](https://fyne.io/) applications written in [go](https://go.dev/). The instructions to build fyne apps are found [here](https://docs.fyne.io/started/packaging.html).

MacOS:

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os darwin
```

Or on Windows:

```cmd
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os windows
```

## Extend

My needs are for _GGA_ and _VTG_ sentences. Yours might be for something else. If so, just create something that implements the `Outputter` interface and add it to the list of outputters in main.go.
