package ui

import (
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type App struct {
	FyneApp fyne.App
	Window  fyne.Window
	
	// Navigation
	currentView fyne.CanvasObject
	
	// Managers
	favoritesManager *services.FavoritesManager
	imageCache       *services.ImageCache
	
	// Référence à la vue liste pour le cleanup
	listView *ArtistListView
}

func NewApp() *App {
	a := app.New()
	
	// Appliquer le thème violet foncé
	a.Settings().SetTheme(&DarkPurpleTheme{})
	
	w := a.NewWindow("Groupie Tracker")

	w.Resize(fyne.NewSize(1200, 700))
	w.CenterOnScreen()

	appInstance := &App{
		FyneApp:          a,
		Window:           w,
		favoritesManager: services.NewFavoritesManager(),
		imageCache:       services.NewImageCache(),
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

	// Afficher l'écran d'accueil avec bouton
	splash := NewSplashScreen(func() {
		// Quand l'utilisateur clique sur "Ouvrir"
		appInstance.ShowArtistList()
	})
	w.SetContent(splash.Container)

	return appInstance
}

func (a *App) ShowArtistList() {
	// Nettoyer l'ancienne vue si elle existe
	if a.listView != nil {
		a.listView.Cleanup()
	}
	
	// Créer la nouvelle vue
	a.listView = NewArtistListView(a.ShowArtistDetails, a.favoritesManager, a.imageCache, a.ShowFavorites)
	a.currentView = a.listView.Container
	a.Window.SetContent(a.currentView)
}

func (a *App) ShowArtistDetails(artistID int) {
	detailsView := NewArtistDetailsView(artistID, a.ShowArtistList, a.favoritesManager)
	a.currentView = detailsView.Container
	a.Window.SetContent(a.currentView)
}

func (a *App) ShowFavorites() {
	if a.listView == nil {
		return
	}
	
	favView := NewFavoritesView(
		a.listView.allArtists,
		a.favoritesManager,
		a.imageCache,
		a.ShowArtistDetails,
		a.ShowArtistList,
	)
	a.currentView = favView.Container
	a.Window.SetContent(a.currentView)
}

func (a *App) Run() {
	a.Window.ShowAndRun()
}
