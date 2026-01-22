package services

import (
	"fmt"
	"groupie-tracker/models"
	"sync"
	"time"
)

// GeocodingPreloader gère le préchargement des géolocalisations
type GeocodingPreloader struct {
	geocoder *GeocodingService
	cache    map[string]*Coordinates
	mu       sync.RWMutex
	loaded   int
	total    int
}

// NewGeocodingPreloader crée un préchargeur
func NewGeocodingPreloader(geocoder *GeocodingService) *GeocodingPreloader {
	return &GeocodingPreloader{
		geocoder: geocoder,
		cache:    make(map[string]*Coordinates),
	}
}

// PreloadAll précharge toutes les géolocalisations pour tous les artistes
func (gp *GeocodingPreloader) PreloadAll(artists []models.Artist, onProgress func(int, int)) error {
	fmt.Println("🌍 Démarrage du préchargement des géolocalisations...")
	startTime := time.Now()

	// Collecter toutes les locations uniques
	uniqueLocations := make(map[string]bool)

	for _, artist := range artists {
		aggregate, err := AggregateArtist(artist)
		if err != nil {
			continue
		}

		for _, location := range aggregate.Locations.Locations {
			uniqueLocations[location] = true
		}
	}

	// Convertir en slice
	locations := make([]string, 0, len(uniqueLocations))
	for loc := range uniqueLocations {
		locations = append(locations, loc)
	}

	gp.total = len(locations)
	fmt.Printf("📍 %d locations uniques à géocoder\n", gp.total)

	// Utiliser un worker pool pour géocoder en parallèle
	numWorkers := 3 // Limiter pour respecter rate limit
	jobs := make(chan string, len(locations))
	results := make(chan bool, len(locations))

	// Lancer les workers
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go gp.worker(w, jobs, results, &wg)
	}

	// Envoyer les jobs
	go func() {
		for _, location := range locations {
			jobs <- location
		}
		close(jobs)
	}()

	// Collecter les résultats
	go func() {
		wg.Wait()
		close(results)
	}()

	// Compter les résultats et notifier la progression
	for range results {
		gp.mu.Lock()
		gp.loaded++
		current := gp.loaded
		gp.mu.Unlock()

		if onProgress != nil {
			onProgress(current, gp.total)
		}

		// Afficher la progression tous les 10%
		if current%(gp.total/10+1) == 0 {
			fmt.Printf("📊 Progression: %d/%d (%.0f%%)\n", current, gp.total, float64(current)*100/float64(gp.total))
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("✅ Préchargement terminé: %d/%d locations en %v\n", gp.loaded, gp.total, duration.Round(time.Second))

	return nil
}

// worker géocode les locations
func (gp *GeocodingPreloader) worker(id int, jobs <-chan string, results chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for location := range jobs {
		// Vérifier le cache d'abord
		if coords, exists := gp.geocoder.GetFromCache(location); exists {
			gp.mu.Lock()
			gp.cache[location] = coords
			gp.mu.Unlock()
			results <- true
			continue
		}

		// Géocoder
		coords, err := gp.geocoder.GeocodeLocation(location)
		if err != nil {
			fmt.Printf("⚠️  Worker %d: Erreur '%s': %v\n", id, location, err)
			results <- false
			continue
		}

		// Sauvegarder dans le cache local
		gp.mu.Lock()
		gp.cache[location] = coords
		gp.mu.Unlock()

		results <- true
	}
}

// GetProgress retourne la progression du chargement
func (gp *GeocodingPreloader) GetProgress() (int, int) {
	gp.mu.RLock()
	defer gp.mu.RUnlock()
	return gp.loaded, gp.total
}

// GetCoordinates retourne les coordonnées d'une location
func (gp *GeocodingPreloader) GetCoordinates(location string) (*Coordinates, bool) {
	gp.mu.RLock()
	defer gp.mu.RUnlock()
	coords, exists := gp.cache[location]
	return coords, exists
}

// GetAllCoordinates retourne toutes les coordonnées
func (gp *GeocodingPreloader) GetAllCoordinates() map[string]*Coordinates {
	gp.mu.RLock()
	defer gp.mu.RUnlock()

	// Copier le cache pour éviter les race conditions
	result := make(map[string]*Coordinates, len(gp.cache))
	for k, v := range gp.cache {
		result[k] = v
	}

	return result
}
