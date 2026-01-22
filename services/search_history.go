package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// SearchHistoryEntry représente une entrée dans l'historique
type SearchHistoryEntry struct {
	Query     string    `json:"query"`
	Timestamp time.Time `json:"timestamp"`
	ResultID  int       `json:"result_id"` // ID de l'artiste sélectionné
}

// SearchHistory gère l'historique des recherches
type SearchHistory struct {
	entries  []SearchHistoryEntry
	maxSize  int
	filePath string
}

// NewSearchHistory crée un nouvel historique de recherche
func NewSearchHistory(maxSize int) *SearchHistory {
	// Déterminer le chemin du fichier d'historique
	homeDir, _ := os.UserHomeDir()
	filePath := filepath.Join(homeDir, ".groupie-tracker", "search_history.json")

	sh := &SearchHistory{
		entries:  []SearchHistoryEntry{},
		maxSize:  maxSize,
		filePath: filePath,
	}

	// Charger l'historique existant
	sh.Load()

	return sh
}

// Add ajoute une recherche à l'historique
func (sh *SearchHistory) Add(query string, resultID int) {
	// Éviter les doublons consécutifs
	if len(sh.entries) > 0 {
		last := sh.entries[len(sh.entries)-1]
		if last.Query == query && last.ResultID == resultID {
			return
		}
	}

	entry := SearchHistoryEntry{
		Query:     query,
		Timestamp: time.Now(),
		ResultID:  resultID,
	}

	sh.entries = append(sh.entries, entry)

	// Limiter la taille
	if len(sh.entries) > sh.maxSize {
		sh.entries = sh.entries[len(sh.entries)-sh.maxSize:]
	}

	// Sauvegarder
	sh.Save()
}

// GetRecent retourne les N dernières recherches
func (sh *SearchHistory) GetRecent(n int) []SearchHistoryEntry {
	if n > len(sh.entries) {
		n = len(sh.entries)
	}

	// Retourner dans l'ordre inverse (plus récent en premier)
	recent := make([]SearchHistoryEntry, n)
	for i := 0; i < n; i++ {
		recent[i] = sh.entries[len(sh.entries)-1-i]
	}

	return recent
}

// GetAll retourne tout l'historique
func (sh *SearchHistory) GetAll() []SearchHistoryEntry {
	return sh.entries
}

// GetMostFrequent retourne les queries les plus fréquentes
func (sh *SearchHistory) GetMostFrequent(n int) []string {
	// Compter les occurrences
	frequency := make(map[string]int)

	for _, entry := range sh.entries {
		frequency[entry.Query]++
	}

	// Créer une liste triée
	type queryFreq struct {
		query string
		count int
	}

	freqList := []queryFreq{}
	for query, count := range frequency {
		freqList = append(freqList, queryFreq{query, count})
	}

	// Trier par fréquence (bubble sort simple)
	for i := 0; i < len(freqList)-1; i++ {
		for j := 0; j < len(freqList)-i-1; j++ {
			if freqList[j].count < freqList[j+1].count {
				freqList[j], freqList[j+1] = freqList[j+1], freqList[j]
			}
		}
	}

	// Retourner les N plus fréquentes
	if n > len(freqList) {
		n = len(freqList)
	}

	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = freqList[i].query
	}

	return result
}

// Clear efface tout l'historique
func (sh *SearchHistory) Clear() {
	sh.entries = []SearchHistoryEntry{}
	sh.Save()
}

// Save sauvegarde l'historique sur disque
func (sh *SearchHistory) Save() error {
	// Créer le répertoire si nécessaire
	dir := filepath.Dir(sh.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Encoder en JSON
	data, err := json.MarshalIndent(sh.entries, "", "  ")
	if err != nil {
		return err
	}

	// Écrire dans le fichier
	return os.WriteFile(sh.filePath, data, 0644)
}

// Load charge l'historique depuis le disque
func (sh *SearchHistory) Load() error {
	// Vérifier si le fichier existe
	if _, err := os.Stat(sh.filePath); os.IsNotExist(err) {
		return nil // Pas d'erreur, juste pas d'historique
	}

	// Lire le fichier
	data, err := os.ReadFile(sh.filePath)
	if err != nil {
		return err
	}

	// Décoder le JSON
	return json.Unmarshal(data, &sh.entries)
}

// GetSuggestions retourne des suggestions basées sur l'historique
func (sh *SearchHistory) GetSuggestions(currentQuery string) []string {
	if currentQuery == "" {
		// Si pas de query, retourner les recherches récentes
		recent := sh.GetRecent(5)
		suggestions := make([]string, len(recent))
		for i, entry := range recent {
			suggestions[i] = entry.Query
		}
		return suggestions
	}

	// Sinon, chercher dans l'historique les queries qui commencent pareil
	suggestions := []string{}
	seen := make(map[string]bool)

	// Parcourir l'historique à l'envers (plus récent en premier)
	for i := len(sh.entries) - 1; i >= 0; i-- {
		entry := sh.entries[i]

		// Si la query commence par currentQuery
		if len(entry.Query) >= len(currentQuery) &&
			entry.Query[:len(currentQuery)] == currentQuery &&
			!seen[entry.Query] {
			suggestions = append(suggestions, entry.Query)
			seen[entry.Query] = true

			// Limiter à 5 suggestions
			if len(suggestions) >= 5 {
				break
			}
		}
	}

	return suggestions
}

// GetStatistics retourne des statistiques sur l'historique
func (sh *SearchHistory) GetStatistics() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["total_searches"] = len(sh.entries)

	if len(sh.entries) > 0 {
		stats["oldest"] = sh.entries[0].Timestamp
		stats["newest"] = sh.entries[len(sh.entries)-1].Timestamp
		stats["most_frequent"] = sh.GetMostFrequent(5)
	}

	return stats
}

// RemoveOldEntries supprime les entrées plus anciennes qu'une durée
func (sh *SearchHistory) RemoveOldEntries(duration time.Duration) {
	cutoff := time.Now().Add(-duration)

	filtered := []SearchHistoryEntry{}
	for _, entry := range sh.entries {
		if entry.Timestamp.After(cutoff) {
			filtered = append(filtered, entry)
		}
	}

	sh.entries = filtered
	sh.Save()
}
