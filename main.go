package main

import (
	"bytes"
	"image"
	"image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"genQ/qr"
	"genQ/ui"

	qrcode "github.com/skip2/go-qrcode"
)

func main() {
	a := app.NewWithID("com.genq.GenQ")
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("genQ — QR Code Generator")
	w.Resize(fyne.NewSize(800, 600))

	preview := ui.NewPreview(400)
	var controls *ui.Controls

	regenerate := func() {
		if controls == nil {
			return
		}
		vals := controls.Values
		if vals.Text == "" {
			return
		}

		level := parseLevel(vals.ErrorCorrection)
		matrix, err := qr.GenerateMatrix(vals.Text, level)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		opts := buildRenderOpts(vals, len(matrix))
		img := qr.RenderImage(matrix, opts)
		preview.SetImage(img)
	}

	controls = ui.NewControls(w, regenerate)

	// Dark/light mode toggle
	isDark := true
	themeBtn := widget.NewButtonWithIcon("Light Mode", theme.ColorPaletteIcon(), nil)
	themeBtn.OnTapped = func() {
		if isDark {
			a.Settings().SetTheme(theme.LightTheme())
			themeBtn.SetText("Dark Mode")
			isDark = false
		} else {
			a.Settings().SetTheme(theme.DarkTheme())
			themeBtn.SetText("Light Mode")
			isDark = true
		}
	}

	// Copy to clipboard
	copyBtn := widget.NewButton("Copy to Clipboard", func() {
		img := renderCurrentImage(controls, w)
		if img == nil {
			return
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			dialog.ShowError(err, w)
			return
		}
		a.Clipboard().SetContent(buf.String())
		dialog.ShowInformation("Copied", "QR code copied to clipboard", w)
	})

	// Save as PNG
	saveBtn := widget.NewButton("Save as PNG", func() {
		img := renderCurrentImage(controls, w)
		if img == nil {
			return
		}
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if writer == nil {
				return
			}
			defer writer.Close()
			if err := png.Encode(writer, img); err != nil {
				dialog.ShowError(err, w)
			}
		}, w)
		dlg.SetFilter(&pngFilter{})
		dlg.SetFileName("qrcode.png")
		dlg.Show()
	})

	// Action buttons row
	actions := container.NewHBox(copyBtn, saveBtn, layout.NewSpacer(), themeBtn)

	// Left panel: controls with scroll
	leftPanel := container.NewVScroll(controls.Container)
	leftPanel.SetMinSize(fyne.NewSize(280, 0))

	// Right panel: preview centred
	previewCentered := container.NewCenter(preview)

	// Main layout
	content := container.NewBorder(
		nil,
		actions,
		leftPanel,
		nil,
		previewCentered,
	)

	w.SetContent(content)

	// Initial render
	regenerate()

	w.ShowAndRun()
}

func renderCurrentImage(controls *ui.Controls, w fyne.Window) *image.RGBA {
	vals := controls.Values
	if vals.Text == "" {
		return nil
	}

	level := parseLevel(vals.ErrorCorrection)
	matrix, err := qr.GenerateMatrix(vals.Text, level)
	if err != nil {
		dialog.ShowError(err, w)
		return nil
	}

	opts := buildRenderOpts(vals, len(matrix))
	return qr.RenderImage(matrix, opts)
}

func buildRenderOpts(vals ui.ControlValues, matrixSize int) qr.RenderOptions {
	pixelSize := vals.SizePreset.Pixels / matrixSize
	if pixelSize < 1 {
		pixelSize = 1
	}
	return qr.RenderOptions{
		PixelSize:     pixelSize,
		DotStyle:      parseDotStyle(vals.DotStyle),
		EyeFrameStyle: parseEyeStyle(vals.EyeStyle),
		EyePupilStyle: parseEyeStyle(vals.EyeStyle),
		Foreground:    vals.Foreground,
		Background:    vals.Background,
	}
}

func parseDotStyle(s string) qr.DotStyle {
	switch s {
	case "Circular":
		return qr.DotCircle
	case "Rounded":
		return qr.DotRounded
	default:
		return qr.DotSquare
	}
}

func parseEyeStyle(s string) qr.EyeStyle {
	switch s {
	case "Rounded":
		return qr.EyeRounded
	case "Circle":
		return qr.EyeCircle
	default:
		return qr.EyeSquare
	}
}

func parseLevel(s string) qrcode.RecoveryLevel {
	switch s {
	case "Low":
		return qrcode.Low
	case "Medium":
		return qrcode.Medium
	case "Quartile":
		return qrcode.High
	case "High":
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

type pngFilter struct{}

func (f *pngFilter) Matches(uri fyne.URI) bool {
	return uri.Extension() == ".png"
}

func (f *pngFilter) String() string {
	return "PNG Images (*.png)"
}
