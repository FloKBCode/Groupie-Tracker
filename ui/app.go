package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"groupie-tracker/models"
)

func StartApp(data []models.ArtistData) {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	w.SetContent(MakeArtistList(data))
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}