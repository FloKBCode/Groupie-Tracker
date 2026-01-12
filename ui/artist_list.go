package ui

import (
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ArtistListView représente la vue liste des artistes
type ArtistListView struct {
	Container fyne.CanvasObject
}

// NewArtistListView crée la vue liste
func NewArtistListView() *ArtistListView {

	artists, err := services.GetArtists()
	if err != nil {
		return &ArtistListView{
			Container: widget.NewLabel("Erreur lors du chargement des artistes"),
		}
	}

	list := widget.NewList(
		func() int {
			return len(artists)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(artists[id].Name)
		},
	)

	title := widget.NewLabelWithStyle(
		"Artistes",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	content := container.NewBorder(
		title,
		nil,
		nil,
		nil,
		list,
	)

	return &ArtistListView{
		Container: content,
	}
}
