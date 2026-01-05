package services

import (
	"groupie-tracker/models"
)

// GetAllArtistData récupère toutes les données regroupées par artiste
func GetAllArtistData() ([]models.ArtistData, error) {
	artists, err := GetArtists()
	if err != nil {
		return nil, err
	}

	var result []models.ArtistData

	for _, artist := range artists {
		location, err := GetLocations(artist.ID)
		if err != nil {
			return nil, err
		}

		date, err := GetDates(artist.ID)
		if err != nil {
			return nil, err
		}

		relation, err := GetRelation(artist.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, models.ArtistData{
			Artist:   artist,
			Location: location,
			Date:     date,
			Relation: relation,
		})
	}

	return result, nil
}
