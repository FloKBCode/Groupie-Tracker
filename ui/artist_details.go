package ui

import (
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"image/color"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// ArtistDetailsView représente une vue détails élégante et entièrement scrollable
type ArtistDetailsView struct {
	Container        fyne.CanvasObject
	aggregate        models.ArtistAggregate
	onBack           func()
	favoritesManager *services.FavoritesManager
	spotifyService   *services.SpotifyService
}

// NewArtistDetailsView crée une vue détails améliorée
func NewArtistDetailsView(artistID int, onBack func(), favMgr *services.FavoritesManager) *ArtistDetailsView {
	view := &ArtistDetailsView{
		onBack:           onBack,
		favoritesManager: favMgr,
		spotifyService:   services.NewSpotifyService(),
	}

	// Chargement des données
	artists, err := services.GetArtists()
	if err != nil {
		view.Container = container.NewCenter(
			widget.NewLabel("❌ Erreur: " + err.Error()),
		)
		return view
	}

	// Trouver l'artiste
	var artist models.Artist
	found := false
	for _, a := range artists {
		if a.ID == artistID {
			artist = a
			found = true
			break
		}
	}

	if !found {
		view.Container = container.NewCenter(
			widget.NewLabel(fmt.Sprintf("❌ Artiste #%d non trouvé", artistID)),
		)
		return view
	}

	// Agréger les données
	aggregate, err := services.AggregateArtist(artist)
	if err != nil {
		view.Container = container.NewCenter(
			widget.NewLabel("❌ Erreur chargement: " + err.Error()),
		)
		return view
	}
	view.aggregate = aggregate

	// Construire l'UI
	view.Container = view.buildUI()
	return view
}

func (v *ArtistDetailsView) buildUI() fyne.CanvasObject {
	

	// Bouton favori
	favBtn := v.createFavoriteButton()

	// Bouton retour
	backButton := widget.NewButton("← Retour à la liste", func() {
		if v.onBack != nil {
			v.onBack()
		}
	})
	backButton.Importance = widget.HighImportance

	// Toolbar avec retour et favori
	toolbar := container.NewHBox(
		backButton,
		widget.NewLabel(""), // Spacer
		favBtn,
	)

	// Header avec image et titre
	header := v.createHeader()

	// Section Spotify
	spotifySection := v.createSpotifySection()

	// Section Informations Générales
	infoSection := v.createInfoSection()

	// Section Membres
	membersSection := v.createMembersSection()

	// Sections Concerts - AMÉLIORÉ avec hiérarchie visuelle
	locationsSection := v.createLocationsSection()
	datesSection := v.createDatesSection()
	scheduleSection := v.createScheduleSection()

	// Layout final avec TOUT scrollable
	content := container.NewVBox(
		toolbar,
		widget.NewSeparator(),
		header,
		widget.NewSeparator(),

		// Spotify
		spotifySection,
		widget.NewSeparator(),

		// Grille 2 colonnes pour infos et membres
		container.NewGridWithColumns(2,
			infoSection,
			membersSection,
		),

		widget.NewSeparator(),

		// Titre concerts avec meilleure hiérarchie
		v.createSectionTitle("🎫 CONCERTS ET TOURNÉES"),
		widget.NewSeparator(),

		locationsSection,
		widget.NewSeparator(),

		datesSection,
		widget.NewSeparator(),

		scheduleSection,
		widget.NewSeparator(),

		// Bouton retour en bas aussi
		container.NewCenter(backButton),
	)

	// TOUT est dans un scroll pour une navigation fluide
	return container.NewVScroll(content)
}

func (v *ArtistDetailsView) createFavoriteButton() *widget.Button {
	artist := v.aggregate.Artist
	isFav := v.favoritesManager.IsFavorite(artist.ID)

	var btn *widget.Button
	btn = widget.NewButton("", func() {
		newState := v.favoritesManager.Toggle(artist.ID)
		if newState {
			btn.SetText("💛 Retirer des favoris")
		} else {
			btn.SetText("⭐ Ajouter aux favoris")
		}
	})

	if isFav {
		btn.SetText("💛 Retirer des favoris")
	} else {
		btn.SetText("⭐ Ajouter aux favoris")
	}

	btn.Importance = widget.HighImportance
	return btn
}

func (v *ArtistDetailsView) createSectionTitle(text string) fyne.CanvasObject {
	title := canvas.NewText(text, color.RGBA{R: 50, G: 50, B: 50, A: 255})
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	return container.NewPadded(title)
}


// createHeader crée l'en-tête avec image et titre
func (v *ArtistDetailsView) createHeader() fyne.CanvasObject {
	artist := v.aggregate.Artist

	// Image de l'artiste
	var artistImage fyne.CanvasObject
	if artist.Image != "" {
		uri, err := storage.ParseURI(artist.Image)
		if err == nil {
			img := canvas.NewImageFromURI(uri)
			img.FillMode = canvas.ImageFillContain
			img.SetMinSize(fyne.NewSize(300, 300))
			artistImage = img
		} else {
			artistImage = v.createPlaceholder()
		}
	} else {
		artistImage = v.createPlaceholder()
	}

	// Titre et sous-titre
	title := canvas.NewText(artist.Name, color.RGBA{R: 40, G: 40, B: 40, A: 255})
	title.TextSize = 36
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabelWithStyle(
		fmt.Sprintf("🎵 Actif depuis %d  •  💿 Premier album: %s", artist.CreationDate, services.FormatDate(artist.FirstAlbum)),
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	// Container header
	return container.NewVBox(
		container.NewCenter(artistImage),
		container.NewPadded(container.NewCenter(title)),
		container.NewCenter(subtitle),
	)
}

// createPlaceholder crée un placeholder pour l'image
func (v *ArtistDetailsView) createPlaceholder() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 100, G: 150, B: 200, A: 255})
	bg.SetMinSize(fyne.NewSize(300, 300))

	icon := canvas.NewText("🎵", color.White)
	icon.TextSize = 100
	icon.Alignment = fyne.TextAlignCenter

	return container.NewStack(bg, container.NewCenter(icon))
}

// createSpotifySection crée la section Spotify
func (v *ArtistDetailsView) createSpotifySection() fyne.CanvasObject {
	artist := v.aggregate.Artist

	spotifyURL := v.spotifyService.GetEmbedURL(artist.Name)

	spotifyBtn := widget.NewButton("🎧 Écouter sur Spotify", func() {
		// Ouvrir dans le navigateur
		fyne.CurrentApp().OpenURL(mustParseURL(spotifyURL))
	})
	spotifyBtn.Importance = widget.HighImportance

	card := widget.NewCard(
		"🎵 Spotify",
		fmt.Sprintf("Écoutez %s sur Spotify", artist.Name),
		container.NewCenter(spotifyBtn),
	)

	return container.NewPadded(card)
}

func mustParseURL(urlStr string) *url.URL {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}
	return parsed
}

// createInfoSection crée la section informations
func (v *ArtistDetailsView) createInfoSection() fyne.CanvasObject {
	artist := v.aggregate.Artist

	card := widget.NewCard(
		"📋 Informations Générales",
		"",
		container.NewVBox(
			v.createBigInfoRow("📅", "Année de création", fmt.Sprintf("%d", artist.CreationDate)),
			widget.NewSeparator(),
			v.createBigInfoRow("💿", "Premier album", services.FormatDate(artist.FirstAlbum)),
			widget.NewSeparator(),
			v.createBigInfoRow("👥", "Nombre de membres", fmt.Sprintf("%d", len(artist.Members))),
			widget.NewSeparator(),
			v.createBigInfoRow("🎤", "Concerts prévus", fmt.Sprintf("%d", len(v.aggregate.Dates.Dates))),
			widget.NewSeparator(),
			v.createBigInfoRow("🌍", "Lieux différents", fmt.Sprintf("%d", len(v.aggregate.Locations.Locations))),
		),
	)

	return container.NewPadded(card)
}

// createBigInfoRow crée une ligne d'information avec police plus grande
func (v *ArtistDetailsView) createBigInfoRow(icon, label, value string) fyne.CanvasObject {
	iconLabel := widget.NewLabelWithStyle(icon, fyne.TextAlignCenter, fyne.TextStyle{})
	iconLabel.Alignment = fyne.TextAlignCenter

	labelWidget := widget.NewLabel(label + ":")
	labelWidget.TextStyle = fyne.TextStyle{Bold: true}

	valueLabel := canvas.NewText(value, color.RGBA{R: 40, G: 100, B: 200, A: 255})
	valueLabel.TextSize = 18
	valueLabel.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewHBox(
		iconLabel,
		labelWidget,
		widget.NewLabel("→"),
		container.NewPadded(valueLabel),
	)
}

// createMembersSection crée la section membres
func (v *ArtistDetailsView) createMembersSection() fyne.CanvasObject {
	artist := v.aggregate.Artist

	membersList := container.NewVBox()
	for i, member := range artist.Members {
		memberLabel := canvas.NewText(member, color.RGBA{R: 40, G: 40, B: 40, A: 255})
		memberLabel.TextSize = 16

		memberRow := container.NewHBox(
			widget.NewLabel(fmt.Sprintf("%d.", i+1)),
			widget.NewLabel("👤"),
			memberLabel,
		)
		membersList.Add(container.NewPadded(memberRow))

		// Séparateur sauf pour le dernier
		if i < len(artist.Members)-1 {
			membersList.Add(widget.NewSeparator())
		}
	}

	card := widget.NewCard(
		"👥 Membres du Groupe",
		fmt.Sprintf("%d musiciens", len(artist.Members)),
		membersList,
	)

	return container.NewPadded(card)
}


// createLocationsSection crée la section lieux - AMÉLIORÉ
func (v *ArtistDetailsView) createLocationsSection() fyne.CanvasObject {
	titleBox := container.NewHBox(
		widget.NewLabelWithStyle("📍 Lieux de Concert", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("(%d lieux)", len(v.aggregate.Locations.Locations))),
	)

	// Grille plus spacieuse 4 colonnes
	grid := container.NewGridWithColumns(4)

	for _, location := range v.aggregate.Locations.Locations {
		city, country := services.ParseLocation(location)

		// Card avec hiérarchie visuelle claire
		cityLabel := canvas.NewText(city, color.RGBA{R: 40, G: 40, B: 40, A: 255})
		cityLabel.TextSize = 16
		cityLabel.TextStyle = fyne.TextStyle{Bold: true}
		cityLabel.Alignment = fyne.TextAlignCenter

		countryLabel := widget.NewLabelWithStyle(strings.ToUpper(country), fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

		locationCard := widget.NewCard(
			"",
			"",
			container.NewVBox(
				widget.NewLabel(""),
				container.NewCenter(cityLabel),
				container.NewCenter(countryLabel),
				widget.NewLabel(""),
			),
		)

		grid.Add(locationCard)
	}

	return container.NewVBox(
		titleBox,
		widget.NewSeparator(),
		grid,
	)
}

// createDatesSection crée la section dates - AMÉLIORÉ avec dates formatées
func (v *ArtistDetailsView) createDatesSection() fyne.CanvasObject {
	titleBox := container.NewHBox(
		widget.NewLabelWithStyle("📅 Dates des Concerts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("(%d dates)", len(v.aggregate.Dates.Dates))),
	)

	// Liste spacieuse avec padding
	datesList := container.NewVBox()

	// Formater les dates
	formattedDates := services.FormatDateList(v.aggregate.Dates.Dates)

	for i, date := range formattedDates {
		dateText := canvas.NewText(date, color.RGBA{R: 40, G: 80, B: 160, A: 255})
		dateText.TextSize = 15
		dateText.TextStyle = fyne.TextStyle{Bold: true}

		dateRow := container.NewHBox(
			widget.NewLabel(fmt.Sprintf("%d.", i+1)),
			widget.NewLabel("📅"),
			dateText,
		)

		// Padding autour de chaque date
		datesList.Add(container.NewPadded(dateRow))

		if i < len(formattedDates)-1 {
			datesList.Add(widget.NewSeparator())
		}
	}

	return container.NewVBox(
		titleBox,
		widget.NewSeparator(),
		datesList,
	)
}

// createScheduleSection crée la section programme - AMÉLIORÉ avec hiérarchie
func (v *ArtistDetailsView) createScheduleSection() fyne.CanvasObject {
	titleBox := container.NewHBox(
		widget.NewLabelWithStyle("🎫 Programme Détaillé", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("(%d concerts)", v.countTotalConcerts())),
	)

	scheduleList := container.NewVBox()

	for location, dates := range v.aggregate.Relation.DatesLocations {
		city, country := services.ParseLocation(location)

		// En-tête du lieu avec hiérarchie visuelle
		locationTitle := canvas.NewText(fmt.Sprintf("%s, %s", city, strings.ToUpper(country)), color.RGBA{R: 40, G: 40, B: 40, A: 255})
		locationTitle.TextSize = 18
		locationTitle.TextStyle = fyne.TextStyle{Bold: true}

		locationHeader := container.NewHBox(
			widget.NewLabel("📍"),
			locationTitle,
			widget.NewLabel(fmt.Sprintf("(%d concerts)", len(dates))),
		)

		// Dates pour ce lieu avec plus d'espace et dates formatées
		datesContainer := container.NewVBox()
		formattedDates := services.FormatDateList(dates)

		for i, date := range formattedDates {
			dateText := canvas.NewText(date, color.RGBA{R: 60, G: 100, B: 180, A: 255})
			dateText.TextSize = 14

			dateRow := container.NewHBox(
				widget.NewLabel("      →"), // Indentation visuelle
				widget.NewLabel("📅"),
				dateText,
			)
			datesContainer.Add(container.NewPadded(dateRow))

			if i < len(formattedDates)-1 {
				datesContainer.Add(widget.NewSeparator())
			}
		}

		// Card pour le lieu
		locationCard := widget.NewCard(
			"",
			"",
			container.NewVBox(
				locationHeader,
				widget.NewSeparator(),
				datesContainer,
			),
		)

		scheduleList.Add(locationCard)
		scheduleList.Add(widget.NewLabel("")) // Espace entre les lieux
	}

	return container.NewVBox(
		titleBox,
		widget.NewSeparator(),
		scheduleList,
	)
}

func (v *ArtistDetailsView) countTotalConcerts() int {
	total := 0
	for _, dates := range v.aggregate.Relation.DatesLocations {
		total += len(dates)
	}
	return total
}
