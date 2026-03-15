package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type SizePreset struct {
	Label  string
	Pixels int
}

var ScreenSizes = []SizePreset{
	{"Small (200px)", 200},
	{"Medium (400px)", 400},
	{"Large (600px)", 600},
	{"XL (900px)", 900},
	{"XXL (1200px)", 1200},
}

var PrintSizes = []SizePreset{
	{"1″ @ 300 DPI (300px)", 300},
	{"2″ @ 300 DPI (600px)", 600},
	{"3″ @ 300 DPI (900px)", 900},
	{"4″ @ 300 DPI (1200px)", 1200},
	{"5″ @ 300 DPI (1500px)", 1500},
	{"6″ @ 300 DPI (1800px)", 1800},
}

type ControlValues struct {
	Text            string
	DotStyle        string
	EyeStyle        string
	Foreground      color.Color
	Background      color.Color
	SizeCategory    string
	SizePreset      SizePreset
	ErrorCorrection string
}

type Controls struct {
	Container *fyne.Container
	Values    ControlValues
	OnChanged func()

	fgRect     *ColorButton
	bgRect     *ColorButton
	sizeSelect *widget.Select
}

type ColorButton struct {
	widget.BaseWidget
	Color     color.Color
	swatch    *canvas.Circle
	swatchBox *fyne.Container
	btn       *widget.Button
	window    fyne.Window
	onChanged func()
}

func NewColorButton(label string, initial color.Color, win fyne.Window, onChanged func()) *ColorButton {
	cb := &ColorButton{
		Color:     initial,
		window:    win,
		onChanged: onChanged,
	}
	cb.swatch = canvas.NewCircle(initial)
	cb.swatch.StrokeColor = color.Gray{Y: 128}
	cb.swatch.StrokeWidth = 1
	cb.swatchBox = container.NewGridWrap(fyne.NewSize(24, 24), cb.swatch)
	cb.btn = widget.NewButton(label, func() {
		picker := dialog.NewColorPicker("Choose "+label, "", func(c color.Color) {
			cb.Color = c
			cb.swatch.FillColor = c
			cb.swatch.Refresh()
			if cb.onChanged != nil {
				cb.onChanged()
			}
		}, win)
		picker.Advanced = true
		picker.SetColor(cb.Color)
		picker.Show()
	})
	return cb
}

func NewControls(win fyne.Window, onChanged func()) *Controls {
	c := &Controls{
		Values: ControlValues{
			Text:            "https://example.com",
			DotStyle:        "Square",
			EyeStyle:        "Square",
			Foreground:      color.Black,
			Background:      color.White,
			SizeCategory:    "Screen",
			SizePreset:      ScreenSizes[1],
			ErrorCorrection: "Medium",
		},
		OnChanged: onChanged,
	}

	// Text input
	textEntry := widget.NewEntry()
	textEntry.SetPlaceHolder("Enter URL or text...")
	textEntry.SetText(c.Values.Text)
	textEntry.OnChanged = func(s string) {
		c.Values.Text = s
		c.fireChanged()
	}

	// Dot style
	dotStyle := widget.NewSelect([]string{"Square", "Circular", "Rounded"}, func(s string) {
		c.Values.DotStyle = s
		c.fireChanged()
	})
	dotStyle.SetSelected("Square")

	// Eye style
	eyeStyle := widget.NewSelect([]string{"Square", "Rounded", "Circle"}, func(s string) {
		c.Values.EyeStyle = s
		c.fireChanged()
	})
	eyeStyle.SetSelected("Square")

	// Colour buttons
	c.fgRect = NewColorButton("Foreground", color.Black, win, func() {
		c.Values.Foreground = c.fgRect.Color
		c.fireChanged()
	})
	c.bgRect = NewColorButton("Background", color.White, win, func() {
		c.Values.Background = c.bgRect.Color
		c.fireChanged()
	})

	// Size category
	sizeLabels := getSizeLabels(ScreenSizes)
	c.sizeSelect = widget.NewSelect(sizeLabels, func(s string) {
		c.Values.SizePreset = findPreset(s, c.Values.SizeCategory)
		c.fireChanged()
	})
	c.sizeSelect.SetSelected(sizeLabels[1])

	sizeCategory := widget.NewRadioGroup([]string{"Screen", "Print"}, func(s string) {
		c.Values.SizeCategory = s
		var presets []SizePreset
		if s == "Print" {
			presets = PrintSizes
		} else {
			presets = ScreenSizes
		}
		labels := getSizeLabels(presets)
		c.sizeSelect.Options = labels
		c.sizeSelect.SetSelected(labels[1])
		c.Values.SizePreset = presets[1]
		c.fireChanged()
	})
	sizeCategory.SetSelected("Screen")
	sizeCategory.Horizontal = true

	// Error correction
	ecLevel := widget.NewSelect([]string{"Low", "Medium", "Quartile", "High"}, func(s string) {
		c.Values.ErrorCorrection = s
		c.fireChanged()
	})
	ecLevel.SetSelected("Medium")

	// Layout
	c.Container = container.NewVBox(
		widget.NewLabel("Content"),
		textEntry,
		widget.NewSeparator(),
		widget.NewLabel("Data Dot Style"),
		dotStyle,
		widget.NewSeparator(),
		widget.NewLabel("Eye Style"),
		eyeStyle,
		widget.NewSeparator(),
		widget.NewLabel("Colours"),
		container.NewHBox(c.fgRect.swatchBox, c.fgRect.btn, c.bgRect.swatchBox, c.bgRect.btn),
		widget.NewSeparator(),
		widget.NewLabel("Output Size"),
		sizeCategory,
		c.sizeSelect,
		widget.NewSeparator(),
		widget.NewLabel("Error Correction"),
		ecLevel,
		layout.NewSpacer(),
	)

	return c
}

func (c *Controls) fireChanged() {
	if c.OnChanged != nil {
		c.OnChanged()
	}
}

func getSizeLabels(presets []SizePreset) []string {
	labels := make([]string, len(presets))
	for i, p := range presets {
		labels[i] = p.Label
	}
	return labels
}

func findPreset(label string, category string) SizePreset {
	var presets []SizePreset
	if category == "Print" {
		presets = PrintSizes
	} else {
		presets = ScreenSizes
	}
	for _, p := range presets {
		if p.Label == label {
			return p
		}
	}
	return presets[1]
}
