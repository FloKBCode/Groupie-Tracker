package ui

import (
	"context"
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ArtistListView struct {
	Container      fyne.CanvasObject
	allArtists     []models.Artist
	filteredArtists []models.Artist
	onSelectArtist func(int)

	// Moteurs de recherche et filtrage
	searchEngine *services.SearchEngine
	filterEngine *services.FilterEngine

	// Widgets
	list         *widget.List
	searchBar    *SearchBar
	statusLabel  *widget.Label
	filtersPanel *FiltersPanel
	
	// Context pour annuler les goroutines
	ctx    context.Context
	cancel context.CancelFunc
}

func NewArtistListViewWithNavigation(onSelectArtist func(int)) *ArtistListView {
	view := &ArtistListView{
		onSelectArtist: onSelectArtist,
	}

	// Créer un context avec cancel
	view.ctx, view.cancel = context.WithCancel(context.Background())

	// Chargement des artistes
	artists, err := services.GetArtists()
	if err != nil {
		view.Container = container.NewCenter(
			widget.NewLabel("❌ Erreur: " + err.Error()),
		)
		return view
	}

	view.allArtists = artists
	view.filteredArtists = artists

	// Initialiser les moteurs
	view.searchEngine = services.NewSearchEngine(artists)
	view.filterEngine = services.NewFilterEngine(artists)

	// Créer le panneau de filtres amélioré
	view.filtersPanel = NewFiltersPanel(func(criteria *services.FilterCriteria) {
		view.applyFilters(criteria)
	})

	// Pré-charger les données agrégées en arrière-plan avec context
	go view.preloadAggregates()

	// Créer les widgets
	view.buildUI()

	return view
}

// preloadAggregates charge les données agrégées en arrière-plan
func (v *ArtistListView) preloadAggregates() {
	for _, artist := range v.allArtists {
		// Vérifier si le context est annulé
		select {
		case <-v.ctx.Done():
			fmt.Println("🛑 Préchargement annulé")
			return
		default:
			v.searchEngine.LoadAggregateData(artist.ID)
			v.filterEngine.LoadAggregateData(artist.ID)
		}
	}
	
	// Une fois chargé, mettre à jour les locations disponibles
	v.filtersPanel.LoadAvailableLocations(v.filterEngine)
	
	fmt.Println("✅ Données agrégées chargées pour la recherche et les filtres")
}

// Cleanup annule les goroutines en cours
func (v *ArtistListView) Cleanup() {
	if v.cancel != nil {
		v.cancel()
		fmt.Println("🧹 Nettoyage des goroutines")
	}
}



// buildUI construit l'interface
func (v *ArtistListView) buildUI() {
	// Titre principal
	title := widget.NewLabelWithStyle(
		"🎵 Groupie Tracker",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	subtitle := widget.NewLabelWithStyle(
		"Découvrez vos artistes préférés",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	// Barre de recherche
	v.searchBar = NewSearchBar(v.searchEngine, v.onSelectArtist)

	// Liste des artistes
	v.list = widget.NewList(
		func() int {
			return len(v.filteredArtists)
		},
		func() fyne.CanvasObject {
			name := widget.NewLabel("")
			name.TextStyle = fyne.TextStyle{Bold: true}
			info := widget.NewLabel("")
			return container.NewVBox(name, info)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			artist := v.filteredArtists[id]
			vbox := obj.(*fyne.Container)
			nameLabel := vbox.Objects[0].(*widget.Label)
			infoLabel := vbox.Objects[1].(*widget.Label)

			nameLabel.SetText(artist.Name)
			infoLabel.SetText(fmt.Sprintf(
				"👥 %d membres | 📅 %d | 💿 %s",
				len(artist.Members),
				artist.CreationDate,
				artist.FirstAlbum,
			))
		},
	)

	v.list.OnSelected = func(id widget.ListItemID) {
		if id < len(v.filteredArtists) {
			artistID := v.filteredArtists[id].ID
			if v.onSelectArtist != nil {
				v.onSelectArtist(artistID)
			}
			v.list.UnselectAll()
		}
	}

	// Label de status
	v.statusLabel = widget.NewLabel(fmt.Sprintf("Affichage de %d artistes", len(v.filteredArtists)))
	v.statusLabel.Alignment = fyne.TextAlignCenter

	// Bouton pour ouvrir les filtres (fenêtre séparée)
	filterButton := widget.NewButton("🔧 Ouvrir les Filtres", func() {
		v.showFiltersWindow()
	})

	// Bouton reset
	resetButton := widget.NewButton("🔄 Réinitialiser", func() {
		v.resetFilters()
	})

	// Bouton info
	infoButton := widget.NewButton("ℹ️ Aide", func() {
		v.showHelpDialog()
	})

	// Toolbar avec tous les boutons
	toolbar := container.NewHBox(
		filterButton,
		resetButton,
		infoButton,
	)

	// Assemblage final
	content := container.NewBorder(
		// Top: Titre + Recherche + Toolbar
		container.NewVBox(
			title,
			subtitle,
			widget.NewSeparator(),
			v.searchBar.Container,
			widget.NewSeparator(),
			toolbar,
			widget.NewSeparator(),
		),
		// Bottom: Status
		v.statusLabel,
		// Left/Right: nil
		nil,
		nil,
		// Center: Liste
		v.list,
	)

	v.Container = content
}

// showFiltersWindow affiche la fenêtre de filtres (déplaçable)
func (v *ArtistListView) showFiltersWindow() {
	if v.filtersPanel == nil {
		v.filtersPanel = NewFiltersPanel(func(criteria *services.FilterCriteria) {
			v.applyFilters(criteria)
		})
		// Charger les locations une fois créé
		v.filtersPanel.LoadAvailableLocations(v.filterEngine)
	}
	v.filtersPanel.Show()
	fmt.Println("🔧 Fenêtre de filtres ouverte")
}

// applyFilters applique les critères de filtrage
func (v *ArtistListView) applyFilters(criteria *services.FilterCriteria) {
	v.filteredArtists = v.filterEngine.ApplyFilters(criteria)
	v.list.Refresh()
	v.statusLabel.SetText(fmt.Sprintf(
		"Affichage de %d artistes sur %d",
		len(v.filteredArtists),
		len(v.allArtists),
	))
	fmt.Printf("✅ Filtres appliqués: %d résultats\n", len(v.filteredArtists))
}

// resetFilters réinitialise tous les filtres
func (v *ArtistListView) resetFilters() {
	v.filteredArtists = v.allArtists
	v.list.Refresh()
	v.statusLabel.SetText(fmt.Sprintf("Affichage de %d artistes", len(v.filteredArtists)))
	v.searchBar.Clear()

	// Réinitialiser aussi le panneau de filtres
	if v.filtersPanel != nil {
		v.filtersPanel.resetFilters()
	}

	fmt.Println("🔄 Tous les filtres réinitialisés")
}

// showHelpDialog affiche une boîte de dialogue d'aide
func (v *ArtistListView) showHelpDialog() {
	helpText := `🎵 Guide d'utilisation

🔍 RECHERCHE
• Tapez dans la barre de recherche
• Suggestions automatiques
• Recherche par artiste, membre, lieu ou date

🔧 FILTRES
• Cliquez sur "Ouvrir les Filtres"
• Activez les filtres souhaités
• Ajustez les valeurs avec les sliders
• Cliquez sur "Appliquer"

📋 NAVIGATION
• Cliquez sur un artiste pour voir les détails
• Utilisez "Retour" pour revenir à la liste

🔄 RÉINITIALISER
• Annule tous les filtres et recherches`

	dialog := widget.NewPopUp(
		container.NewVBox(
			widget.NewLabelWithStyle(
				"ℹ️ Aide",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true},
			),
			widget.NewSeparator(),
			widget.NewLabel(helpText),
			widget.NewSeparator(),
			widget.NewButton("OK", func() {
				// Le dialog se fermera automatiquement
			}),
		),
		fyne.CurrentApp().Driver().AllWindows()[0].Canvas(),
	)

	dialog.Resize(fyne.NewSize(400, 400))
	dialog.Show()
}
