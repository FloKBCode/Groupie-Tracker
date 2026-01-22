package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// FavoritesManager gère les artistes favoris
type FavoritesManager struct {
	mu        sync.RWMutex
	favorites map[int]bool // ID artiste -> favori
	filePath  string
}

// NewFavoritesManager crée un nouveau gestionnaire de favoris
func NewFavoritesManager() *FavoritesManager {
	homeDir, _ := os.UserHomeDir()
	filePath := filepath.Join(homeDir, ".groupie-tracker-favorites.json")
	
	fm := &FavoritesManager{
		favorites: make(map[int]bool),
		filePath:  filePath,
	}
	
	// Charger les favoris existants
	fm.load()
	
	return fm
}

// Toggle ajoute ou retire un artiste des favoris
func (fm *FavoritesManager) Toggle(artistID int) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	fm.favorites[artistID] = !fm.favorites[artistID]
	
	// Sauvegarder
	fm.save()
	
	return fm.favorites[artistID]
}

// IsFavorite vérifie si un artiste est dans les favoris
func (fm *FavoritesManager) IsFavorite(artistID int) bool {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	
	return fm.favorites[artistID]
}

// GetFavorites retourne la liste des IDs favoris
func (fm *FavoritesManager) GetFavorites() []int {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	
	var favs []int
	for id, isFav := range fm.favorites {
		if isFav {
			favs = append(favs, id)
		}
	}
	
	return favs
}

// Count retourne le nombre de favoris
func (fm *FavoritesManager) Count() int {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	
	count := 0
	for _, isFav := range fm.favorites {
		if isFav {
			count++
		}
	}
	
	return count
}

// Clear supprime tous les favoris
func (fm *FavoritesManager) Clear() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	fm.favorites = make(map[int]bool)
	fm.save()
}

// save sauvegarde les favoris dans un fichier
func (fm *FavoritesManager) save() {
	data, err := json.Marshal(fm.favorites)
	if err != nil {
		return
	}
	
	os.WriteFile(fm.filePath, data, 0644)
}

// load charge les favoris depuis le fichier
func (fm *FavoritesManager) load() {
	data, err := os.ReadFile(fm.filePath)
	if err != nil {
		return
	}
	
	json.Unmarshal(data, &fm.favorites)
}
