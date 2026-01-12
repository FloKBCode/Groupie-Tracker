package ui

import (
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ArtistListView représente la vue liste des artistes
type ArtistListView struct {
	Container     fyne.CanvasObject
	artists       []models.Artist
	onSelectArtist func(int) // Callback pour la navigation
}

// NewArtistListView crée la vue liste (compatible avec ton code actuel)
func NewArtistListView() *ArtistListView {
	return NewArtistListViewWithNavigation(nil)
}

// NewArtistListViewWithNavigation crée la vue liste avec navigation
func NewArtistListViewWithNavigation(onSelectArtist func(int)) *ArtistListView {
	view := &ArtistListView{
		onSelectArtist: onSelectArtist,
	}

	// Chargement des artistes
	artists, err := services.GetArtists()
	if err != nil {
		view.Container = container.NewCenter(
			widget.NewLabel("❌ Erreur lors du chargement des artistes: " + err.Error()),
		)
		return view
	}
	view.artists = artists

	// Création de la liste avec callback de sélection
	list := widget.NewList(
		func() int {
			return len(view.artists)
		},
		func() fyne.CanvasObject {
			// Template pour chaque ligne
			name := widget.NewLabel("")
			name.TextStyle = fyne.TextStyle{Bold: true}
			
			info := widget.NewLabel("")
			
			return container.NewVBox(
				name,
				info,
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			// Mise à jour de chaque ligne
			artist := view.artists[id]
			
			vbox := obj.(*fyne.Container)
			nameLabel := vbox.Objects[0].(*widget.Label)
			infoLabel := vbox.Objects[1].(*widget.Label)
			
			nameLabel.SetText(artist.Name)
			infoLabel.SetText(fmt.Sprintf(
				"👥 %d membres | 📅 Créé en %d | 💿 Premier album: %s",
				len(artist.Members),
				artist.CreationDate,
				artist.FirstAlbum,
			))
		},
	)

	// Gestion du clic sur un artiste
	list.OnSelected = func(id widget.ListItemID) {
		artistID := view.artists[id].ID
		fmt.Printf("✅ Artiste sélectionné: %s (ID: %d)\n", view.artists[id].Name, artistID)
		
		// Appeler le callback de navigation si défini
		if view.onSelectArtist != nil {
			view.onSelectArtist(artistID)
		}
		
		list.UnselectAll() // Désélectionner après le clic
	}

	// Titre
	title := widget.NewLabelWithStyle(
		"🎵 Groupie Tracker - Liste des Artistes",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Info footer
	footer := widget.NewLabel(fmt.Sprintf("Total: %d artistes | Cliquez sur un artiste pour voir les détails", len(artists)))
	footer.Alignment = fyne.TextAlignCenter

	// Assemblage
	content := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()),
		footer,
		nil,
		nil,
		list,
	)

	view.Container = content
	return view
}