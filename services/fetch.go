package services

import (
		"fmt"
    "groupie-tracker/api"
    "groupie-tracker/models"
)

const API_BASE = "https://groupietrackers.herokuapp.com/api/"

func GetArtists() ([]models.Artist, error) {
    var artists []models.Artist
    err := api.FetchJSON(API_BASE+"artists", &artists)
    return artists, err
}

func GetLocations(id int) (models.Location, error) {
    var loc models.Location
    url := fmt.Sprintf("%slocations/%d", API_BASE, id)
    err := api.FetchJSON(url, &loc)
    return loc, err
}

func GetDates(id int) (models.Date, error) {
    var date models.Date
    url := fmt.Sprintf("%s/dates/%d", API_BASE, id)
    err := api.FetchJSON(url, &date)
    return date, err
}

func GetRelation(id int) (models.Relation, error) {
	var relation models.Relation
	url := fmt.Sprintf("%s/relation/%d", API_BASE, id)
	err := api.FetchJSON(url, &relation)
	return relation, err
}
