package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	serialv "go.bug.st/serial"

	"github.com/duncanvanzyl/xplane-serial-gps-connector/nmea"
	"github.com/duncanvanzyl/xplane-serial-gps-connector/serial"
)

// SettingsMenu returns the settings menu
func (ui *AppUI) SettingsMenu(w fyne.Window) *fyne.Menu {
	prMenu := fyne.NewMenuItem("Precision", func() {
		prGrp := widget.NewRadioGroup([]string{"Standard", "Enhanced"}, func(value string) {
			switch value {
			case "Standard":
				nmea.Formats = nmea.DEFAULTS
			case "Enhanced":
				nmea.Formats = nmea.ENHANCED
			}
			ui.Logger.Debug("Precision Changed", "precision", value)
		})
		if nmea.Formats == nmea.ENHANCED {
			prGrp.SetSelected("Enhanced")
		} else {
			prGrp.SetSelected("Standard")
		}

		info := widget.NewLabel(
			"Sets the number of decimal places in the NMEA\n" +
				"Sentences.\n" +
				"• Standard follows the NMEA standard.\n" +
				"• Enhanced adds more decimal places to the NMEA\n" +
				"  Sentences, but may not be supported by all devices.")
		// info.Alignment = fyne.TextAlignLeading

		c := container.NewVBox(
			info,
			widget.NewSeparator(),
			prGrp,
		)

		dialog.ShowCustom("Precision", "Done", c, w)

	})
	spMenu := fyne.NewMenuItem("Serial Port", func() {
		ser, ok := ui.app.Serial.(*serial.Serial)
		if !ok {
			return
		}
		mode := ser.Mode()

		dataBits := widget.NewSelect([]string{"8", "7", "6", "5"}, func(value string) {
			v, err := strconv.Atoi(value)
			if err != nil {
				ui.Logger.Error("Failed to convert DataBits", "err", err)
				return
			}
			if v < 5 || v > 8 {
				ui.Logger.Error("Invalid DataBits", "databits", v)
				return
			}
			ser.SetDataBits(v)
			ui.Logger.Debug("Set DataBits", "databits", value)
		})

		parity := widget.NewSelect([]string{"None", "Odd", "Even", "Mark", "Space"}, func(value string) {
			switch value {
			case "Odd":
				ser.SetParity(serialv.OddParity)
			case "Even":
				ser.SetParity(serialv.EvenParity)
			case "Mark":
				ser.SetParity(serialv.MarkParity)
			case "Space":
				ser.SetParity(serialv.SpaceParity)
			default:
				ser.SetParity(serialv.NoParity)
			}
			ui.Logger.Debug("Set Parity", "parity", value)
		})

		stopBits := widget.NewSelect([]string{"1", "1.5", "2"}, func(value string) {
			switch value {
			case "1.5":
				ser.SetStopBits(serialv.OnePointFiveStopBits)
			case "2":
				ser.SetStopBits(serialv.TwoStopBits)
			default:
				ser.SetStopBits(serialv.OneStopBit)
				ui.Logger.Debug("Set StopBits", "stopbits", value)
			}
		})

		info := widget.NewLabel("These are the advanced settings for the serial port.")
		info.Alignment = fyne.TextAlignCenter
		// info.Wrapping = fyne.TextWrapWord

		warning := widget.NewLabel(
			"You probably don't want to change these settings\n" +
				"unless you know what you are doing.",
		)
		warning.Alignment = fyne.TextAlignCenter
		// warning.Wrapping = fyne.TextWrapWord
		warning.Importance = widget.HighImportance

		dialog.ShowCustom("Serial Port", "Done",
			container.NewVBox(
				info,
				warning,
				widget.NewSeparator(),
				container.New(layout.NewFormLayout(),
					widget.NewLabel("Baude Rate"), ui.baudRate,
					widget.NewLabel("Data Bits"), dataBits,
					widget.NewLabel("Parity"), parity,
					widget.NewLabel("Stop Bits"), stopBits,
				),
			),
			w,
		)

		switch mode.DataBits {
		case 7:
			dataBits.SetSelected("7")
		case 6:
			dataBits.SetSelected("6")
		case 5:
			dataBits.SetSelected("5")
		default:
			dataBits.SetSelected("8")
		}

		switch mode.Parity {
		case serialv.OddParity:
			parity.SetSelected("Odd")
		case serialv.EvenParity:
			parity.SetSelected("Even")
		case serialv.MarkParity:
			parity.SetSelected("Mark")
		case serialv.SpaceParity:
			parity.SetSelected("Space")
		default:
			parity.SetSelected("None")
		}

		switch mode.StopBits {
		case serialv.OnePointFiveStopBits:
			stopBits.SetSelected("1.5")
		case serialv.TwoStopBits:
			stopBits.SetSelected("2")
		default:
			stopBits.SetSelected("1")
		}
	})
	return fyne.NewMenu("Settings",
		prMenu,
		spMenu,
	)
}
