package services

import (
	"testing"
	"time"
)

func TestNewGeocodingService(t *testing.T) {
	gs := NewGeocodingService()
	
	if gs == nil {
		t.Fatal("NewGeocodingService devrait retourner un service valide")
	}
	
	if gs.cache == nil {
		t.Error("Le cache devrait être initialisé")
	}
	
	if gs.httpClient == nil {
		t.Error("Le client HTTP devrait être initialisé")
	}
}

func TestParseCoordinates(t *testing.T) {
	testCases := []struct {
		latStr   string
		lonStr   string
		wantLat  float64
		wantLon  float64
		wantErr  bool
	}{
		{"48.8566", "2.3522", 48.8566, 2.3522, false},  // Paris
		{"51.5074", "-0.1278", 51.5074, -0.1278, false}, // Londres
		{"invalid", "2.0", 0, 0, true},                  // Erreur
		{"48.0", "invalid", 0, 0, true},                 // Erreur
	}
	
	for _, tc := range testCases {
		lat, lon, err := parseCoordinates(tc.latStr, tc.lonStr)
		
		if tc.wantErr {
			if err == nil {
				t.Errorf("parseCoordinates('%s', '%s') devrait retourner une erreur", tc.latStr, tc.lonStr)
			}
			continue
		}
		
		if err != nil {
			t.Errorf("parseCoordinates('%s', '%s') erreur inattendue: %v", tc.latStr, tc.lonStr, err)
			continue
		}
		
		if lat != tc.wantLat || lon != tc.wantLon {
			t.Errorf("parseCoordinates('%s', '%s') = (%.4f, %.4f), want (%.4f, %.4f)",
				tc.latStr, tc.lonStr, lat, lon, tc.wantLat, tc.wantLon)
		}
	}
}

func TestGeocodeCache(t *testing.T) {
	gs := NewGeocodingService()
	
	// Simuler une entrée en cache
	testLocation := "Paris, France"
	testCoords := Coordinates{
		Latitude:    48.8566,
		Longitude:   2.3522,
		DisplayName: "Paris, Île-de-France, France",
	}
	
	gs.cache[testLocation] = LocationCache{
		Location:    testLocation,
		Coordinates: testCoords,
		Timestamp:   time.Now(),
	}
	
	// Tester GetFromCache
	coords, found := gs.GetFromCache(testLocation)
	if !found {
		t.Error("Devrait trouver Paris dans le cache")
	}
	
	if coords.Latitude != testCoords.Latitude {
		t.Errorf("Latitude = %.4f, want %.4f", coords.Latitude, testCoords.Latitude)
	}
}

func TestGetCacheSize(t *testing.T) {
	gs := NewGeocodingService()
	
	if gs.GetCacheSize() != 0 {
		t.Error("Cache devrait être vide au départ")
	}
	
	// Ajouter des entrées
	gs.cache["Paris"] = LocationCache{
		Location:  "Paris",
		Timestamp: time.Now(),
	}
	gs.cache["London"] = LocationCache{
		Location:  "London",
		Timestamp: time.Now(),
	}
	
	if gs.GetCacheSize() != 2 {
		t.Errorf("Cache size = %d, want 2", gs.GetCacheSize())
	}
}

func TestClearCache(t *testing.T) {
	gs := NewGeocodingService()
	
	// Ajouter des entrées
	gs.cache["Paris"] = LocationCache{Location: "Paris", Timestamp: time.Now()}
	gs.cache["London"] = LocationCache{Location: "London", Timestamp: time.Now()}
	
	gs.ClearCache()
	
	if gs.GetCacheSize() != 0 {
		t.Errorf("Cache devrait être vide après Clear, got %d entries", gs.GetCacheSize())
	}
}

func TestClearOldCache(t *testing.T) {
	gs := NewGeocodingService()
	
	// Ajouter une entrée récente
	gs.cache["Paris"] = LocationCache{
		Location:  "Paris",
		Timestamp: time.Now(),
	}
	
	// Ajouter une entrée ancienne (2 jours)
	gs.cache["London"] = LocationCache{
		Location:  "London",
		Timestamp: time.Now().Add(-48 * time.Hour),
	}
	
	// Supprimer les entrées > 24h
	gs.ClearOldCache(24 * time.Hour)
	
	// Paris devrait rester, London devrait être supprimé
	if _, found := gs.GetFromCache("Paris"); !found {
		t.Error("Paris devrait rester dans le cache")
	}
	
	if _, found := gs.GetFromCache("London"); found {
		t.Error("London devrait être supprimé du cache")
	}
}

func TestGetBounds(t *testing.T) {
	coords := []*Coordinates{
		{Latitude: 48.8566, Longitude: 2.3522},   // Paris
		{Latitude: 51.5074, Longitude: -0.1278},  // Londres
		{Latitude: 40.7128, Longitude: -74.0060}, // New York
	}
	
	minLat, maxLat, minLon, maxLon := GetBounds(coords)
	
	// Vérifier que les limites sont correctes
	if minLat > 40.7128 || maxLat < 51.5074 {
		t.Error("Limites de latitude incorrectes")
	}
	
	if minLon > -74.0060 || maxLon < 2.3522 {
		t.Error("Limites de longitude incorrectes")
	}
}

func TestGetCenter(t *testing.T) {
	coords := []*Coordinates{
		{Latitude: 0, Longitude: 0},
		{Latitude: 10, Longitude: 10},
	}
	
	lat, lon := GetCenter(coords)
	
	// Centre devrait être (5, 5)
	if lat != 5 || lon != 5 {
		t.Errorf("Centre = (%.2f, %.2f), want (5, 5)", lat, lon)
	}
}

func TestGetCenterEmpty(t *testing.T) {
	coords := []*Coordinates{}
	
	lat, lon := GetCenter(coords)
	
	if lat != 0 || lon != 0 {
		t.Error("Centre d'une liste vide devrait être (0, 0)")
	}
}

func TestDistanceBetween(t *testing.T) {
	// Paris à Londres : environ 344 km
	parisLat, parisLon := 48.8566, 2.3522
	londonLat, londonLon := 51.5074, -0.1278
	
	distance := DistanceBetween(parisLat, parisLon, londonLat, londonLon)
	
	// La distance devrait être autour de 344 km (avec une marge d'erreur)
	if distance < 300 || distance > 400 {
		t.Errorf("Distance Paris-Londres = %.2f km, devrait être ~344 km", distance)
	}
}

func TestDistanceSamePoint(t *testing.T) {
	distance := DistanceBetween(48.8566, 2.3522, 48.8566, 2.3522)
	
	if distance > 1 {
		t.Errorf("Distance entre le même point devrait être ~0, got %.2f", distance)
	}
}

// Test d'intégration (nécessite une connexion internet)
// Décommenter pour tester avec l'API réelle
/*
func TestGeocodeIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	gs := NewGeocodingService()
	
	coords, err := gs.Geocode("Paris, France")
	if err != nil {
		t.Fatalf("Erreur géocodage: %v", err)
	}
	
	// Paris devrait être autour de 48.85, 2.35
	if coords.Latitude < 48.8 || coords.Latitude > 48.9 {
		t.Errorf("Latitude de Paris incorrecte: %.4f", coords.Latitude)
	}
	
	if coords.Longitude < 2.3 || coords.Longitude > 2.4 {
		t.Errorf("Longitude de Paris incorrecte: %.4f", coords.Longitude)
	}
}

func TestGeocodeLocationIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	gs := NewGeocodingService()
	
	// Format "city-country"
	coords, err := gs.GeocodeLocation("paris-france")
	if err != nil {
		t.Fatalf("Erreur géocodage location: %v", err)
	}
	
	if coords.Latitude < 48 || coords.Latitude > 49 {
		t.Errorf("Coordonnées de Paris incorrectes: %.4f", coords.Latitude)
	}
}

func TestBatchGeocodeIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	gs := NewGeocodingService()
	
	locations := []string{
		"paris-france",
		"london-uk",
		"berlin-germany",
	}
	
	results := gs.BatchGeocode(locations)
	
	if len(results) != 3 {
		t.Errorf("Devrait avoir 3 résultats, got %d", len(results))
	}
}
*/

// Benchmarks
func BenchmarkParseCoordinates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseCoordinates("48.8566", "2.3522")
	}
}

func BenchmarkGetCenter(b *testing.B) {
	coords := []*Coordinates{
		{Latitude: 48.8566, Longitude: 2.3522},
		{Latitude: 51.5074, Longitude: -0.1278},
		{Latitude: 40.7128, Longitude: -74.0060},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCenter(coords)
	}
}

func BenchmarkDistanceBetween(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DistanceBetween(48.8566, 2.3522, 51.5074, -0.1278)
	}
}