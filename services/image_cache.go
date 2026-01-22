package services

import (
	"fmt"
	"groupie-tracker/models"
	"image"
	"io"
	"net/http"
	"sync"
	"time"

	_ "image/jpeg"
	_ "image/png"
)

// ImageCache gère le cache des images d'artistes
type ImageCache struct {
	mu     sync.RWMutex
	cache  map[int]image.Image // ID artiste -> Image
	errors map[int]error       // ID artiste -> Erreur de chargement
}

// NewImageCache crée un nouveau cache d'images
func NewImageCache() *ImageCache {
	return &ImageCache{
		cache:  make(map[int]image.Image),
		errors: make(map[int]error),
	}
}

// PreloadImages précharge toutes les images des artistes
func (ic *ImageCache) PreloadImages(artists []models.Artist, progressCallback func(current, total int)) error {
	total := len(artists)
	
	// Limiter le nombre de goroutines simultanées
	semaphore := make(chan struct{}, 5)
	var wg sync.WaitGroup
	
	for i, artist := range artists {
		wg.Add(1)
		
		go func(idx int, a models.Artist) {
			defer wg.Done()
			
			semaphore <- struct{}{} // Acquérir
			defer func() { <-semaphore }() // Libérer
			
			// Charger l'image
			img, err := ic.loadImage(a.Image)
			
			ic.mu.Lock()
			if err != nil {
				ic.errors[a.ID] = err
			} else {
				ic.cache[a.ID] = img
			}
			ic.mu.Unlock()
			
			// Callback de progression
			if progressCallback != nil {
				progressCallback(idx+1, total)
			}
		}(i, artist)
	}
	
	wg.Wait()
	return nil
}

// GetImage retourne une image depuis le cache
func (ic *ImageCache) GetImage(artistID int) (image.Image, bool) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	
	img, ok := ic.cache[artistID]
	return img, ok
}

// HasError vérifie si le chargement d'une image a échoué
func (ic *ImageCache) HasError(artistID int) bool {
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	
	_, hasErr := ic.errors[artistID]
	return hasErr
}

// GetProgress retourne la progression du chargement
func (ic *ImageCache) GetProgress() (loaded, total int) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	
	loaded = len(ic.cache) + len(ic.errors)
	return loaded, loaded
}

// loadImage charge une image depuis une URL
func (ic *ImageCache) loadImage(url string) (image.Image, error) {
	if url == "" {
		return nil, fmt.Errorf("URL vide")
	}
	
	// Client HTTP avec timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur HTTP: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status HTTP: %d", resp.StatusCode)
	}
	
	// Lire et décoder l'image
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture: %w", err)
	}
	
	img, _, err := image.Decode(io.NopCloser(io.Reader(newBytesReader(data))))
	if err != nil {
		return nil, fmt.Errorf("erreur décodage: %w", err)
	}
	
	return img, nil
}

// bytesReader est un wrapper pour les bytes
type bytesReader struct {
	data []byte
	pos  int
}

func newBytesReader(data []byte) *bytesReader {
	return &bytesReader{data: data, pos: 0}
}

func (br *bytesReader) Read(p []byte) (n int, err error) {
	if br.pos >= len(br.data) {
		return 0, io.EOF
	}
	n = copy(p, br.data[br.pos:])
	br.pos += n
	return n, nil
}
