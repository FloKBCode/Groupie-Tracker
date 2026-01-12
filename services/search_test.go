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

	// Test recherche case-insensitive
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

	// Vérifier qu'on a trouvé le bon type
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

	// Test différentes casses
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

	// Recherche partielle
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

// Benchmark pour tester les performances
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
