package ui

import (
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ArtistDetailsView représente la vue détails d'un artiste
type ArtistDetailsView struct {
	Container fyne.CanvasObject
	aggregate models.ArtistAggregate
	onBack    func() // Callback pour retourner à la liste
}

// NewArtistDetailsView crée la vue détails (compatible avec ton code actuel)
func NewArtistDetailsView(artistID int) *ArtistDetailsView {
	return NewArtistDetailsViewWithNavigation(artistID, nil)
}

// NewArtistDetailsViewWithNavigation crée la vue détails avec navigation
func NewArtistDetailsViewWithNavigation(artistID int, onBack func()) *ArtistDetailsView {
	view := &ArtistDetailsView{
		onBack: onBack,
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
			widget.NewLabel(fmt.Sprintf("❌ Artiste avec ID %d non trouvé", artistID)),
		)
		return view
	}

	// Agréger les données (locations, dates, relations)
	aggregate, err := services.AggregateArtist(artist)
	if err != nil {
		view.Container = container.NewCenter(
			widget.NewLabel("❌ Erreur lors du chargement des détails: " + err.Error()),
		)
		return view
	}
	view.aggregate = aggregate

	// Construction de l'interface
	view.Container = view.buildUI()
	return view
}

func (v *ArtistDetailsView) buildUI() fyne.CanvasObject {
	artist := v.aggregate.Artist

	// Titre avec nom de l'artiste
	title := widget.NewLabelWithStyle(
		artist.Name,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Informations générales
	infoCard := widget.NewCard(
		"📋 Informations",
		"",
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("📅 Année de création: %d", artist.CreationDate)),
			widget.NewLabel(fmt.Sprintf("💿 Premier album: %s", artist.FirstAlbum)),
			widget.NewLabel(fmt.Sprintf("👥 Nombre de membres: %d", len(artist.Members))),
		),
	)

	// Liste des membres
	membersList := widget.NewList(
		func() int {
			return len(artist.Members)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText("• " + artist.Members[id])
		},
	)

	membersCard := widget.NewCard(
		"👥 Membres",
		"",
		container.NewVBox(membersList),
	)

	// Lieux de concert
	locationsList := widget.NewList(
		func() int {
			return len(v.aggregate.Locations.Locations)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			loc := v.aggregate.Locations.Locations[id]
			// Formater: "los_angeles-usa" -> "Los Angeles, USA"
			city, country := services.ParseLocation(loc)
			city = strings.Title(strings.ToLower(city))
			country = strings.ToUpper(country)
			
			formatted := city
			if country != "" {
				formatted = fmt.Sprintf("%s, %s", city, country)
			}
			
			obj.(*widget.Label).SetText("📍 " + formatted)
		},
	)

	locationsCard := widget.NewCard(
		"🗺️ Lieux de Concert",
		fmt.Sprintf("%d lieux", len(v.aggregate.Locations.Locations)),
		container.NewVBox(locationsList),
	)

	// Dates de concert
	datesList := widget.NewList(
		func() int {
			return len(v.aggregate.Dates.Dates)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			date := v.aggregate.Dates.Dates[id]
			obj.(*widget.Label).SetText("📅 " + date)
		},
	)

	datesCard := widget.NewCard(
		"📅 Dates de Concert",
		fmt.Sprintf("%d concerts", len(v.aggregate.Dates.Dates)),
		container.NewVBox(datesList),
	)

	// Relations (dates par location)
	relationsContent := container.NewVBox()
	for location, dates := range v.aggregate.Relation.DatesLocations {
		city, country := services.ParseLocation(location)
		city = strings.Title(strings.ToLower(city))
		country = strings.ToUpper(country)
		
		locationFormatted := city
		if country != "" {
			locationFormatted = fmt.Sprintf("%s, %s", city, country)
		}
		
		locLabel := widget.NewLabelWithStyle(
			fmt.Sprintf("📍 %s", locationFormatted),
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		)
		
		datesStr := strings.Join(dates, ", ")
		datesLabel := widget.NewLabel("   " + datesStr)
		datesLabel.Wrapping = fyne.TextWrapWord
		
		relationsContent.Add(locLabel)
		relationsContent.Add(datesLabel)
		relationsContent.Add(widget.NewSeparator())
	}

	relationsCard := widget.NewCard(
		"🎫 Concerts Programmés",
		"Dates par lieu",
		container.NewVScroll(relationsContent),
	)

	// Bouton retour
	backButton := widget.NewButton("← Retour à la liste", func() {
		if v.onBack != nil {
			v.onBack()
		} else {
			fmt.Println("⚠️ Callback de retour non défini")
		}
	})

	// Layout en colonnes
	leftColumn := container.NewVBox(
		infoCard,
		membersCard,
	)

	rightColumn := container.NewVBox(
		locationsCard,
		datesCard,
	)

	columns := container.NewGridWithColumns(
		2,
		leftColumn,
		rightColumn,
	)

	// Layout final
	content := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()),
		backButton,
		nil,
		nil,
		container.NewVScroll(
			container.NewVBox(
				columns,
				relationsCard,
			),
		),
	)

	return content
}