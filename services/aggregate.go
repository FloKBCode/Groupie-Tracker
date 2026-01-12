package services

import "groupie-tracker/models"

func AggregateArtist(artist models.Artist) (models.ArtistAggregate, error) {

	locations, err := GetLocation(artist.ID)
	if err != nil {
		return models.ArtistAggregate{}, err
	}

	dates, err := GetDate(artist.ID)
	if err != nil {
		return models.ArtistAggregate{}, err
	}

	relation, err := GetRelation(artist.ID)
	if err != nil {
		return models.ArtistAggregate{}, err
	}

	return models.ArtistAggregate{
		Artist:    artist,
		Locations: locations,
		Dates:     dates,
		Relation:  relation,
	}, nil
}
