package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"groupie-tracker/models"
)

func MakeArtistList(data []models.ArtistData) fyne.CanvasObject {
	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i].Artist.Name)
		},
	)

	return container.NewBorder(
		widget.NewLabel("🎵 Artistes"),
		nil, nil, nil,
		list,
	)
}
