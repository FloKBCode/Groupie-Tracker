package main

import (
    "fmt"
    "log"
    "groupie-tracker/services"
)

func main() {
	data, err := services.GetAllArtistData()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range data {
		fmt.Println(d.Artist.Name)
		fmt.Println("Locations:", d.Location.Locations)
		fmt.Println("Dates:", d.Date.Dates)
		fmt.Println("Relations:", d.Relation.DatesLocations)
		fmt.Println("------")
	}
}
