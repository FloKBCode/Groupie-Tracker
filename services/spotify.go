package services

import (
	"fmt"
	"net/url"
	"strings"
)

// SpotifyService gère l'intégration Spotify
type SpotifyService struct {
	// Cache des URIs Spotify par nom d'artiste
	cache map[string]string
}

// NewSpotifyService crée un nouveau service Spotify
func NewSpotifyService() *SpotifyService {
	return &SpotifyService{
		cache: make(map[string]string),
	}
}

// GetEmbedURL retourne l'URL d'embed Spotify pour un artiste
func (s *SpotifyService) GetEmbedURL(artistName string) string {
	// Vérifier le cache
	if uri, ok := s.cache[artistName]; ok {
		return uri
	}
	
	// Mapping manuel des artistes connus (basé sur l'API Groupie Tracker)
	// Dans une vraie application, on utiliserait l'API Spotify
	knownArtists := map[string]string{
		"Queen": "spotify:artist:1dfeR4HaWDbWqFHLkxsg1d",
		"Pink Floyd": "spotify:artist:0k17h0D3J5VfsdmQ1iZtE9",
		"Led Zeppelin": "spotify:artist:36QJpDe2go2KgaRleHCDTp",
		"Scorpions": "spotify:artist:27T030eWyCQRmDyuvr1kxY",
		"Soja": "spotify:artist:6yy6iCKGTP5VcsXCLv4Iwd",
		"John Mayer": "spotify:artist:0hEurMDQu99nJRq8pTxO14",
		"Imagine Dragons": "spotify:artist:53XhwfbYqKCa1cC15pYq2q",
		"Linkin Park": "spotify:artist:6XyY86QOPPrYVGvF9ch6wz",
		"Gorillaz": "spotify:artist:3AA28KZvwAUcZuOKwyblJQ",
		"AC/DC": "spotify:artist:711MCceyCBcFnzjGY4Q7Un",
		"Aerosmith": "spotify:artist:7Ey4PD4MYsKc5I2dolUwbH",
		"Eagles": "spotify:artist:0ECwFtbIWEVNwjlrfc6xoL",
		"The Rolling Stones": "spotify:artist:22bE4uQ6baNwSHPVcDxLCe",
		"Metallica": "spotify:artist:2ye2Wgw4gimLv2eAKyk1NB",
		"Red Hot Chili Peppers": "spotify:artist:0L8ExT028jH3ddEcZwqJJ5",
		"Muse": "spotify:artist:12Chz98pHFMPJEknJQMWvI",
		"Coldplay": "spotify:artist:4gzpq5DPGxSnKTe4SA8HAU",
		"Arctic Monkeys": "spotify:artist:7Ln80lUS6He07XvHI8qqHH",
		"Twenty One Pilots": "spotify:artist:3YQKmKGau1PzlVlkL1iodx",
		"Bon Jovi": "spotify:artist:58lV9VcRSjABbAbfWS6skp",
	}
	
	// Chercher dans le mapping
	if uri, ok := knownArtists[artistName]; ok {
		s.cache[artistName] = uri
		return s.formatEmbedURL(uri)
	}
	
	// Pour les artistes non listés, créer une URL de recherche
	searchURL := fmt.Sprintf("https://open.spotify.com/search/%s", url.QueryEscape(artistName))
	s.cache[artistName] = searchURL
	
	return searchURL
}

// formatEmbedURL convertit un URI Spotify en URL d'embed
func (s *SpotifyService) formatEmbedURL(spotifyURI string) string {
	// Convertir spotify:artist:ID en https://open.spotify.com/embed/artist/ID
	parts := strings.Split(spotifyURI, ":")
	if len(parts) == 3 {
		return fmt.Sprintf("https://open.spotify.com/embed/%s/%s", parts[1], parts[2])
	}
	
	return spotifyURI
}

// GetSearchURL retourne une URL de recherche Spotify
func (s *SpotifyService) GetSearchURL(artistName string) string {
	return fmt.Sprintf("https://open.spotify.com/search/%s", url.QueryEscape(artistName))
}
