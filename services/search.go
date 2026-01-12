package services

import (
	"groupie-tracker/models"
	"strings"
)

// SearchType représente le type de résultat de recherche
type SearchType string

const (
	SearchTypeArtist   SearchType = "artist"
	SearchTypeMember   SearchType = "member"
	SearchTypeLocation SearchType = "location"
	SearchTypeDate     SearchType = "date"
)

// SearchResult représente un résultat de recherche avec son type
type SearchResult struct {
	ArtistID    int        // ID de l'artiste correspondant
	ArtistName  string     // Nom de l'artiste
	MatchedText string     // Texte qui a matché
	Type        SearchType // Type de match
}

// SearchEngine gère la recherche dans les artistes
type SearchEngine struct {
	artists    []models.Artist
	aggregates map[int]models.ArtistAggregate // Cache des données agrégées
}

// NewSearchEngine crée une nouvelle instance du moteur de recherche
func NewSearchEngine(artists []models.Artist) *SearchEngine {
	return &SearchEngine{
		artists:    artists,
		aggregates: make(map[int]models.ArtistAggregate),
	}
}

// LoadAggregateData charge les données agrégées pour un artiste (lazy loading)
func (se *SearchEngine) LoadAggregateData(artistID int) error {
	// Si déjà en cache, ne rien faire
	if _, exists := se.aggregates[artistID]; exists {
		return nil
	}

	// Trouver l'artiste
	var artist models.Artist
	for _, a := range se.artists {
		if a.ID == artistID {
			artist = a
			break
		}
	}

	// Charger les données agrégées
	aggregate, err := AggregateArtist(artist)
	if err != nil {
		return err
	}

	se.aggregates[artistID] = aggregate
	return nil
}

// Search effectue une recherche case-insensitive sur tous les champs
func (se *SearchEngine) Search(query string) []SearchResult {
	if query == "" {
		return []SearchResult{}
	}

	query = strings.ToLower(strings.TrimSpace(query))
	results := []SearchResult{}
	seen := make(map[string]bool) // Pour éviter les doublons

	for _, artist := range se.artists {
		// Recherche dans le nom de l'artiste
		if strings.Contains(strings.ToLower(artist.Name), query) {
			key := artist.Name + "-artist"
			if !seen[key] {
				results = append(results, SearchResult{
					ArtistID:    artist.ID,
					ArtistName:  artist.Name,
					MatchedText: artist.Name,
					Type:        SearchTypeArtist,
				})
				seen[key] = true
			}
		}

		// Recherche dans les membres
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				key := artist.Name + "-" + member
				if !seen[key] {
					results = append(results, SearchResult{
						ArtistID:    artist.ID,
						ArtistName:  artist.Name,
						MatchedText: member,
						Type:        SearchTypeMember,
					})
					seen[key] = true
				}
			}
		}

		// Recherche dans la date du premier album
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			key := artist.Name + "-" + artist.FirstAlbum
			if !seen[key] {
				results = append(results, SearchResult{
					ArtistID:    artist.ID,
					ArtistName:  artist.Name,
					MatchedText: artist.FirstAlbum,
					Type:        SearchTypeDate,
				})
				seen[key] = true
			}
		}

		// Recherche dans les locations (nécessite le chargement des données)
		if aggregate, exists := se.aggregates[artist.ID]; exists {
			for _, location := range aggregate.Locations.Locations {
				city, country := ParseLocation(location)
				fullLocation := city + " " + country
				
				if strings.Contains(strings.ToLower(fullLocation), query) {
					key := artist.Name + "-" + location
					if !seen[key] {
						results = append(results, SearchResult{
							ArtistID:    artist.ID,
							ArtistName:  artist.Name,
							MatchedText: city + ", " + strings.ToUpper(country),
							Type:        SearchTypeLocation,
						})
						seen[key] = true
					}
				}
			}
		}
	}

	return results
}

// GetSuggestions retourne des suggestions de recherche limitées à maxResults
func (se *SearchEngine) GetSuggestions(query string, maxResults int) []SearchResult {
	results := se.Search(query)
	
	if len(results) > maxResults {
		return results[:maxResults]
	}
	
	return results
}

// SearchByType effectue une recherche filtrée par type
func (se *SearchEngine) SearchByType(query string, searchType SearchType) []SearchResult {
	allResults := se.Search(query)
	filtered := []SearchResult{}

	for _, result := range allResults {
		if result.Type == searchType {
			filtered = append(filtered, result)
		}
	}

	return filtered
}