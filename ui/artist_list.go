package ui

import (
	"context"
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type ViewMode int

const (
	ViewModeList ViewMode = iota
	ViewModeGallery
	ViewModeMap
)

type ArtistListView struct {
	Container       fyne.CanvasObject
	allArtists      []models.Artist
	filteredArtists []models.Artist
	onSelectArtist  func(int)
	onShowFavorites func()

	searchEngine     *services.SearchEngine
	filterEngine     *services.FilterEngine
	geocoder         *services.GeocodingService
	geoPreloader     *services.GeocodingPreloader
	favoritesManager *services.FavoritesManager
	imageCache       *services.ImageCache

	listView      *widget.List
	galleryView   fyne.CanvasObject
	currentView   fyne.CanvasObject
	searchBar     *SearchBar
	statusLabel   *widget.Label
	filtersPanel  *FiltersPanel
	viewContainer *fyne.Container

	viewMode ViewMode
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewArtistListView(onSelectArtist func(int), favMgr *services.FavoritesManager, imgCache *services.ImageCache, onShowFavorites func()) *ArtistListView {
	view := &ArtistListView{
		onSelectArtist:   onSelectArtist,
		favoritesManager: favMgr,
		imageCache:       imgCache,
		onShowFavorites:  onShowFavorites,
		viewMode:         ViewModeList,
	}

	view.ctx, view.cancel = context.WithCancel(context.Background())

	artists, err := services.GetArtists()
	if err != nil {
		view.Container = container.NewCenter(
			widget.NewLabel("❌ Erreur: " + err.Error()),
		)
		return view
	}

	view.allArtists = artists
	view.filteredArtists = artists

	view.searchEngine = services.NewSearchEngine(artists)
	view.filterEngine = services.NewFilterEngine(artists)
	view.geocoder = services.NewGeocodingService()
	view.geoPreloader = services.NewGeocodingPreloader(view.geocoder)

	view.filtersPanel = NewFiltersPanel(func(criteria *services.FilterCriteria) {
		view.applyFilters(criteria)
	})

	view.buildUI()
	go view.preload()

	return view
}

func (v *ArtistListView) preload() {
	for i, artist := range v.allArtists {
		select {
		case <-v.ctx.Done():
			return
		default:
			v.searchEngine.LoadAggregateData(artist.ID)
			v.filterEngine.LoadAggregateData(artist.ID)
			
			if i%(len(v.allArtists)/10+1) == 0 {
				fmt.Printf("📊 Données: %d/%d\n", i, len(v.allArtists))
			}
		}
	}
	
	v.filtersPanel.LoadAvailableLocations(v.filterEngine)
	fmt.Println("✅ Données agrégées OK")

	fmt.Println("🖼️ Préchargement des images...")
	err := v.imageCache.PreloadImages(v.allArtists, func(current, total int) {
		if current%(total/10+1) == 0 {
			fmt.Printf("🖼️ Images: %d/%d (%.0f%%)\n", current, total, float64(current)*100/float64(total))
		}
	})
	
	if err != nil {
		fmt.Printf("⚠️ Erreur images: %v\n", err)
	} else {
		fmt.Println("✅ Images préchargées OK")
	}

	fmt.Println("💡 Géolocalisation: à la demande (clic sur Carte)")
}

func (v *ArtistListView) preloadGeoOnDemand() {
	loaded, total := v.geoPreloader.GetProgress()
	if loaded == total && total > 0 {
		fmt.Println("✅ Géolocalisation déjà chargée")
		return
	}

	fmt.Println("🌍 Démarrage géolocalisation...")
	
	err := v.geoPreloader.PreloadAll(v.allArtists, func(current, total int) {
		if current%(total/10+1) == 0 {
			fmt.Printf("🌍 Géo: %d/%d (%.0f%%)\n", current, total, float64(current)*100/float64(total))
		}
	})

	if err != nil {
		fmt.Printf("⚠️ Erreur géo: %v\n", err)
	} else {
		fmt.Println("✅ Géolocalisation terminée")
	}
}

func (v *ArtistListView) Cleanup() {
	if v.cancel != nil {
		v.cancel()
	}
}

func (v *ArtistListView) buildUI() {
	title := widget.NewLabelWithStyle(
		"🎵 Groupie Tracker",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	v.searchBar = NewSearchBar(v.searchEngine, v.onSelectArtist)
	v.statusLabel = widget.NewLabel(fmt.Sprintf("📋 %d artistes", len(v.filteredArtists)))
	v.statusLabel.Alignment = fyne.TextAlignCenter

	listBtn := widget.NewButton("📋 Liste", func() { v.switchView(ViewModeList) })
	galleryBtn := widget.NewButton("🖼️ Galerie", func() { v.switchView(ViewModeGallery) })
	mapBtn := widget.NewButton("🗺️ Carte", func() { v.switchView(ViewModeMap) })
	favBtn := widget.NewButton(fmt.Sprintf("⭐ Favoris (%d)", v.favoritesManager.Count()), func() {
		if v.onShowFavorites != nil {
			v.onShowFavorites()
		}
	})
	favBtn.Importance = widget.HighImportance
	
	filterBtn := widget.NewButton("🔧 Filtres", func() { v.showFiltersWindow() })
	resetBtn := widget.NewButton("🔄 Reset", func() { v.resetFilters() })
	helpBtn := widget.NewButton("ℹ️ Aide", func() { v.showHelpDialog() })

	viewToolbar := container.NewHBox(
		widget.NewLabel("Affichage:"),
		listBtn, galleryBtn, mapBtn,
	)

	actionToolbar := container.NewHBox(favBtn, filterBtn, resetBtn, helpBtn)

	v.createAllViews()
	v.viewContainer = container.NewMax(v.currentView)

	content := container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			v.searchBar.Container,
			widget.NewSeparator(),
			viewToolbar,
			actionToolbar,
			widget.NewSeparator(),
		),
		v.statusLabel,
		nil,
		nil,
		v.viewContainer,
	)

	v.Container = content
}

func (v *ArtistListView) createAllViews() {
	v.listView = v.createListView()
	v.currentView = v.listView
}

func (v *ArtistListView) createListView() *widget.List {
	list := widget.NewList(
		func() int { return len(v.filteredArtists) },
		func() fyne.CanvasObject {
			nameLabel := widget.NewLabel("")
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
			nameLabel.Wrapping = fyne.TextWrapWord
			
			infoLabel := widget.NewLabel("")
			infoLabel.Wrapping = fyne.TextWrapWord
			
			favStar := widget.NewLabel("")
			
			content := container.NewBorder(
				nil,
				nil,
				favStar,
				nil,
				container.NewVBox(
					nameLabel,
					infoLabel,
					widget.NewSeparator(),
				),
			)
			
			return container.NewPadded(content)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(v.filteredArtists) {
				return
			}
			artist := v.filteredArtists[id]
			
			padded := obj.(*fyne.Container)
			bordered := padded.Objects[0].(*fyne.Container)
			vbox := bordered.Objects[0].(*fyne.Container)
			nameLabel := vbox.Objects[0].(*widget.Label)
			infoLabel := vbox.Objects[1].(*widget.Label)
			favStar := bordered.Objects[1].(*widget.Label)

			nameLabel.SetText(artist.Name)
			infoLabel.SetText(fmt.Sprintf(
				"👥 %d membres  •  📅 Créé en %d  •  💿 Premier album: %s",
				len(artist.Members),
				artist.CreationDate,
				services.FormatDate(artist.FirstAlbum),
			))
			
			if v.favoritesManager.IsFavorite(artist.ID) {
				favStar.SetText("⭐")
			} else {
				favStar.SetText("")
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id < len(v.filteredArtists) && v.onSelectArtist != nil {
			v.onSelectArtist(v.filteredArtists[id].ID)
			list.UnselectAll()
		}
	}

	return list
}

func (v *ArtistListView) createGalleryView() fyne.CanvasObject {
	// Grille avec 3 colonnes fixes
	cards := container.NewGridWithColumns(3)

	for i := range v.filteredArtists {
		artist := v.filteredArtists[i]

		var artistImage fyne.CanvasObject
		
		if img, ok := v.imageCache.GetImage(artist.ID); ok {
			canvasImg := canvas.NewImageFromImage(img)
			canvasImg.FillMode = canvas.ImageFillContain
			canvasImg.SetMinSize(fyne.NewSize(280, 220))
			artistImage = canvasImg
		} else if artist.Image != "" {
			uri, err := storage.ParseURI(artist.Image)
			if err == nil {
				img := canvas.NewImageFromURI(uri)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(280, 220))
				artistImage = img
			} else {
				artistImage = createImagePlaceholder(artist.Name)
			}
		} else {
			artistImage = createImagePlaceholder(artist.Name)
		}

		favBadge := ""
		if v.favoritesManager.IsFavorite(artist.ID) {
			favBadge = " ★"
		}

		infoBox := container.NewVBox(
			widget.NewLabelWithStyle(artist.Name+favBadge, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel(""),
			widget.NewLabel(fmt.Sprintf("Créé en %d", artist.CreationDate)),
			widget.NewLabel(fmt.Sprintf("%s", services.FormatDate(artist.FirstAlbum))),
			widget.NewLabel(fmt.Sprintf("%d membres", len(artist.Members))),
		)

		detailsButton := widget.NewButton("Voir détails", func() {
			artistID := artist.ID
			if v.onSelectArtist != nil {
				v.onSelectArtist(artistID)
			}
		})
		detailsButton.Importance = widget.LowImportance

		// Contenu de la card avec padding pour éviter débordement
		cardContent := container.NewVBox(
			artistImage,
			infoBox,
			container.NewPadded(detailsButton),
		)

		card := widget.NewCard("", "", cardContent)
		cards.Add(card)
	}

	// Centrer la grille horizontalement
	return container.NewVScroll(container.NewCenter(cards))
}

func createImagePlaceholder(name string) fyne.CanvasObject {
	firstChar := 'A'
	if len(name) > 0 {
		firstChar = rune(name[0])
	}
	
	colors := []color.Color{
		color.RGBA{R: 100, G: 150, B: 200, A: 255},
		color.RGBA{R: 150, G: 100, B: 200, A: 255},
		color.RGBA{R: 200, G: 100, B: 150, A: 255},
		color.RGBA{R: 100, G: 200, B: 150, A: 255},
		color.RGBA{R: 200, G: 150, B: 100, A: 255},
	}
	
	bg := canvas.NewRectangle(colors[int(firstChar)%len(colors)])
	bg.SetMinSize(fyne.NewSize(270, 220))
	
	initial := canvas.NewText(string(firstChar), color.White)
	initial.TextSize = 80
	initial.Alignment = fyne.TextAlignCenter
	
	return container.NewStack(bg, container.NewCenter(initial))
}

func (v *ArtistListView) switchView(mode ViewMode) {
	v.viewMode = mode

	switch mode {
	case ViewModeList:
		v.currentView = v.listView
		v.listView.Refresh()

	case ViewModeGallery:
		v.galleryView = v.createGalleryView()
		v.currentView = v.galleryView

	case ViewModeMap:
		// Note: Géolocalisation désactivée pour performances
		// Les coordonnées seront chargées à la demande si nécessaire
		v.createGlobalMapView()
	}

	v.viewContainer.Objects = []fyne.CanvasObject{v.currentView}
	v.viewContainer.Refresh()
}

// AMÉLIORATION MAJEURE: Grille d'artistes au lieu de liste scrollable
func (v *ArtistListView) createGlobalMapView() {
	content := container.NewVBox(
		widget.NewLabelWithStyle(
			"🗺️ Carte des Concerts",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true},
		),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Sélectionnez un artiste:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	// NOUVELLE GRILLE: 4 colonnes pour afficher beaucoup d'artistes à la fois
	artistGrid := container.NewGridWithColumns(4)

	for _, artist := range v.filteredArtists {
		artistCopy := artist // Copie pour la closure
		
		// Card avec image miniature
		var artistThumb fyne.CanvasObject
		if img, ok := v.imageCache.GetImage(artistCopy.ID); ok {
			thumbImg := canvas.NewImageFromImage(img)
			thumbImg.FillMode = canvas.ImageFillContain
			thumbImg.SetMinSize(fyne.NewSize(80, 80))
			artistThumb = thumbImg
		} else {
			artistThumb = createSmallPlaceholder(artistCopy.Name)
		}

		nameLabel := widget.NewLabelWithStyle(artistCopy.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		nameLabel.Wrapping = fyne.TextWrapWord

		viewMapBtn := widget.NewButton("🗺️ Voir carte", func() {
			v.showArtistMap(artistCopy.ID)
		})

		artistCard := widget.NewCard(
			"",
			"",
			container.NewVBox(
				container.NewCenter(artistThumb),
				nameLabel,
				viewMapBtn,
			),
		)

		artistGrid.Add(artistCard)
	}

	content.Add(artistGrid)
	v.currentView = container.NewVScroll(content)
}

// Placeholder miniature pour la grille
func createSmallPlaceholder(name string) fyne.CanvasObject {
	firstChar := 'A'
	if len(name) > 0 {
		firstChar = rune(name[0])
	}
	
	colors := []color.Color{
		color.RGBA{R: 100, G: 150, B: 200, A: 255},
		color.RGBA{R: 150, G: 100, B: 200, A: 255},
		color.RGBA{R: 200, G: 100, B: 150, A: 255},
	}
	
	bg := canvas.NewRectangle(colors[int(firstChar)%len(colors)])
	bg.SetMinSize(fyne.NewSize(80, 80))
	
	initial := canvas.NewText(string(firstChar), color.White)
	initial.TextSize = 40
	initial.Alignment = fyne.TextAlignCenter
	
	return container.NewStack(bg, container.NewCenter(initial))
}

func (v *ArtistListView) showArtistMap(artistID int) {
	var artist models.Artist
	for _, a := range v.allArtists {
		if a.ID == artistID {
			artist = a
			break
		}
	}

	aggregate, err := services.AggregateArtist(artist)
	if err != nil {
		return
	}

	mapView := NewMapView(aggregate, v.geocoder)

	backButton := widget.NewButton("← Retour", func() {
		v.switchView(ViewModeMap)
	})

	content := container.NewBorder(nil, backButton, nil, nil, mapView.Container)
	v.currentView = content
	v.viewContainer.Objects = []fyne.CanvasObject{v.currentView}
	v.viewContainer.Refresh()
}

func (v *ArtistListView) applyFilters(criteria *services.FilterCriteria) {
	v.filteredArtists = v.filterEngine.ApplyFilters(criteria)
	v.refreshCurrentView()
	v.statusLabel.SetText(fmt.Sprintf("📋 %d/%d artistes", len(v.filteredArtists), len(v.allArtists)))
}

func (v *ArtistListView) refreshCurrentView() {
	switch v.viewMode {
	case ViewModeList:
		v.listView.Refresh()
	case ViewModeGallery:
		v.switchView(ViewModeGallery)
	case ViewModeMap:
		v.switchView(ViewModeMap)
	}
}

func (v *ArtistListView) resetFilters() {
	v.filteredArtists = v.allArtists
	v.refreshCurrentView()
	v.statusLabel.SetText(fmt.Sprintf("📋 %d artistes", len(v.filteredArtists)))
	v.searchBar.Clear()
	if v.filtersPanel != nil {
		v.filtersPanel.resetFilters()
	}
}

func (v *ArtistListView) showFiltersWindow() {
	if v.filtersPanel == nil {
		v.filtersPanel = NewFiltersPanel(func(criteria *services.FilterCriteria) {
			v.applyFilters(criteria)
		})
		v.filtersPanel.LoadAvailableLocations(v.filterEngine)
	}
	v.filtersPanel.Show()
}

func (v *ArtistListView) showHelpDialog() {
	helpWindow := fyne.CurrentApp().NewWindow("ℹ️ Aide")
	helpWindow.Resize(fyne.NewSize(600, 500))
	helpWindow.CenterOnScreen()

	helpText := `🎵 GUIDE D'UTILISATION

🔍 RECHERCHE
• Tapez le nom d'un artiste, membre, lieu ou date
• Recherche par initiales: "fm" → Freddie Mercury
• Recherche floue: "qeen" → Queen

🎨 AFFICHAGE
• Liste: Vue détaillée classique avec séparateurs
• Galerie: Grille avec images préchargées
• Carte: Géolocalisation des concerts

⭐ FAVORIS
• Ajoutez des artistes favoris depuis leur page détail
• Accédez à vos favoris via le bouton "Favoris"
• Les favoris sont sauvegardés automatiquement

🔧 FILTRES
• Date de création
• Date premier album  
• Nombre de membres
• Lieux de concert

💡 ASTUCES
• Cliquez sur un artiste pour voir ses détails
• Les images se chargent en arrière-plan pour une navigation fluide
• La géolocalisation se charge à la première ouverture de carte
• Les dates sont au format JJ/MM/AAAA`

	content := container.NewVBox(
		widget.NewLabelWithStyle(
			"ℹ️ Aide",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true},
		),
		widget.NewSeparator(),
		widget.NewLabel(helpText),
		widget.NewSeparator(),
		widget.NewButton("OK", func() {
			helpWindow.Close()
		}),
	)

	helpWindow.SetContent(container.NewPadded(content))
	helpWindow.Show()
}
