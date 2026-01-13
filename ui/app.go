package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type App struct {
	FyneApp fyne.App
	Window  fyne.Window
	
	// Navigation
	currentView fyne.CanvasObject
	
	// Référence à la vue liste pour le cleanup
	listView *ArtistListView
}

func NewApp() *App {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	w.Resize(fyne.NewSize(1200, 700))
	w.CenterOnScreen()

	appInstance := &App{
		FyneApp: a,
		Window:  w,
	}

	// Hook de fermeture pour nettoyer les ressources
	w.SetOnClosed(func() {
		if appInstance.listView != nil {
			appInstance.listView.Cleanup()
		}
		// Fermer aussi le panneau de filtres s'il est ouvert
		if appInstance.listView != nil && appInstance.listView.filtersPanel != nil {
			appInstance.listView.filtersPanel.Hide()
		}
		a.Quit()
	})

	// Afficher la vue liste au démarrage
	appInstance.ShowArtistList()

	return appInstance
}

func (a *App) ShowArtistList() {
	// Nettoyer l'ancienne vue si elle existe
	if a.listView != nil {
		a.listView.Cleanup()
	}
	
	// Créer la nouvelle vue
	a.listView = NewArtistListViewWithNavigation(a.ShowArtistDetails)
	a.currentView = a.listView.Container
	a.Window.SetContent(a.currentView)
}

func (a *App) ShowArtistDetails(artistID int) {
	detailsView := NewArtistDetailsViewWithNavigation(artistID, a.ShowArtistList)
	a.currentView = detailsView.Container
	a.Window.SetContent(a.currentView)
}

func (a *App) Run() {
	a.Window.ShowAndRun()
}