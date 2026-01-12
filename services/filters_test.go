package services

import (
	"groupie-tracker/models"
	"testing"
)

func TestFilterEngine_FilterByCreationDate(t *testing.T) {
	artists := createTestArtists()
	engine := NewFilterEngine(artists)

	criteria := NewFilterCriteria()
	criteria.EnableCreationDateFilter = true
	criteria.CreationDateMin = 1965
	criteria.CreationDateMax = 1970

	filtered := engine.ApplyFilters(criteria)

	// Devrait trouver Queen (1970) et Pink Floyd (1965)
	if len(filtered) != 2 {
		t.Errorf("Devrait trouver 2 artistes, got %d", len(filtered))
	}

	// Vérifier que Beatles (1960) n'est pas inclus
	for _, artist := range filtered {
		if artist.Name == "The Beatles" {
			t.Error("The Beatles ne devrait pas être dans les résultats (1960 < 1965)")
		}
	}
}

func TestFilterEngine_FilterByMembers(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Duo", Members: []string{"A", "B"}, CreationDate: 2000},
		{ID: 2, Name: "Trio", Members: []string{"A", "B", "C"}, CreationDate: 2000},
		{ID: 3, Name: "Quartet", Members: []string{"A", "B", "C", "D"}, CreationDate: 2000},
	}
	engine := NewFilterEngine(artists)

	criteria := NewFilterCriteria()
	criteria.EnableMembersFilter = true
	criteria.MembersMin = 2
	criteria.MembersMax = 3

	filtered := engine.ApplyFilters(criteria)

	// Devrait trouver Duo (2) et Trio (3), pas Quartet (4)
	if len(filtered) != 2 {
		t.Errorf("Devrait trouver 2 artistes, got %d", len(filtered))
	}

	for _, artist := range filtered {
		if artist.Name == "Quartet" {
			t.Error("Quartet ne devrait pas être dans les résultats (4 membres)")
		}
	}
}

func TestFilterEngine_NoFilterEnabled(t *testing.T) {
	artists := createTestArtists()
	engine := NewFilterEngine(artists)

	// Critères par défaut (tous désactivés)
	criteria := NewFilterCriteria()

	filtered := engine.ApplyFilters(criteria)

	// Sans filtre, on devrait avoir tous les artistes
	if len(filtered) != len(artists) {
		t.Errorf("Devrait retourner tous les artistes (%d), got %d", len(artists), len(filtered))
	}
}

func TestFilterEngine_MultipleFilters(t *testing.T) {
	artists := createTestArtists()
	engine := NewFilterEngine(artists)

	criteria := NewFilterCriteria()
	criteria.EnableCreationDateFilter = true
	criteria.CreationDateMin = 1960
	criteria.CreationDateMax = 1970
	criteria.EnableMembersFilter = true
	criteria.MembersMin = 4
	criteria.MembersMax = 4

	filtered := engine.ApplyFilters(criteria)

	// Tous les artistes de test ont 4 membres et sont entre 1960-1970
	if len(filtered) != 3 {
		t.Errorf("Devrait trouver 3 artistes, got %d", len(filtered))
	}
}

func TestFilterEngine_FilterByFirstAlbum(t *testing.T) {
	artists := createTestArtists()
	engine := NewFilterEngine(artists)

	criteria := NewFilterCriteria()
	criteria.EnableFirstAlbumFilter = true
	criteria.FirstAlbumYearMin = 1970
	criteria.FirstAlbumYearMax = 1975

	filtered := engine.ApplyFilters(criteria)

	// Queen (1973) devrait être trouvé
	// Beatles (1963) et Pink Floyd (1967) non
	if len(filtered) != 1 {
		t.Errorf("Devrait trouver 1 artiste (Queen), got %d", len(filtered))
	}

	if len(filtered) > 0 && filtered[0].Name != "Queen" {
		t.Errorf("Devrait trouver Queen, got %s", filtered[0].Name)
	}
}

func TestFilterEngine_GetDateRange(t *testing.T) {
	artists := createTestArtists()
	engine := NewFilterEngine(artists)

	min, max := engine.GetDateRange()

	if min != 1960 {
		t.Errorf("Min devrait être 1960 (Beatles), got %d", min)
	}

	if max != 1970 {
		t.Errorf("Max devrait être 1970 (Queen), got %d", max)
	}
}

func TestFilterEngine_GetMembersRange(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Solo", Members: []string{"A"}},
		{ID: 2, Name: "Large Band", Members: []string{"A", "B", "C", "D", "E", "F"}},
	}
	engine := NewFilterEngine(artists)

	min, max := engine.GetMembersRange()

	if min != 1 {
		t.Errorf("Min devrait être 1, got %d", min)
	}

	if max != 6 {
		t.Errorf("Max devrait être 6, got %d", max)
	}
}

func TestFilterEngine_ExtractYearFromFirstAlbum(t *testing.T) {
	engine := NewFilterEngine([]models.Artist{})

	testCases := []struct {
		input    string
		expected int
	}{
		{"14-12-1973", 1973},
		{"*23-08-2019", 2019},
		{"01-01-2000", 2000},
	}

	for _, tc := range testCases {
		year := engine.extractYearFromFirstAlbum(tc.input)
		if year != tc.expected {
			t.Errorf("extractYearFromFirstAlbum(%s) = %d, want %d", tc.input, year, tc.expected)
		}
	}
}

func TestFilterEngine_FilterLocations(t *testing.T) {
	// Ce test nécessite des données agrégées
	artists := []models.Artist{
		{ID: 1, Name: "Artist1", CreationDate: 2000},
	}
	engine := NewFilterEngine(artists)

	// Simuler des données agrégées
	engine.aggregates[1] = models.ArtistAggregate{
		Artist: artists[0],
		Locations: models.Location{
			ID:        1,
			Locations: []string{"los_angeles-usa", "paris-france", "tokyo-japan"},
		},
	}

	criteria := NewFilterCriteria()
	criteria.EnableLocationsFilter = true
	criteria.Locations = []string{"usa", "france"}

	filtered := engine.ApplyFilters(criteria)

	// Devrait trouver l'artiste car il a des concerts aux USA et en France
	if len(filtered) != 1 {
		t.Errorf("Devrait trouver 1 artiste, got %d", len(filtered))
	}
}

func TestFilterEngine_FilterLocationsNoMatch(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Artist1", CreationDate: 2000},
	}
	engine := NewFilterEngine(artists)

	engine.aggregates[1] = models.ArtistAggregate{
		Artist: artists[0],
		Locations: models.Location{
			ID:        1,
			Locations: []string{"los_angeles-usa"},
		},
	}

	criteria := NewFilterCriteria()
	criteria.EnableLocationsFilter = true
	criteria.Locations = []string{"france"} // Pas de concerts en France

	filtered := engine.ApplyFilters(criteria)

	// Ne devrait pas trouver l'artiste
	if len(filtered) != 0 {
		t.Errorf("Devrait trouver 0 artiste, got %d", len(filtered))
	}
}

func TestFilterEngine_GetAvailableLocations(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Artist1"},
		{ID: 2, Name: "Artist2"},
	}
	engine := NewFilterEngine(artists)

	engine.aggregates[1] = models.ArtistAggregate{
		Locations: models.Location{
			Locations: []string{"los_angeles-usa", "paris-france"},
		},
	}
	engine.aggregates[2] = models.ArtistAggregate{
		Locations: models.Location{
			Locations: []string{"tokyo-japan", "paris-france"},
		},
	}

	locations := engine.GetAvailableLocations()

	// Devrait avoir USA, FRANCE, JAPAN (pas de doublons)
	if len(locations) != 3 {
		t.Errorf("Devrait avoir 3 pays uniques, got %d", len(locations))
	}
}

// Benchmark
func BenchmarkFilterEngine_ApplyFilters(b *testing.B) {
	artists := createTestArtists()
	engine := NewFilterEngine(artists)
	criteria := NewFilterCriteria()
	criteria.EnableCreationDateFilter = true
	criteria.CreationDateMin = 1960
	criteria.CreationDateMax = 1970

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.ApplyFilters(criteria)
	}
}
