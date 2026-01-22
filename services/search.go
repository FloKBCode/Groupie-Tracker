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
	Score       int        // Score de pertinence (plus élevé = plus pertinent)
	MatchStart  int        // Position de début du match (pour highlighting)
	MatchEnd    int        // Position de fin du match
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

// Search effectue une recherche case-insensitive sur tous les champs avec scoring
func (se *SearchEngine) Search(query string) []SearchResult {
	if query == "" {
		return []SearchResult{}
	}

	query = strings.ToLower(strings.TrimSpace(query))
	results := []SearchResult{}
	seen := make(map[string]bool) // Pour éviter les doublons

	for _, artist := range se.artists {
		// Recherche dans le nom de l'artiste
		if matchPos := strings.Index(strings.ToLower(artist.Name), query); matchPos != -1 {
			key := artist.Name + "-artist"
			if !seen[key] {
				score := se.calculateScore(artist.Name, query, matchPos, SearchTypeArtist)
				results = append(results, SearchResult{
					ArtistID:    artist.ID,
					ArtistName:  artist.Name,
					MatchedText: artist.Name,
					Type:        SearchTypeArtist,
					Score:       score,
					MatchStart:  matchPos,
					MatchEnd:    matchPos + len(query),
				})
				seen[key] = true
			}
		}

		// Recherche dans les membres
		for _, member := range artist.Members {
			if matchPos := strings.Index(strings.ToLower(member), query); matchPos != -1 {
				key := artist.Name + "-" + member
				if !seen[key] {
					score := se.calculateScore(member, query, matchPos, SearchTypeMember)
					results = append(results, SearchResult{
						ArtistID:    artist.ID,
						ArtistName:  artist.Name,
						MatchedText: member,
						Type:        SearchTypeMember,
						Score:       score,
						MatchStart:  matchPos,
						MatchEnd:    matchPos + len(query),
					})
					seen[key] = true
				}
			}
		}

		// Recherche dans la date du premier album
		if matchPos := strings.Index(strings.ToLower(artist.FirstAlbum), query); matchPos != -1 {
			key := artist.Name + "-" + artist.FirstAlbum
			if !seen[key] {
				score := se.calculateScore(artist.FirstAlbum, query, matchPos, SearchTypeDate)
				results = append(results, SearchResult{
					ArtistID:    artist.ID,
					ArtistName:  artist.Name,
					MatchedText: artist.FirstAlbum,
					Type:        SearchTypeDate,
					Score:       score,
					MatchStart:  matchPos,
					MatchEnd:    matchPos + len(query),
				})
				seen[key] = true
			}
		}

		// Recherche dans les locations (nécessite le chargement des données)
		if aggregate, exists := se.aggregates[artist.ID]; exists {
			for _, location := range aggregate.Locations.Locations {
				city, country := ParseLocation(location)
				fullLocation := city + " " + country

				if matchPos := strings.Index(strings.ToLower(fullLocation), query); matchPos != -1 {
					key := artist.Name + "-" + location
					if !seen[key] {
						displayLocation := city + ", " + strings.ToUpper(country)
						score := se.calculateScore(fullLocation, query, matchPos, SearchTypeLocation)
						results = append(results, SearchResult{
							ArtistID:    artist.ID,
							ArtistName:  artist.Name,
							MatchedText: displayLocation,
							Type:        SearchTypeLocation,
							Score:       score,
							MatchStart:  matchPos,
							MatchEnd:    matchPos + len(query),
						})
						seen[key] = true
					}
				}
			}
		}
	}

	// Trier par score (plus pertinent en premier)
	results = se.sortByScore(results)

	return results
}

// calculateScore calcule un score de pertinence pour un match
func (se *SearchEngine) calculateScore(text, query string, matchPos int, searchType SearchType) int {
	score := 100

	textLower := strings.ToLower(text)
	queryLower := strings.ToLower(query)

	// Match exact = score maximum
	if textLower == queryLower {
		score += 1000
	}

	// Match au début du texte = bonus
	if matchPos == 0 {
		score += 500
	}

	// Match après un espace = bonus (début de mot)
	if matchPos > 0 && text[matchPos-1] == ' ' {
		score += 300
	}

	// Plus le texte est court, plus il est pertinent
	score += (100 - len(text))

	// Bonus selon le type
	switch searchType {
	case SearchTypeArtist:
		score += 400 // Artistes en priorité
	case SearchTypeMember:
		score += 300
	case SearchTypeLocation:
		score += 200
	case SearchTypeDate:
		score += 100
	}

	// Bonus si la query couvre une grande partie du texte
	coveragePercent := (len(query) * 100) / len(text)
	score += coveragePercent * 2

	return score
}

// sortByScore trie les résultats par score décroissant
func (se *SearchEngine) sortByScore(results []SearchResult) []SearchResult {
	// Bubble sort simple (suffisant pour un petit nombre de résultats)
	n := len(results)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if results[j].Score < results[j+1].Score {
				results[j], results[j+1] = results[j+1], results[j]
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

// GetTopSuggestion retourne la meilleure suggestion (score le plus élevé)
func (se *SearchEngine) GetTopSuggestion(query string) *SearchResult {
	results := se.Search(query)
	if len(results) > 0 {
		return &results[0]
	}
	return nil
}

// SearchWithMinScore retourne uniquement les résultats au-dessus d'un score minimum
func (se *SearchEngine) SearchWithMinScore(query string, minScore int) []SearchResult {
	allResults := se.Search(query)
	filtered := []SearchResult{}

	for _, result := range allResults {
		if result.Score >= minScore {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

// HighlightMatch retourne le texte avec le match mis en évidence
func (se *SearchEngine) HighlightMatch(result SearchResult) (before, match, after string) {
	text := result.MatchedText
	
	if result.MatchStart < 0 || result.MatchEnd > len(text) {
		return text, "", ""
	}

	before = text[:result.MatchStart]
	match = text[result.MatchStart:result.MatchEnd]
	after = text[result.MatchEnd:]

	return before, match, after
}