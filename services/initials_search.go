package services

import (
	"strings"
)

// InitialsSearchEngine gère la recherche par initiales
type InitialsSearchEngine struct {
	baseEngine *SearchEngine
}

// NewInitialsSearchEngine crée un moteur de recherche par initiales
func NewInitialsSearchEngine(baseEngine *SearchEngine) *InitialsSearchEngine {
	return &InitialsSearchEngine{
		baseEngine: baseEngine,
	}
}

// SearchByInitials recherche par initiales
// Exemples: "fm" → "Freddie Mercury", "qotsa" → "Queen of the Stone Age"
func (ise *InitialsSearchEngine) SearchByInitials(initials string) []SearchResult {
	initials = strings.ToLower(strings.TrimSpace(initials))

	if initials == "" {
		return []SearchResult{}
	}

	results := []SearchResult{}
	seen := make(map[string]bool)

	for _, artist := range ise.baseEngine.artists {
		// Recherche dans le nom de l'artiste
		if ise.matchesInitials(artist.Name, initials) {
			key := artist.Name + "-artist"
			if !seen[key] {
				score := ise.calculateInitialsScore(artist.Name, initials, SearchTypeArtist)
				results = append(results, SearchResult{
					ArtistID:    artist.ID,
					ArtistName:  artist.Name,
					MatchedText: artist.Name,
					Type:        SearchTypeArtist,
					Score:       score,
					MatchStart:  0,
					MatchEnd:    len(artist.Name),
				})
				seen[key] = true
			}
		}

		// Recherche dans les membres
		for _, member := range artist.Members {
			if ise.matchesInitials(member, initials) {
				key := artist.Name + "-" + member
				if !seen[key] {
					score := ise.calculateInitialsScore(member, initials, SearchTypeMember)
					results = append(results, SearchResult{
						ArtistID:    artist.ID,
						ArtistName:  artist.Name,
						MatchedText: member,
						Type:        SearchTypeMember,
						Score:       score,
						MatchStart:  0,
						MatchEnd:    len(member),
					})
					seen[key] = true
				}
			}
		}
	}

	// Trier par score
	return ise.baseEngine.sortByScore(results)
}

// matchesInitials vérifie si un texte correspond aux initiales
func (ise *InitialsSearchEngine) matchesInitials(text, initials string) bool {
	text = strings.ToLower(text)
	initials = strings.ToLower(initials)

	// Extraire les initiales du texte
	words := strings.Fields(text)
	textInitials := ""

	for _, word := range words {
		if len(word) > 0 {
			textInitials += string(word[0])
		}
	}

	// Vérifier si les initiales correspondent
	return textInitials == initials || strings.HasPrefix(textInitials, initials)
}

// extractInitials extrait les initiales d'un texte
func (ise *InitialsSearchEngine) extractInitials(text string) string {
	words := strings.Fields(strings.ToLower(text))
	initials := ""

	for _, word := range words {
		if len(word) > 0 {
			initials += string(word[0])
		}
	}

	return initials
}

// calculateInitialsScore calcule le score pour une recherche par initiales
func (ise *InitialsSearchEngine) calculateInitialsScore(text, initials string, searchType SearchType) int {
	score := 400

	textInitials := ise.extractInitials(text)

	// Match exact des initiales = bonus
	if textInitials == initials {
		score += 500
	}

	// Match partiel (début) = bonus moindre
	if strings.HasPrefix(textInitials, initials) {
		score += 300
	}

	// Bonus selon le type
	switch searchType {
	case SearchTypeArtist:
		score += 300
	case SearchTypeMember:
		score += 200
	case SearchTypeLocation:
		score += 100
	case SearchTypeDate:
		score += 50
	}

	// Bonus si peu de mots (initiales plus spécifiques)
	wordCount := len(strings.Fields(text))
	if wordCount == len(initials) {
		score += 200
	}

	return score
}

// SmartInitialsSearch combine recherche normale et initiales
func (ise *InitialsSearchEngine) SmartInitialsSearch(query string) []SearchResult {
	query = strings.ToLower(strings.TrimSpace(query))

	// Si la query a des espaces, utiliser recherche normale
	if strings.Contains(query, " ") {
		return ise.baseEngine.Search(query)
	}

	// Si la query est courte (2-5 caractères), essayer initiales
	if len(query) >= 2 && len(query) <= 5 {
		initialsResults := ise.SearchByInitials(query)

		// Si on trouve des résultats par initiales, les combiner avec recherche normale
		normalResults := ise.baseEngine.Search(query)

		// Combiner et éliminer doublons
		seen := make(map[string]bool)
		combined := []SearchResult{}

		// Priorité aux résultats normaux
		for _, r := range normalResults {
			key := r.ArtistName + "-" + r.MatchedText
			if !seen[key] {
				combined = append(combined, r)
				seen[key] = true
			}
		}

		// Ajouter résultats par initiales
		for _, r := range initialsResults {
			key := r.ArtistName + "-" + r.MatchedText
			if !seen[key] {
				combined = append(combined, r)
				seen[key] = true
			}
		}

		return combined
	}

	// Sinon, recherche normale
	return ise.baseEngine.Search(query)
}
