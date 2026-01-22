package services

import (
	"groupie-tracker/models"
	"testing"
)

// createTestArtists crée des artistes de test
func createTestArtists() []models.Artist {
	return []models.Artist{
		{
			ID:           1,
			Name:         "Queen",
			Members:      []string{"Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"},
			CreationDate: 1970,
			FirstAlbum:   "14-12-1973",
		},
		{
			ID:           2,
			Name:         "The Beatles",
			Members:      []string{"John Lennon", "Paul McCartney", "George Harrison", "Ringo Starr"},
			CreationDate: 1960,
			FirstAlbum:   "22-03-1963",
		},
		{
			ID:           3,
			Name:         "Pink Floyd",
			Members:      []string{"Roger Waters", "David Gilmour", "Nick Mason", "Richard Wright"},
			CreationDate: 1965,
			FirstAlbum:   "05-08-1967",
		},
	}
}

func TestSearchEngine_SearchByArtistName(t *testing.T) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	results := engine.Search("queen")

	if len(results) == 0 {
		t.Error("Devrait trouver Queen")
	}

	if results[0].Type != SearchTypeArtist {
		t.Errorf("Type devrait être 'artist', got '%s'", results[0].Type)
	}

	if results[0].ArtistID != 1 {
		t.Errorf("ID devrait être 1, got %d", results[0].ArtistID)
	}
}

func TestSearchEngine_SearchByMember(t *testing.T) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	results := engine.Search("freddie")

	if len(results) == 0 {
		t.Error("Devrait trouver Freddie Mercury")
	}

	foundMember := false
	for _, result := range results {
		if result.Type == SearchTypeMember {
			foundMember = true
			if result.ArtistID != 1 {
				t.Errorf("Freddie devrait être dans Queen (ID=1), got %d", result.ArtistID)
			}
		}
	}

	if !foundMember {
		t.Error("Devrait trouver un résultat de type 'member'")
	}
}

func TestSearchEngine_CaseInsensitive(t *testing.T) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	testCases := []string{"QUEEN", "queen", "QuEeN", "qUeEn"}

	for _, query := range testCases {
		results := engine.Search(query)
		if len(results) == 0 {
			t.Errorf("Devrait trouver Queen pour la query '%s'", query)
		}
	}
}

func TestSearchEngine_PartialMatch(t *testing.T) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	results := engine.Search("beat")

	if len(results) == 0 {
		t.Error("Devrait trouver 'The Beatles' avec 'beat'")
	}
}

func TestSearchEngine_EmptyQuery(t *testing.T) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	results := engine.Search("")

	if len(results) != 0 {
		t.Error("Une query vide devrait retourner 0 résultats")
	}
}

func TestSearchEngine_ScoreCalculation(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Queen"},
		{ID: 2, Name: "Queen of the Stone Age"},
		{ID: 3, Name: "The Queen"},
	}
	engine := NewSearchEngine(artists)

	results := engine.Search("queen")

	// "Queen" devrait avoir le meilleur score (match exact)
	if len(results) > 0 {
		topResult := results[0]
		if topResult.ArtistID != 1 {
			t.Errorf("Match exact 'Queen' devrait être premier, got ID %d", topResult.ArtistID)
		}
		
		// Vérifier que le score est élevé
		if topResult.Score < 1000 {
			t.Errorf("Score pour match exact devrait être > 1000, got %d", topResult.Score)
		}
	}
}

func TestSearchEngine_SortByScore(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "The Queen"},          // Match au milieu
		{ID: 2, Name: "Queen"},              // Match exact
		{ID: 3, Name: "Queen of Hearts"},    // Match au début
	}
	engine := NewSearchEngine(artists)

	results := engine.Search("queen")

	// Vérifier que les résultats sont triés par score
	for i := 0; i < len(results)-1; i++ {
		if results[i].Score < results[i+1].Score {
			t.Errorf("Résultats mal triés: result[%d].Score=%d < result[%d].Score=%d",
				i, results[i].Score, i+1, results[i+1].Score)
		}
	}
}

func TestSearchEngine_GetSuggestions(t *testing.T) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	// Limiter à 2 suggestions
	suggestions := engine.GetSuggestions("roger", 2)

	if len(suggestions) > 2 {
		t.Errorf("Devrait retourner max 2 suggestions, got %d", len(suggestions))
	}
}

func TestSearchEngine_SearchByType(t *testing.T) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	// Chercher seulement les membres
	results := engine.SearchByType("roger", SearchTypeMember)

	// On devrait trouver "Roger Waters" et "Roger Taylor"
	if len(results) < 1 {
		t.Error("Devrait trouver au moins un membre nommé Roger")
	}

	for _, result := range results {
		if result.Type != SearchTypeMember {
			t.Errorf("Tous les résultats devraient être de type 'member', got '%s'", result.Type)
		}
	}
}

func TestSearchEngine_GetTopSuggestion(t *testing.T) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	top := engine.GetTopSuggestion("queen")

	if top == nil {
		t.Error("Devrait retourner une suggestion")
	}

	if top != nil && top.ArtistID != 1 {
		t.Errorf("Top suggestion pour 'queen' devrait être Queen (ID=1), got %d", top.ArtistID)
	}
}

func TestSearchEngine_SearchWithMinScore(t *testing.T) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	// Chercher avec un score minimum élevé
	results := engine.SearchWithMinScore("que", 500)

	// Avec un score min de 500, on devrait avoir moins de résultats
	allResults := engine.Search("que")

	if len(results) > len(allResults) {
		t.Error("SearchWithMinScore devrait retourner moins ou autant de résultats")
	}

	// Vérifier que tous les résultats ont un score >= 500
	for _, result := range results {
		if result.Score < 500 {
			t.Errorf("Résultat avec score %d < 500 ne devrait pas être retourné", result.Score)
		}
	}
}

func TestSearchEngine_HighlightMatch(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Queen"},
	}
	engine := NewSearchEngine(artists)

	results := engine.Search("uee")

	if len(results) == 0 {
		t.Skip("Pas de résultats pour tester highlight")
	}

	before, match, after := engine.HighlightMatch(results[0])

	// "Queen" -> "Q" + "uee" + "n"
	if before != "Q" {
		t.Errorf("Before devrait être 'Q', got '%s'", before)
	}

	if match != "uee" {
		t.Errorf("Match devrait être 'uee', got '%s'", match)
	}

	if after != "n" {
		t.Errorf("After devrait être 'n', got '%s'", after)
	}
}

func TestSearchEngine_NoDuplicates(t *testing.T) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	results := engine.Search("roger")

	// Vérifier qu'il n'y a pas de doublons
	seen := make(map[string]bool)
	for _, result := range results {
		key := result.ArtistName + "-" + result.MatchedText + "-" + string(result.Type)
		if seen[key] {
			t.Errorf("Doublon détecté: %s", key)
		}
		seen[key] = true
	}
}

func TestSearchEngine_MatchPositions(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Queen"},
	}
	engine := NewSearchEngine(artists)

	results := engine.Search("queen")

	if len(results) == 0 {
		t.Skip("Pas de résultats")
	}

	result := results[0]

	// Pour "queen" dans "Queen", match devrait commencer à 0
	if result.MatchStart != 0 {
		t.Errorf("MatchStart devrait être 0, got %d", result.MatchStart)
	}

	// Et finir à 5 (longueur de "queen")
	if result.MatchEnd != 5 {
		t.Errorf("MatchEnd devrait être 5, got %d", result.MatchEnd)
	}
}

// Benchmarks
func BenchmarkSearchEngine_Search(b *testing.B) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Search("queen")
	}
}

func BenchmarkSearchEngine_GetSuggestions(b *testing.B) {
	artists := createTestArtists()
	engine := NewSearchEngine(artists)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.GetSuggestions("roger", 10)
	}
}

func BenchmarkSearchEngine_CalculateScore(b *testing.B) {
	engine := NewSearchEngine([]models.Artist{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.calculateScore("Queen", "queen", 0, SearchTypeArtist)
	}
}