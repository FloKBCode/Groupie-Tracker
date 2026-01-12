package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// App contient les éléments racine de l'application
type App struct {
	FyneApp fyne.App
	Window  fyne.Window
}

// NewApp initialise l'application et la fenêtre principale
func NewApp() *App {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	w.Resize(fyne.NewSize(900, 600))
	w.CenterOnScreen()

	return &App{
		FyneApp: a,
		Window:  w,
	}
}

// Run lance l'application
func (a *App) Run() {
	a.Window.ShowAndRun()
}
