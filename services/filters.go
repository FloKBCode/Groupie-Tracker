package services

import (
	"groupie-tracker/models"
	"strings"
)

// FilterCriteria contient tous les critères de filtrage possibles
type FilterCriteria struct {
	// Filtres de range
	CreationDateMin int // Année minimum de création
	CreationDateMax int // Année maximum de création
	
	FirstAlbumYearMin int // Année minimum du premier album
	FirstAlbumYearMax int // Année maximum du premier album
	
	MembersMin int // Nombre minimum de membres
	MembersMax int // Nombre maximum de membres
	
	// Filtres par location (liste de pays/villes sélectionnés)
	Locations []string // Si vide, tous les lieux sont acceptés
	
	// Filtres booléens
	EnableCreationDateFilter  bool
	EnableFirstAlbumFilter    bool
	EnableMembersFilter       bool
	EnableLocationsFilter     bool
}

// NewFilterCriteria crée des critères de filtrage par défaut (tous désactivés)
func NewFilterCriteria() *FilterCriteria {
	return &FilterCriteria{
		CreationDateMin:   1900,
		CreationDateMax:   2025,
		FirstAlbumYearMin: 1900,
		FirstAlbumYearMax: 2025,
		MembersMin:        1,
		MembersMax:        10,
		Locations:         []string{},
		
		EnableCreationDateFilter:  false,
		EnableFirstAlbumFilter:    false,
		EnableMembersFilter:       false,
		EnableLocationsFilter:     false,
	}
}

// FilterEngine gère le filtrage des artistes
type FilterEngine struct {
	artists    []models.Artist
	aggregates map[int]models.ArtistAggregate
}

// NewFilterEngine crée une nouvelle instance du moteur de filtrage
func NewFilterEngine(artists []models.Artist) *FilterEngine {
	return &FilterEngine{
		artists:    artists,
		aggregates: make(map[int]models.ArtistAggregate),
	}
}

// LoadAggregateData charge les données agrégées pour un artiste
func (fe *FilterEngine) LoadAggregateData(artistID int) error {
	if _, exists := fe.aggregates[artistID]; exists {
		return nil
	}

	var artist models.Artist
	for _, a := range fe.artists {
		if a.ID == artistID {
			artist = a
			break
		}
	}

	aggregate, err := AggregateArtist(artist)
	if err != nil {
		return err
	}

	fe.aggregates[artistID] = aggregate
	return nil
}

// ApplyFilters applique les critères de filtrage aux artistes
func (fe *FilterEngine) ApplyFilters(criteria *FilterCriteria) []models.Artist {
	filtered := []models.Artist{}

	for _, artist := range fe.artists {
		if fe.matchesCriteria(artist, criteria) {
			filtered = append(filtered, artist)
		}
	}

	return filtered
}

// matchesCriteria vérifie si un artiste correspond aux critères
func (fe *FilterEngine) matchesCriteria(artist models.Artist, criteria *FilterCriteria) bool {
	// Filtre par date de création
	if criteria.EnableCreationDateFilter {
		if artist.CreationDate < criteria.CreationDateMin || artist.CreationDate > criteria.CreationDateMax {
			return false
		}
	}

	// Filtre par année du premier album
	if criteria.EnableFirstAlbumFilter {
		albumYear := fe.extractYearFromFirstAlbum(artist.FirstAlbum)
		if albumYear < criteria.FirstAlbumYearMin || albumYear > criteria.FirstAlbumYearMax {
			return false
		}
	}

	// Filtre par nombre de membres
	if criteria.EnableMembersFilter {
		memberCount := len(artist.Members)
		if memberCount < criteria.MembersMin || memberCount > criteria.MembersMax {
			return false
		}
	}

	// Filtre par location (nécessite les données agrégées)
	if criteria.EnableLocationsFilter && len(criteria.Locations) > 0 {
		if !fe.matchesLocations(artist.ID, criteria.Locations) {
			return false
		}
	}

	return true
}

// extractYearFromFirstAlbum extrait l'année d'une date au format "dd-mm-yyyy"
func (fe *FilterEngine) extractYearFromFirstAlbum(dateStr string) int {
	date, err := ParseDate(dateStr)
	if err != nil {
		return 0
	}
	return date.Year()
}

// matchesLocations vérifie si un artiste a des concerts dans les locations spécifiées
func (fe *FilterEngine) matchesLocations(artistID int, wantedLocations []string) bool {
	aggregate, exists := fe.aggregates[artistID]
	if !exists {
		return false // Si pas de données, on considère que ça ne match pas
	}

	// Normaliser les locations recherchées
	normalizedWanted := make(map[string]bool)
	for _, loc := range wantedLocations {
		normalizedWanted[strings.ToLower(strings.TrimSpace(loc))] = true
	}

	// Vérifier si au moins une location de l'artiste match
	for _, location := range aggregate.Locations.Locations {
		city, country := ParseLocation(location)
		
		// Chercher match par ville
		if normalizedWanted[strings.ToLower(city)] {
			return true
		}
		
		// Chercher match par pays
		if normalizedWanted[strings.ToLower(country)] {
			return true
		}
		
		// Chercher match par location complète
		fullLoc := city + ", " + country
		if normalizedWanted[strings.ToLower(fullLoc)] {
			return true
		}
	}

	return false
}

// GetAvailableLocations retourne toutes les locations uniques disponibles
func (fe *FilterEngine) GetAvailableLocations() []string {
	locationSet := make(map[string]bool)
	locations := []string{}

	for artistID := range fe.aggregates {
		aggregate := fe.aggregates[artistID]
		for _, location := range aggregate.Locations.Locations {
			_, country := ParseLocation(location)
			
			// Ajouter le pays
			if country != "" && !locationSet[country] {
				locationSet[country] = true
				locations = append(locations, strings.ToUpper(country))
			}
		}
	}

	return locations
}

// GetDateRange retourne le range de dates de création disponible
func (fe *FilterEngine) GetDateRange() (min, max int) {
	if len(fe.artists) == 0 {
		return 1900, 2025
	}

	min = fe.artists[0].CreationDate
	max = fe.artists[0].CreationDate

	for _, artist := range fe.artists {
		if artist.CreationDate < min {
			min = artist.CreationDate
		}
		if artist.CreationDate > max {
			max = artist.CreationDate
		}
	}

	return min, max
}

// GetMembersRange retourne le range de nombre de membres disponible
func (fe *FilterEngine) GetMembersRange() (min, max int) {
	if len(fe.artists) == 0 {
		return 1, 10
	}

	min = len(fe.artists[0].Members)
	max = len(fe.artists[0].Members)

	for _, artist := range fe.artists {
		memberCount := len(artist.Members)
		if memberCount < min {
			min = memberCount
		}
		if memberCount > max {
			max = memberCount
		}
	}

	return min, max
}