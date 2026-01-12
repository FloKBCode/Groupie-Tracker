package models

type ArtistAggregate struct {
	Artist    Artist
	Locations Location
	Dates     ConcertDate
	Relation  Relation
}
