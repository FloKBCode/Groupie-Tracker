package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SplashScreen représente l'écran d'accueil
type SplashScreen struct {
	Container fyne.CanvasObject
	onStart   func()
}

// NewSplashScreen crée un écran d'accueil
func NewSplashScreen(onStart func()) *SplashScreen {
	ss := &SplashScreen{
		onStart: onStart,
	}
	ss.Container = ss.buildUI()
	return ss
}

func (ss *SplashScreen) buildUI() fyne.CanvasObject {
	// Fond violet foncé
	bg := canvas.NewRectangle(color.RGBA{R: 25, G: 20, B: 35, A: 255})
	bg.SetMinSize(fyne.NewSize(800, 600))

	// Logo/Titre
	title := canvas.NewText("Groupie Tracker", color.RGBA{R: 200, G: 150, B: 255, A: 255})
	title.TextSize = 52
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// Sous-titre
	subtitle := canvas.NewText("Explorez vos artistes musicaux préférés", color.RGBA{R: 180, G: 160, B: 200, A: 255})
	subtitle.TextSize = 18
	subtitle.Alignment = fyne.TextAlignCenter

	// Description
	desc1 := widget.NewLabel("Découvrez les informations sur 52 artistes")
	desc1.Alignment = fyne.TextAlignCenter
	
	desc2 := widget.NewLabel("Consultez leurs dates de concert et lieux")
	desc2.Alignment = fyne.TextAlignCenter
	
	desc3 := widget.NewLabel("Explorez une carte interactive des tournées")
	desc3.Alignment = fyne.TextAlignCenter

	// Bouton de démarrage
	startBtn := widget.NewButton("Ouvrir l'application", func() {
		if ss.onStart != nil {
			ss.onStart()
		}
	})
	startBtn.Importance = widget.HighImportance

	// Version
	version := widget.NewLabel("Version 1.0.0")
	version.Alignment = fyne.TextAlignCenter

	// Layout
	content := container.NewVBox(
		widget.NewLabel(""),
		widget.NewLabel(""),
		widget.NewLabel(""),
		container.NewCenter(title),
		widget.NewLabel(""),
		container.NewCenter(subtitle),
		widget.NewLabel(""),
		widget.NewLabel(""),
		desc1,
		desc2,
		desc3,
		widget.NewLabel(""),
		widget.NewLabel(""),
		container.NewCenter(startBtn),
		widget.NewLabel(""),
		widget.NewLabel(""),
		widget.NewLabel(""),
		container.NewCenter(version),
	)

	return container.NewStack(bg, content)
}

// LoadingOverlay crée un overlay de chargement
func NewLoadingOverlay(message string) fyne.CanvasObject {
	// Fond semi-transparent
	bg := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 150})
	bg.SetMinSize(fyne.NewSize(800, 600))

	// Barre de progression
	progress := widget.NewProgressBarInfinite()

	// Message
	msg := widget.NewLabel(message)
	msg.Alignment = fyne.TextAlignCenter
	msg.TextStyle = fyne.TextStyle{Bold: true}

	// Card de chargement
	loadingCard := widget.NewCard(
		"",
		"",
		container.NewVBox(
			progress,
			widget.NewLabel(""),
			msg,
		),
	)

	return container.NewStack(
		bg,
		container.NewCenter(
			container.NewPadded(loadingCard),
		),
	)
}
