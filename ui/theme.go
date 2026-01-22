package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// DarkPurpleTheme est le thème personnalisé violet foncé
type DarkPurpleTheme struct{}

var _ fyne.Theme = (*DarkPurpleTheme)(nil)

func (t DarkPurpleTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 25, G: 20, B: 35, A: 255} // Violet très foncé
	case theme.ColorNameButton:
		return color.RGBA{R: 100, G: 70, B: 140, A: 255} // Violet moyen
	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 60, G: 50, B: 70, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 150, G: 100, B: 200, A: 255} // Violet clair
	case theme.ColorNameHover:
		return color.RGBA{R: 120, G: 90, B: 160, A: 255}
	case theme.ColorNameFocus:
		return color.RGBA{R: 170, G: 120, B: 220, A: 255}
	case theme.ColorNameForeground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // Blanc pur
	case theme.ColorNameDisabled:
		return color.RGBA{R: 150, G: 150, B: 150, A: 255} // Gris clair pour désactivé
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255} // Gris très clair
	case theme.ColorNamePressed:
		return color.RGBA{R: 130, G: 100, B: 180, A: 255}
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 80, G: 60, B: 100, A: 255}
	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 100}
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 35, G: 30, B: 50, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 130, G: 90, B: 170, A: 255}
	case theme.ColorNameSuccess:
		return color.RGBA{R: 100, G: 200, B: 150, A: 255}
	case theme.ColorNameWarning:
		return color.RGBA{R: 255, G: 180, B: 100, A: 255}
	case theme.ColorNameError:
		return color.RGBA{R: 220, G: 80, B: 100, A: 255}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t DarkPurpleTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t DarkPurpleTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t DarkPurpleTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
