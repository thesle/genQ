package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type ThemeMode int

const (
	ThemeLight ThemeMode = iota
	ThemeDark
)

func ApplyTheme(app fyne.App, mode ThemeMode) {
	switch mode {
	case ThemeLight:
		app.Settings().SetTheme(theme.LightTheme())
	case ThemeDark:
		app.Settings().SetTheme(theme.DarkTheme())
	}
}
