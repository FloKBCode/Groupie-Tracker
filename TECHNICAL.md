# 🔧 Documentation Technique - Groupie Tracker

Cette documentation technique décrit l'architecture, les patterns utilisés et les détails d'implémentation.

---

## 📐 Architecture Globale

### Pattern MVC Adapté

```
┌─────────────┐
│   UI Layer  │ ← Interface utilisateur (Fyne)
└──────┬──────┘
       │
┌──────▼──────┐
│  Services   │ ← Logique métier
└──────┬──────┘
       │
┌──────▼──────┐
│   Models    │ ← Structures de données
└──────┬──────┘
       │
┌──────▼──────┐
│     API     │ ← Communication externe
└─────────────┘
```

---

## 🗂 Structure des Packages

### 📦 **Package `models`**

Définit les structures de données de l'application.

#### `Artist`
```go
type Artist struct {
    ID           int      `json:"id"`
    Image        string   `json:"image"`
    Name         string   `json:"name"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`
    FirstAlbum   string   `json:"firstAlbum"`
}
```

#### `Location`
```go
type Location struct {
    ID        int      `json:"id"`
    Locations []string `json:"locations"` // Format: "city-country"
    Dates     string   `json:"dates"`
}
```

#### `Date`
```go
type Date struct {
    ID    int      `json:"id"`
    Dates []string `json:"dates"` // Format: "DD-MM-YYYY"
}
```

#### `Relation`
```go
type Relation struct {
    ID             int                 `json:"id"`
    DatesLocations map[string][]string `json:"datesLocations"`
}
```

#### `ArtistAggregate`
```go
type ArtistAggregate struct {
    Artist    Artist
    Locations Location
    Dates     Date
    Relation  Relation
}
```

---

### 🔧 **Package `services`**

Contient toute la logique métier.

#### 1. **Fetch Service** (`fetch.go`)

Gestion des appels API.

```go
func GetArtists() ([]Artist, error)
func GetLocations() ([]Location, error)
func GetDates() ([]Date, error)
func GetRelation() ([]Relation, error)
```

**Caractéristiques** :
- Client HTTP avec timeout de 10s
- Retry automatique en cas d'erreur
- Décodage JSON optimisé

#### 2. **Search Engine** (`search.go`)

Moteur de recherche multicritère.

```go
type SearchEngine struct {
    artists      []Artist
    artistsData  map[int]ArtistAggregate
    searchIndex  *SearchIndex
}

func (se *SearchEngine) Search(query string) []SearchResult
```

**Algorithmes** :
- **Recherche exacte** : Match direct sur noms
- **Recherche par initiales** : "fm" → "Freddie Mercury"
- **Recherche floue** : Distance de Levenshtein < 3
- **Recherche membres** : Dans les listes de membres
- **Recherche lieux** : Parsing "city-country"
- **Recherche dates** : Format flexible

#### 3. **Filter Engine** (`filters.go`)

Système de filtrage avancé.

```go
type FilterCriteria struct {
    CreationDateMin  *int
    CreationDateMax  *int
    FirstAlbumMin    *time.Time
    FirstAlbumMax    *time.Time
    MembersMin       *int
    MembersMax       *int
    Locations        []string
}

func (fe *FilterEngine) ApplyFilters(criteria *FilterCriteria) []Artist
```

**Optimisations** :
- Filtrage paresseux (lazy evaluation)
- Cache des résultats intermédiaires
- Préchargement des données agrégées

#### 4. **Geocoding Service** (`geocoding.go`)

Service de géolocalisation avec cache.

```go
type GeocodingService struct {
    cache      map[string]*Coordinates
    client     *http.Client
    rateLimiter *time.Ticker
}

func (gs *GeocodingService) Geocode(location string) (*Coordinates, error)
```

**Caractéristiques** :
- Cache en mémoire (évite les appels répétés)
- Rate limiting : 1 req/s (respect API Nominatim)
- Parsing intelligent "city-country"
- User-Agent personnalisé

#### 5. **Image Cache** (`image_cache.go`)

Cache d'images en mémoire.

```go
type ImageCache struct {
    cache  map[int]image.Image
    mu     sync.RWMutex
}

func (ic *ImageCache) PreloadImages(artists []Artist, progress func(int, int)) error
```

**Optimisations** :
- Mutex pour accès concurrent
- Préchargement asynchrone
- Callback de progression
- Libération mémoire si erreur

#### 6. **Favorites Manager** (`favorites.go`)

Gestion des favoris avec persistance.

```go
type FavoritesManager struct {
    favorites map[int]bool
    mu        sync.RWMutex
    filepath  string
}

func (fm *FavoritesManager) AddFavorite(artistID int) error
func (fm *FavoritesManager) RemoveFavorite(artistID int) error
```

**Caractéristiques** :
- Sauvegarde JSON automatique
- Thread-safe avec RWMutex
- Chargement au démarrage

---

### 🎨 **Package `ui`**

Composants d'interface utilisateur.

#### 1. **App** (`app.go`)

Point d'entrée de l'UI.

```go
type App struct {
    fyneApp         fyne.App
    mainWindow      fyne.Window
    currentView     fyne.CanvasObject
    favoritesManager *services.FavoritesManager
    imageCache       *services.ImageCache
}

func NewApp() *App
func (a *App) Run()
```

#### 2. **Artist List View** (`artist_list.go`)

Vue principale avec 3 modes d'affichage.

```go
type ArtistListView struct {
    Container       fyne.CanvasObject
    allArtists      []models.Artist
    filteredArtists []models.Artist
    viewMode        ViewMode // List, Gallery, Map
}

func (v *ArtistListView) switchView(mode ViewMode)
```

**Modes** :
- **Liste** : `widget.List` avec template personnalisé
- **Galerie** : `container.NewGridWrap` avec cards
- **Carte** : Grille de sélection → MapView

#### 3. **Map View** (`map_view.go`) 

Vue carte interactive.

```go
type MapView struct {
    Container   fyne.CanvasObject
    geocoder    *services.GeocodingService
    coordinates map[string]*services.Coordinates
    mapWidget   *xwidget.Map
}

func (mv *MapView) loadCoordinates()
func (mv *MapView) adjustZoom(centerLat, centerLon float64)
```

**Améliorations** :
- ✅ Centrage avec `Center(lat, lon)` au lieu de `Move()`
- ✅ Zoom adaptatif selon dispersion
- ✅ Calcul du centre moyen
- ✅ Boutons "Centrer" fonctionnels

#### 4. **Artist Details View** (`artist_details.go`) 

Vue détaillée d'un artiste.

```go
type ArtistDetailsView struct {
    Container        fyne.CanvasObject
    aggregate        models.ArtistAggregate
    favoritesManager *services.FavoritesManager
    spotifyService   *services.SpotifyService
}
```

**Sections** :
- Header avec image
- Spotify integration
- Infos générales (Card)
- Membres (Card)
- Lieux (Grille 4 colonnes) 
- Dates (Liste spacieuse) 
- Programme détaillé (Cards par lieu) 

---

## 🔄 Flux de Données

### Chargement Initial

```
1. main.go
   ↓
2. ui.NewApp()
   ↓
3. services.GetArtists() → API Call
   ↓
4. ArtistListView.preload()
   ├→ SearchEngine.LoadAggregateData()
   ├→ FilterEngine.LoadAggregateData()
   └→ ImageCache.PreloadImages()
```

### Recherche

```
1. SearchBar.OnChanged
   ↓
2. SearchEngine.Search(query)
   ├→ Exact Match
   ├→ Initials Search
   ├→ Fuzzy Search
   ├→ Members Search
   ├→ Locations Search
   └→ Dates Search
   ↓
3. ArtistListView.Update(results)
```

### Filtrage

```
1. FiltersPanel.OnApply
   ↓
2. FilterEngine.ApplyFilters(criteria)
   ├→ Creation Date Filter
   ├→ First Album Filter
   ├→ Members Count Filter
   └→ Locations Filter
   ↓
3. ArtistListView.refreshCurrentView()
```

---

## ⚡ Optimisations Performances

### 1. **Préchargement Intelligent**

```go
// Chargement asynchrone avec contexte
go view.preload()

// Préchargement progressif
for i, artist := range v.allArtists {
    select {
    case <-v.ctx.Done():
        return // Arrêt propre
    default:
        v.searchEngine.LoadAggregateData(artist.ID)
    }
}
```

### 2. **Cache Multi-niveaux**

- **Image Cache** : Évite téléchargements répétés
- **Geo Cache** : Évite géocodage répété
- **Aggregate Cache** : Données précalculées

### 3. **Lazy Loading**

- Géolocalisation : Chargée uniquement si vue carte activée
- Images : Préchargées en arrière-plan sans bloquer UI

### 4. **Concurrent Access**

```go
type ImageCache struct {
    cache map[int]image.Image
    mu    sync.RWMutex // Read/Write Mutex
}

func (ic *ImageCache) GetImage(id int) (image.Image, bool) {
    ic.mu.RLock()         // Lock en lecture
    defer ic.mu.RUnlock()
    img, ok := ic.cache[id]
    return img, ok
}
```

---

## 🧪 Tests

### Tests Unitaires

```bash
# Services
go test ./services -v

# Models
go test ./models -v
```

### Couverture de Code

```bash
go test -cover ./...
```

**Fichiers de tests** :
- `services/filters_test.go`
- `services/geocoding_test.go`
- `services/search_test.go`
- `services/utils_test.go`
- `models/models_test.go`

---

## 🔐 Sécurité

### Rate Limiting

```go
// Geocoding : 1 req/s
rateLimiter := time.NewTicker(1 * time.Second)
```

### User-Agent

```go
req.Header.Set("User-Agent", "GroupieTracker/1.0")
```

### Timeouts

```go
client := &http.Client{
    Timeout: 10 * time.Second,
}
```

---

## 📊 Patterns Utilisés

### 1. **Singleton** (Services)
- GeocodingService
- ImageCache
- FavoritesManager

### 2. **Observer** (UI Updates)
- SearchBar → ArtistListView
- FiltersPanel → ArtistListView

### 3. **Strategy** (View Modes)
- ListMode
- GalleryMode
- MapMode

### 4. **Factory** (View Creation)
```go
func NewArtistListView(...) *ArtistListView
func NewMapView(...) *MapView
```

---

## 🚀 Points d'Extension

### Ajouter un nouveau filtre

1. Ajouter le champ dans `FilterCriteria`
2. Implémenter la logique dans `ApplyFilters()`
3. Ajouter l'UI dans `FiltersPanel`

### Ajouter un nouveau mode d'affichage

1. Ajouter une constante dans `ViewMode`
2. Implémenter `create...View()`
3. Gérer dans `switchView()`

### Ajouter un nouveau service

1. Créer le fichier dans `services/`
2. Définir l'interface
3. Implémenter avec cache si nécessaire
4. Intégrer dans `ArtistListView`

---

## 📝 Conventions de Code

### Nommage

- **Packages** : lowercase, singular
- **Types** : PascalCase
- **Fonctions publiques** : PascalCase
- **Fonctions privées** : camelCase
- **Constantes** : PascalCase

### Commentaires

```go
// FunctionName fait quelque chose
// Elle retourne une erreur si...
func FunctionName() error {
    // ...
}
```

### Gestion des erreurs

```go
if err != nil {
    return fmt.Errorf("contexte: %w", err)
}
```

---

## 🔗 Dépendances Externes

### Fyne Framework
```
fyne.io/fyne/v2@v2.5.2
fyne.io/x/fyne@v0.0.0...
```

### APIs
- **Groupie Trackers API** : https://groupietrackers.herokuapp.com
- **Nominatim** : https://nominatim.openstreetmap.org

