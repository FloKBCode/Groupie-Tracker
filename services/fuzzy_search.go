package services

import (
	"strings"
)

// FuzzySearchEngine gère la recherche floue (tolérante aux fautes)
type FuzzySearchEngine struct {
	baseEngine *SearchEngine
}

// NewFuzzySearchEngine crée un moteur de recherche floue
func NewFuzzySearchEngine(baseEngine *SearchEngine) *FuzzySearchEngine {
	return &FuzzySearchEngine{
		baseEngine: baseEngine,
	}
}

// FuzzySearch effectue une recherche tolérante aux fautes de frappe
func (fse *FuzzySearchEngine) FuzzySearch(query string, maxDistance int) []SearchResult {
	query = strings.ToLower(strings.TrimSpace(query))

	if query == "" {
		return []SearchResult{}
	}

	// D'abord, essayer la recherche exacte
	exactResults := fse.baseEngine.Search(query)

	// Si on a des résultats exacts, les retourner
	if len(exactResults) > 0 {
		return exactResults
	}

	// Sinon, chercher avec distance de Levenshtein
	allResults := []SearchResult{}
	seen := make(map[string]bool)

	for _, artist := range fse.baseEngine.artists {
		// Recherche floue sur le nom de l'artiste
		if distance := levenshteinDistance(query, strings.ToLower(artist.Name)); distance <= maxDistance {
			key := artist.Name + "-artist"
			if !seen[key] {
				score := fse.calculateFuzzyScore(artist.Name, query, distance, SearchTypeArtist)
				allResults = append(allResults, SearchResult{
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

		// Recherche floue sur les membres
		for _, member := range artist.Members {
			if distance := levenshteinDistance(query, strings.ToLower(member)); distance <= maxDistance {
				key := artist.Name + "-" + member
				if !seen[key] {
					score := fse.calculateFuzzyScore(member, query, distance, SearchTypeMember)
					allResults = append(allResults, SearchResult{
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
	return fse.baseEngine.sortByScore(allResults)
}

// calculateFuzzyScore calcule le score pour une recherche floue
func (fse *FuzzySearchEngine) calculateFuzzyScore(text, query string, distance int, searchType SearchType) int {
	// Score de base
	score := 500

	// Pénalité pour la distance (plus la distance est grande, moins le score est élevé)
	score -= distance * 100

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

	// Bonus si les mots ont la même longueur
	if len(text) == len(query) {
		score += 100
	}

	return score
}

// levenshteinDistance calcule la distance de Levenshtein entre deux chaînes
// (nombre minimum d'opérations pour transformer s1 en s2)
func levenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	len1 := len(s1)
	len2 := len(s2)

	// Matrice de distances
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	// Initialisation
	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// Calcul des distances
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min3(
				matrix[i-1][j]+1,      // Suppression
				matrix[i][j-1]+1,      // Insertion
				matrix[i-1][j-1]+cost, // Substitution
			)
		}
	}

	return matrix[len1][len2]
}

// min3 retourne le minimum de trois entiers
func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// SmartSearch combine recherche exacte et floue
func (fse *FuzzySearchEngine) SmartSearch(query string) []SearchResult {
	// D'abord recherche exacte
	exactResults := fse.baseEngine.Search(query)

	// Si beaucoup de résultats exacts, les retourner
	if len(exactResults) >= 3 {
		return exactResults
	}

	// Sinon, ajouter des résultats flous avec distance max 2
	fuzzyResults := fse.FuzzySearch(query, 2)

	// Combiner les résultats (éviter doublons)
	seen := make(map[string]bool)
	combined := []SearchResult{}

	// Ajouter d'abord les résultats exacts
	for _, r := range exactResults {
		key := r.ArtistName + "-" + r.MatchedText
		if !seen[key] {
			combined = append(combined, r)
			seen[key] = true
		}
	}

	// Puis ajouter les résultats flous
	for _, r := range fuzzyResults {
		key := r.ArtistName + "-" + r.MatchedText
		if !seen[key] {
			combined = append(combined, r)
			seen[key] = true
		}
	}

	return combined
}
