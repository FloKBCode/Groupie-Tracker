package main

import (
	"log"

	"groupie-tracker/services"
	"groupie-tracker/ui"
)

func main() {
	data, err := services.GetAllArtistData()
	if err != nil {
		log.Fatal(err)
	}

	ui.StartApp(data)
}
