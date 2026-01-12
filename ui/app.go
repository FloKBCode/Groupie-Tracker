package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// App contient les éléments racine de l'application
type App struct {
	FyneApp fyne.App
	Window  fyne.Window
	
	// Navigation
	currentView fyne.CanvasObject
}

// NewApp initialise l'application et la fenêtre principale
func NewApp() *App {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	w.Resize(fyne.NewSize(1200, 700))
	w.CenterOnScreen()

	appInstance := &App{
		FyneApp: a,
		Window:  w,
	}

	// Afficher la vue liste au démarrage
	appInstance.ShowArtistList()

	return appInstance
}

// ShowArtistList affiche la vue liste des artistes
func (a *App) ShowArtistList() {
	// Passer les callbacks de navigation
	listView := NewArtistListViewWithNavigation(a.ShowArtistDetails)
	a.currentView = listView.Container
	a.Window.SetContent(a.currentView)
}

// ShowArtistDetails affiche les détails d'un artiste
func (a *App) ShowArtistDetails(artistID int) {
	// Passer le callback de retour
	detailsView := NewArtistDetailsViewWithNavigation(artistID, a.ShowArtistList)
	a.currentView = detailsView.Container
	a.Window.SetContent(a.currentView)
}

// Run lance l'application
func (a *App) Run() {
	a.Window.ShowAndRun()
}