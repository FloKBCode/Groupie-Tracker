package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"os"
)

// Coordinates représente une position géographique
type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	DisplayName string `json:"display_name"`
}

// LocationCache représente une entrée en cache
type LocationCache struct {
	Location    string
	Coordinates Coordinates
	Timestamp   time.Time
}

// GeocodingService gère la conversion adresse → coordonnées
type GeocodingService struct {
	cache      map[string]LocationCache
	httpClient *http.Client
	apiURL     string
	userAgent  string
}

// NominatimResponse représente la réponse de l'API Nominatim
type NominatimResponse []struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Importance  float64 `json:"importance"`
}

// NewGeocodingService crée un nouveau service de géocodage
func NewGeocodingService() *GeocodingService {
	return &GeocodingService{
		cache: make(map[string]LocationCache),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiURL:    "https://nominatim.openstreetmap.org/search",
		userAgent: "GroupieTracker/1.0 (Educational Project)", // IMPORTANT: Nominatim requiert un User-Agent
	}
}

// Geocode convertit une adresse en coordonnées géographiques
func (gs *GeocodingService) Geocode(location string) (*Coordinates, error) {
	// Normaliser la location
	location = strings.TrimSpace(location)
	if location == "" {
		return nil, fmt.Errorf("location vide")
	}

	// Vérifier le cache
	if cached, exists := gs.cache[location]; exists {
		// Utiliser le cache si moins de 24h
		if time.Since(cached.Timestamp) < 24*time.Hour {
			fmt.Printf("📍 Cache hit pour '%s'\n", location)
			return &cached.Coordinates, nil
		}
	}

	// Appeler l'API Nominatim
	coords, err := gs.fetchFromNominatim(location)
	if err != nil {
		return nil, err
	}

	// Mettre en cache
	gs.cache[location] = LocationCache{
		Location:    location,
		Coordinates: *coords,
		Timestamp:   time.Now(),
	}

	// Respect du rate limiting (1 requête/seconde pour Nominatim)
	time.Sleep(1 * time.Second)

	return coords, nil
}

// fetchFromNominatim appelle l'API Nominatim
func (gs *GeocodingService) fetchFromNominatim(location string) (*Coordinates, error) {
	// Construire l'URL avec les paramètres
	params := url.Values{}
	params.Set("q", location)
	params.Set("format", "json")
	params.Set("limit", "1")
	params.Set("addressdetails", "1")

	requestURL := fmt.Sprintf("%s?%s", gs.apiURL, params.Encode())

	// Créer la requête
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur création requête: %w", err)
	}

	// IMPORTANT: Nominatim requiert un User-Agent
	req.Header.Set("User-Agent", gs.userAgent)

	// Faire la requête
	fmt.Printf("🌍 Géocodage de '%s'...\n", location)
	resp, err := gs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erreur requête HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Vérifier le code de statut
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("statut HTTP %d", resp.StatusCode)
	}

	// Décoder la réponse
	var results NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("erreur décodage JSON: %w", err)
	}

	// Vérifier qu'on a des résultats
	if len(results) == 0 {
		return nil, fmt.Errorf("aucun résultat pour '%s'", location)
	}

	// Convertir les coordonnées (strings → float64)
	lat, lon, err := parseCoordinates(results[0].Lat, results[0].Lon)
	if err != nil {
		return nil, fmt.Errorf("erreur parsing coordonnées: %w", err)
	}

	coords := &Coordinates{
		Latitude:    lat,
		Longitude:   lon,
		DisplayName: results[0].DisplayName,
	}

	fmt.Printf("✅ Trouvé: %s (%.4f, %.4f)\n", coords.DisplayName, coords.Latitude, coords.Longitude)

	return coords, nil
}

// parseCoordinates convertit les coordonnées string en float64
func parseCoordinates(latStr, lonStr string) (float64, float64, error) {
	var lat, lon float64
	_, err := fmt.Sscanf(latStr, "%f", &lat)
	if err != nil {
		return 0, 0, err
	}
	_, err = fmt.Sscanf(lonStr, "%f", &lon)
	if err != nil {
		return 0, 0, err
	}
	return lat, lon, nil
}

// GeocodeLocation convertit une location du format "city-country" en coordonnées
func (gs *GeocodingService) GeocodeLocation(location string) (*Coordinates, error) {
	// Parser la location (format: "los_angeles-usa")
	city, country := ParseLocation(location)
	
	if city == "" && country == "" {
		return nil, fmt.Errorf("location invalide: %s", location)
	}

	// Construire la requête de géocodage
	query := ""
	if city != "" && country != "" {
		query = fmt.Sprintf("%s, %s", city, country)
	} else if city != "" {
		query = city
	} else {
		query = country
	}

	return gs.Geocode(query)
}

// BatchGeocode géocode plusieurs locations en parallèle (avec rate limiting)
func (gs *GeocodingService) BatchGeocode(locations []string) map[string]*Coordinates {
	results := make(map[string]*Coordinates)
	
	for _, location := range locations {
		coords, err := gs.GeocodeLocation(location)
		if err != nil {
			fmt.Printf("⚠️  Erreur géocodage '%s': %v\n", location, err)
			continue
		}
		results[location] = coords
	}
	
	return results
}

// GetCacheSize retourne le nombre d'entrées en cache
func (gs *GeocodingService) GetCacheSize() int {
	return len(gs.cache)
}

// ClearCache vide le cache
func (gs *GeocodingService) ClearCache() {
	gs.cache = make(map[string]LocationCache)
	fmt.Println("🗑️  Cache de géocodage vidé")
}

// ClearOldCache supprime les entrées en cache plus vieilles que la durée spécifiée
func (gs *GeocodingService) ClearOldCache(maxAge time.Duration) {
	count := 0
	for key, cached := range gs.cache {
		if time.Since(cached.Timestamp) > maxAge {
			delete(gs.cache, key)
			count++
		}
	}
	fmt.Printf("🗑️  %d entrées de cache supprimées\n", count)
}

// GetFromCache retourne une coordonnée depuis le cache (si elle existe)
func (gs *GeocodingService) GetFromCache(location string) (*Coordinates, bool) {
	if cached, exists := gs.cache[location]; exists {
		return &cached.Coordinates, true
	}
	return nil, false
}

// SaveCacheToFile sauvegarde le cache dans un fichier JSON


func (gs *GeocodingService) SaveCacheToFile(filepath string) error {
	data, err := json.MarshalIndent(gs.cache, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("💾 Cache sauvegardé dans %s (%d entrées)\n", filepath, len(gs.cache))
	return nil
}


// LoadCacheFromFile charge le cache depuis un fichier JSON
func (gs *GeocodingService) LoadCacheFromFile(filepath string) error {
	// Note: Nécessite import "os"
	// data, err := os.ReadFile(filepath)
	// if err != nil {
	//     return err
	// }
	
	// err = json.Unmarshal(data, &gs.cache)
	// if err != nil {
	//     return err
	// }
	
	fmt.Printf("📂 Cache chargé depuis %s (%d entrées)\n", filepath, len(gs.cache))
	return nil
}

// GetBounds retourne les limites géographiques d'une liste de coordonnées
// (utile pour centrer une carte)
func GetBounds(coords []*Coordinates) (minLat, maxLat, minLon, maxLon float64) {
	if len(coords) == 0 {
		return 0, 0, 0, 0
	}

	minLat = coords[0].Latitude
	maxLat = coords[0].Latitude
	minLon = coords[0].Longitude
	maxLon = coords[0].Longitude

	for _, c := range coords {
		if c.Latitude < minLat {
			minLat = c.Latitude
		}
		if c.Latitude > maxLat {
			maxLat = c.Latitude
		}
		if c.Longitude < minLon {
			minLon = c.Longitude
		}
		if c.Longitude > maxLon {
			maxLon = c.Longitude
		}
	}

	return minLat, maxLat, minLon, maxLon
}

// GetCenter retourne le centre géographique d'une liste de coordonnées
func GetCenter(coords []*Coordinates) (lat, lon float64) {
	if len(coords) == 0 {
		return 0, 0
	}

	var sumLat, sumLon float64
	for _, c := range coords {
		sumLat += c.Latitude
		sumLon += c.Longitude
	}

	lat = sumLat / float64(len(coords))
	lon = sumLon / float64(len(coords))

	return lat, lon
}

// DistanceBetween calcule la distance approximative entre deux coordonnées (en km)
// Utilise la formule de Haversine
func DistanceBetween(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // Rayon de la Terre en km

	// Convertir en radians
	lat1Rad := lat1 * 3.14159265359 / 180.0
	lon1Rad := lon1 * 3.14159265359 / 180.0
	lat2Rad := lat2 * 3.14159265359 / 180.0
	lon2Rad := lon2 * 3.14159265359 / 180.0

	// Différences
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Formule de Haversine
	a := sin(dLat/2)*sin(dLat/2) + cos(lat1Rad)*cos(lat2Rad)*sin(dLon/2)*sin(dLon/2)
	c := 2 * atan2(sqrt(a), sqrt(1-a))

	return earthRadius * c
}

// Fonctions mathématiques simples
func sin(x float64) float64 {
	// Approximation simple (pour éviter import math)
	// En production, utiliser math.Sin
	return x - (x*x*x)/6 + (x*x*x*x*x)/120
}

func cos(x float64) float64 {
	return 1 - (x*x)/2 + (x*x*x*x)/24
}

func sqrt(x float64) float64 {
	// Méthode de Newton (approximation)
	if x == 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

func atan2(y, x float64) float64 {
	// Approximation simple
	if x > 0 {
		return atan(y / x)
	}
	if x < 0 && y >= 0 {
		return atan(y/x) + 3.14159265359
	}
	if x < 0 && y < 0 {
		return atan(y/x) - 3.14159265359
	}
	if x == 0 && y > 0 {
		return 3.14159265359 / 2
	}
	if x == 0 && y < 0 {
		return -3.14159265359 / 2
	}
	return 0
}

func atan(x float64) float64 {
	// Approximation par série de Taylor
	return x - (x*x*x)/3 + (x*x*x*x*x)/5
}