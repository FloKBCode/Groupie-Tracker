package ui

import (
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

// FavoritesView affiche les artistes favoris
type FavoritesView struct {
	Container        fyne.CanvasObject
	allArtists       []models.Artist
	favoritesManager *services.FavoritesManager
	imageCache       *services.ImageCache
	onSelectArtist   func(int)
	onBack           func()
}

// NewFavoritesView crée une nouvelle vue favoris
func NewFavoritesView(allArtists []models.Artist, favMgr *services.FavoritesManager, imgCache *services.ImageCache, onSelectArtist func(int), onBack func()) *FavoritesView {
	view := &FavoritesView{
		allArtists:       allArtists,
		favoritesManager: favMgr,
		imageCache:       imgCache,
		onSelectArtist:   onSelectArtist,
		onBack:           onBack,
	}
	
	view.Container = view.buildUI()
	return view
}

func (v *FavoritesView) buildUI() fyne.CanvasObject {
	// En-tête
	title := widget.NewLabelWithStyle(
		"⭐ Mes Artistes Favoris",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	
	count := widget.NewLabel(fmt.Sprintf("💛 %d favoris", v.favoritesManager.Count()))
	count.Alignment = fyne.TextAlignCenter
	
	backBtn := widget.NewButton("← Retour", func() {
		if v.onBack != nil {
			v.onBack()
		}
	})
	
	clearBtn := widget.NewButton("🗑️ Tout supprimer", func() {
		v.showClearConfirmation()
	})
	
	// Obtenir les artistes favoris
	favoriteArtists := v.getFavoriteArtists()
	
	var content fyne.CanvasObject
	
	if len(favoriteArtists) == 0 {
		// Aucun favori
		emptyMsg := container.NewVBox(
			widget.NewLabel(""),
			widget.NewLabel(""),
			widget.NewLabelWithStyle(
				"💔 Aucun artiste favori",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true},
			),
			widget.NewLabel(""),
			widget.NewLabel("Ajoutez des artistes à vos favoris en cliquant sur l'étoile ⭐"),
			widget.NewLabel(""),
			widget.NewLabel(""),
		)
		
		content = container.NewVBox(
			backBtn,
			widget.NewSeparator(),
			title,
			count,
			widget.NewSeparator(),
			container.NewCenter(emptyMsg),
		)
	} else {
		// Afficher les favoris en grille
		gallery := v.createFavoritesGallery(favoriteArtists)
		
		toolbar := container.NewHBox(
			backBtn,
			clearBtn,
		)
		
		content = container.NewBorder(
			container.NewVBox(
				toolbar,
				widget.NewSeparator(),
				title,
				count,
				widget.NewSeparator(),
			),
			nil,
			nil,
			nil,
			gallery,
		)
	}
	
	return content
}

func (v *FavoritesView) getFavoriteArtists() []models.Artist {
	favoriteIDs := v.favoritesManager.GetFavorites()
	
	favorites := make([]models.Artist, 0, len(favoriteIDs))
	for _, artist := range v.allArtists {
		for _, favID := range favoriteIDs {
			if artist.ID == favID {
				favorites = append(favorites, artist)
				break
			}
		}
	}
	
	return favorites
}

func (v *FavoritesView) createFavoritesGallery(artists []models.Artist) fyne.CanvasObject {
	cards := container.NewGridWrap(fyne.NewSize(300, 420))
	
	for i := range artists {
		artist := artists[i]
		
		var artistImage fyne.CanvasObject
		
		// Essayer de charger depuis le cache
		if img, ok := v.imageCache.GetImage(artist.ID); ok {
			canvasImg := canvas.NewImageFromImage(img)
			canvasImg.FillMode = canvas.ImageFillContain
			canvasImg.SetMinSize(fyne.NewSize(270, 220))
			artistImage = canvasImg
		} else if artist.Image != "" {
			// Charger depuis l'URL
			uri, err := storage.ParseURI(artist.Image)
			if err == nil {
				img := canvas.NewImageFromURI(uri)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(270, 220))
				artistImage = img
			} else {
				artistImage = createFavImagePlaceholder(artist.Name)
			}
		} else {
			artistImage = createFavImagePlaceholder(artist.Name)
		}
		
		// Informations
		nameLabel := widget.NewLabelWithStyle(
			artist.Name,
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true},
		)
		
		infoBox := container.NewVBox(
			nameLabel,
			widget.NewLabel(fmt.Sprintf("📅 Créé en %d", artist.CreationDate)),
			widget.NewLabel(fmt.Sprintf("💿 Premier album: %s", services.FormatDate(artist.FirstAlbum))),
			widget.NewLabel(fmt.Sprintf("👥 %d membres", len(artist.Members))),
		)
		
		// Boutons
		detailsBtn := widget.NewButton("📋 Voir détails", func() {
			artistID := artist.ID
			if v.onSelectArtist != nil {
				v.onSelectArtist(artistID)
			}
		})
		
		removeBtn := widget.NewButton("💔 Retirer", func() {
			artistID := artist.ID
			v.favoritesManager.Toggle(artistID)
			v.Container = v.buildUI() // Rafraîchir
			v.onBack()                 // Retourner pour rafraîchir
		})
		removeBtn.Importance = widget.DangerImportance
		
		buttonsBox := container.NewGridWithColumns(2,
			detailsBtn,
			removeBtn,
		)
		
		// Contenu de la card avec boutons centrés
		cardContent := container.NewVBox(
			artistImage,
			infoBox,
			container.NewPadded(buttonsBox),
		)
		
		card := widget.NewCard("", "", cardContent)
		
		cards.Add(card)
	}
	
	return container.NewVScroll(cards)
}

func createFavImagePlaceholder(name string) fyne.CanvasObject {
	firstChar := 'A'
	if len(name) > 0 {
		firstChar = rune(name[0])
	}
	
	colors := []color.Color{
		color.RGBA{R: 255, G: 193, B: 7, A: 255},   // Gold
		color.RGBA{R: 233, G: 30, B: 99, A: 255},   // Pink
		color.RGBA{R: 156, G: 39, B: 176, A: 255},  // Purple
		color.RGBA{R: 63, G: 81, B: 181, A: 255},   // Indigo
		color.RGBA{R: 0, G: 150, B: 136, A: 255},   // Teal
	}
	
	bg := canvas.NewRectangle(colors[int(firstChar)%len(colors)])
	bg.SetMinSize(fyne.NewSize(270, 220))
	
	initial := canvas.NewText(string(firstChar), color.White)
	initial.TextSize = 100
	initial.Alignment = fyne.TextAlignCenter
	
	star := canvas.NewText("⭐", color.White)
	star.TextSize = 40
	star.Alignment = fyne.TextAlignCenter
	
	return container.NewStack(
		bg,
		container.NewVBox(
			container.NewCenter(initial),
			container.NewCenter(star),
		),
	)
}

func (v *FavoritesView) showClearConfirmation() {
	var dialog *widget.PopUp
	
	dialog = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabelWithStyle(
				"⚠️ Confirmation",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true},
			),
			widget.NewLabel(""),
			widget.NewLabel("Êtes-vous sûr de vouloir supprimer tous vos favoris ?"),
			widget.NewLabel("Cette action est irréversible."),
			widget.NewLabel(""),
			container.NewGridWithColumns(2,
				widget.NewButton("Annuler", func() {
					dialog.Hide()
				}),
				widget.NewButton("Supprimer tout", func() {
					v.favoritesManager.Clear()
					dialog.Hide()
					v.onBack()
				}),
			),
		),
		fyne.CurrentApp().Driver().AllWindows()[0].Canvas(),
	)
	
	dialog.Resize(fyne.NewSize(400, 200))
	dialog.Show()
}
