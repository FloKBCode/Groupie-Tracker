package ui

import (
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ArtistListView - Version améliorée avec recherche et filtres
type ArtistListView struct {
	Container      fyne.CanvasObject
	allArtists     []models.Artist      
	filteredArtists []models.Artist     
	onSelectArtist func(int)
	
	// Moteurs de recherche et filtrage
	searchEngine *services.SearchEngine
	filterEngine *services.FilterEngine
	
	// Widgets
	list      *widget.List
	searchBar *SearchBar
	statusLabel *widget.Label
	filtersPanel *FiltersPanel
}

// NewArtistListView crée la vue liste avec recherche et filtres
func NewArtistListViewWithNavigation(onSelectArtist func(int)) *ArtistListView {
	view := &ArtistListView{
		onSelectArtist: onSelectArtist,
	}

	// Chargement des artistes
	artists, err := services.GetArtists()
	if err != nil {
		view.Container = container.NewCenter(
			widget.NewLabel("❌ Erreur: " + err.Error()),
		)
		return view
	}
	
	view.allArtists = artists
	view.filteredArtists = artists // Au départ, tous les artistes sont affichés

	// Initialiser les moteurs
	view.searchEngine = services.NewSearchEngine(artists)
	view.filterEngine = services.NewFilterEngine(artists)

	// Créer le panneau de filtres avec callback
	view.filtersPanel = NewFiltersPanel(func(criteria *services.FilterCriteria) {
		view.applyFilters(criteria)
	})


	// Pré-charger les données agrégées pour la recherche de locations (optionnel)
	// Tu peux le faire en background ou à la demande
	go view.preloadAggregates()

	// Créer les widgets
	view.buildUI()

	return view
}

// preloadAggregates charge les données agrégées en arrière-plan
func (v *ArtistListView) preloadAggregates() {
	for _, artist := range v.allArtists {
		v.searchEngine.LoadAggregateData(artist.ID)
		v.filterEngine.LoadAggregateData(artist.ID)
	}
	fmt.Println("✅ Données agrégées chargées pour la recherche")
}

// buildUI construit l'interface
func (v *ArtistListView) buildUI() {
	// Titre
	title := widget.NewLabelWithStyle(
		"🎵 Groupie Tracker",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
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

	// Bouton pour ouvrir le panneau de filtres (à implémenter)
	filterButton := widget.NewButton("🔧 Filtres", func() {
		v.showFiltersDialog()
	})

	// Bouton reset
	resetButton := widget.NewButton("🔄 Réinitialiser", func() {
		v.resetFilters()
	})

	// Toolbar
	toolbar := container.NewHBox(
		filterButton,
		resetButton,
	)

	// Assemblage
	content := container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			v.searchBar.Container,
			toolbar,
			widget.NewSeparator(),
		),
		v.statusLabel,
		nil,
		nil,
		v.list,
	)

	v.Container = content
}

// showFiltersDialog affiche une popup avec les options de filtrage
func (v *ArtistListView) showFiltersDialog() {
	if v.filtersPanel == nil {
		return
	}

	// Affiche le panel dans une popup
	w := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog := widget.NewModalPopUp(v.filtersPanel.Container, w.Canvas())
	dialog.Show()
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
	fmt.Println("🔄 Filtres réinitialisés")
}


