# 🎵 Groupie Tracker

**Groupie Tracker** est une application de bureau moderne développée en **Go** avec le framework **Fyne**, permettant d'explorer et de visualiser des informations détaillées sur des artistes musicaux, leurs tournées, concerts et emplacements géographiques.

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)
![Fyne](https://img.shields.io/badge/Fyne-v2-00ADD8)
![License](https://img.shields.io/badge/license-MIT-green)

---

## 📋 Table des matières

- [Fonctionnalités](#-fonctionnalités)
- [Technologies](#-technologies)
- [Installation](#-installation)
- [Utilisation](#-utilisation)
- [Architecture](#-architecture)
- [API](#-api)
- - [Difficultés techniques rencontrées](#difficultes-techniques)


---

## ✨ Fonctionnalités

### 🔍 **Recherche Avancée**
- **Recherche intelligente** : Artistes, membres, lieux, dates
- **Recherche par initiales** : Exemple : "fm" → Freddie Mercury
- **Recherche floue** : Tolérance aux fautes de frappe ("qeen" → Queen)
- **Historique de recherche** : Accès rapide aux recherches récentes

### 🎨 **Modes d'affichage**
- **📋 Vue Liste** : Affichage détaillé classique avec séparateurs élégants
- **🖼️ Vue Galerie** : Grille moderne avec images préchargées
- **🗺️ Vue Carte** : Visualisation géographique interactive des concerts avec OpenStreetMap

### ⭐ **Système de Favoris**
- Ajout/suppression rapide d'artistes favoris
- Vue dédiée aux favoris
- Sauvegarde automatique persistante

### 🔧 **Filtres Puissants**
- Filtrage par date de création
- Filtrage par date du premier album
- Filtrage par nombre de membres
- Filtrage par lieux de concert

### 🗺️ **Carte Interactive**
- Géolocalisation automatique des concerts
- Centrage intelligent sur les zones de concerts
- Zoom adaptatif selon la dispersion géographique

### 🎧 **Intégration Spotify**
- Liens directs vers les artistes sur Spotify
- Recherche directe sur Spotify
- Boutons d'écoute rapide

### 📊 **Affichage des Détails**
- Informations générales (création, premier album, membres)
- Liste complète des membres du groupe
- Programme détaillé des concerts avec dates formatées


### ⚡ **Performances**
- Préchargement intelligent des images en arrière-plan
- Cache d'images pour navigation fluide
- Géolocalisation à la demande (charge uniquement si nécessaire)
- Traitement asynchrone des données

---

## 🛠 Technologies

### Backend
- **Go 1.23** : Langage principal
- **net/http** : Client HTTP pour les API
- **encoding/json** : Manipulation JSON

### Frontend
- **Fyne v2** : Framework UI moderne pour Go
- **fyne.io/x/fyne** : Extensions Fyne (carte OpenStreetMap)

### Services
- **Groupie Trackers API** : Source de données artistes
- **Nominatim API** : Géocodage OpenStreetMap
- **Spotify** : Intégration musicale

### Fonctionnalités avancées
- Cache d'images en mémoire
- Géocodage avec mise en cache
- Système de préchargement progressif
- Recherche fuzzy avec algorithme de Levenshtein

---

## 📥 Installation

### Prérequis
- **Go 1.23+** installé ([télécharger Go](https://go.dev/dl/))
- Connexion Internet (pour l'API et la carte)

### Étapes

1. **Cloner le dépôt**
```bash
git clone https://github.com/votre-username/groupie-tracker.git
cd groupie-tracker
```

2. **Installer les dépendances**
```bash
go mod download
```

3. **Compiler et lancer**
```bash
go run main.go
```

Ou compiler un exécutable :
```bash
go build -o groupie-tracker
./groupie-tracker
```

---

## 🚀 Utilisation

### Navigation principale

1. **Recherche** : Tapez dans la barre de recherche en haut
2. **Filtres** : Cliquez sur le bouton "🔧 Filtres" pour affiner les résultats
3. **Modes d'affichage** :
   - 📋 **Liste** : Vue détaillée
   - 🖼️ **Galerie** : Grille avec images
   - 🗺️ **Carte** : Visualisation géographique
4. **Favoris** : Cliquez sur ⭐ pour ajouter/retirer des favoris
5. **Détails** : Cliquez sur un artiste pour voir toutes ses informations

### Vue Carte (Améliorée ✅)

1. Cliquez sur "🗺️ Carte"
2. **Sélection d'artiste en grille** : Visualisation de plusieurs artistes simultanément (au lieu d'une liste scrollable)
3. Cliquez sur "🗺️ Voir carte" pour un artiste
4. La carte se centre automatiquement sur une carte mondiale
5. Cliquez sur "voir" pour zoomer sur un lieu spécifique

---

## 🏗 Architecture

```
groupie-tracker/
├── main.go                 # Point d'entrée
├── api/
│   └── client.go          # Client HTTP pour l'API
├── models/
│   ├── artist.go          # Modèle Artist
│   ├── location.go        # Modèle Location
│   ├── date.go            # Modèle Date
│   ├── relation.go        # Modèle Relation
│   └── artist_data.go     # Agrégation des données
├── services/
│   ├── fetch.go           # Récupération API
│   ├── search.go          # Moteur de recherche
│   ├── filters.go         # Système de filtres
│   ├── geocoding.go       # Géolocalisation
│   ├── favorites.go       # Gestion favoris
│   ├── image_cache.go     # Cache d'images
│   ├── spotify.go         # Intégration Spotify
│   └── utils.go           # Utilitaires
└── ui/
    ├── app.go             # Application principale
    ├── artist_list.go     # Vue liste 
    ├── artist_details.go  # Vue détails 
    ├── map_view.go        # Vue carte 
    ├── search_bar.go      # Barre de recherche
    ├── filters_panel.go   # Panneau de filtres
    └── favorites_view.go  # Vue favoris
```

---

## 🌐 API

### Groupie Trackers API

**Base URL** : `https://groupietrackers.herokuapp.com/api`

#### Endpoints utilisés

| Endpoint | Description |
|----------|-------------|
| `/artists` | Liste complète des artistes |
| `/locations` | Lieux des concerts |
| `/dates` | Dates des concerts |
| `/relation` | Relations dates/lieux |

### Nominatim (OpenStreetMap)

**Base URL** : `https://nominatim.openstreetmap.org`

- Géocodage des adresses
- Rate limit : 1 requête/seconde
- User-Agent requis : `GroupieTracker/1.0`

---

## ⚠️ Difficultés Techniques Rencontrées <a id="difficultes-techniques"></a>

Durant le développement de ce projet, nous avons fait face à de nombreux défis techniques liés à Fyne, au backend, et à l'architecture globale. Voici un récapitulatif détaillé des problèmes et de leurs solutions.

### 🗺️ **Widget Map de Fyne - Problème Majeur**

**Problème** : Le widget `Map` de Fyne (`fyne.io/x/fyne/widget`) est **extrêmement limité** et mal documenté.

#### Tentatives Infructueuses

1. **API inexistante**
   ```go
   // ❌ Ces méthodes n'existent PAS
   mv.mapWidget.SetZoom(10)
   mv.mapWidget.Center(lat, lon)
   mv.mapWidget.Latitude = lat  // ❌ Pas de champs publics accessibles
   ```
   - Erreurs : `SetZoom undefined`, `Center undefined`, `Latitude undefined`
   - La documentation suggère ces méthodes, mais elles ne sont **pas implémentées**

2. **Champs publics inaccessibles**
   ```go
   // ❌ Tentative d'accès direct
   mv.mapWidget.Zoom = 8
   // Erreur: cannot assign to mv.mapWidget.Zoom 
   // (neither addressable nor a map index expression)
   ```

#### Solution Finale

**Abandon du widget Map** et implémentation d'une **carte personnalisée** :
- ✅ Téléchargement direct des tuiles OpenStreetMap (256x256px)
- ✅ Assemblage manuel de 9 tuiles (grille 3x3)
- ✅ Calcul mathématique des coordonnées → pixels
- ✅ Dessin manuel des marqueurs avec `image.RGBA`

```go
// Téléchargement des tuiles OSM
url := fmt.Sprintf("https://tile.openstreetmap.org/%d/%d/%d.png", 
    zoom, tileX, tileY)

// Conversion GPS → pixel
tileX, tileY := latLonToTile(lat, lon, zoom)
px := (dx+1)*256 + 128
py := (dy+1)*256 + 128

// Dessin du marqueur pixel par pixel
for dy := -radius; dy <= radius; dy++ {
    for dx := -radius; dx <= radius; dx++ {
        img.Set(px+dx, py+dy, markerColor)
    }
}
```

**Temps perdu** : ~6 heures de debugging pour finalement tout recoder manuellement.

---

### 📦 **Problèmes avec les Documents (DOCX/PPTX/XLSX)**

**Problème** : Création de fichiers Office complexes avec Go.

#### Difficultés Rencontrées

1. **Bibliothèque `unioffice` limitée**
   - Pas de support natif pour les styles avancés
   - Problèmes de formatage des tableaux
   - Gestion des couleurs incohérente

2. **Tracked Changes dans DOCX**
   ```go
   // ❌ API complexe et peu intuitive
   run.Properties().SetHighlight(wml.ST_HighlightColorYellow)
   // Nécessite de comprendre la structure XML interne
   ```

3. **Formules Excel**
   - Recalcul manuel nécessaire
   - Les formules ne s'auto-évaluent pas
   - Besoin de `f.UpdateLinkedValue()` partout

#### Solutions Appliquées

- ✅ Création de **skills** (guides de bonnes pratiques)
- ✅ Abstraction des opérations complexes
- ✅ Documentation exhaustive des patterns

**Temps investi** : ~10 heures pour maîtriser les APIs.

---

### 🎨 **Interface Utilisateur avec Fyne**

#### 1. **Boutons qui Débordent des Cards**

**Problème** : Les boutons sortaient du cadre des cards dans la vue galerie.

```go
// ❌ Avant : Bouton déborde
container.NewVBox(
    artistImage,
    infoBox,
    detailsButton  // Déborde de la card !
)
```

**Solution** :
```go
// ✅ Après : Padding contrôlé
container.NewVBox(
    artistImage,
    infoBox,
    container.NewPadded(detailsButton)  // Reste dans la card
)
```

**Ajustement supplémentaire** :
- Taille de card : 300x420 → 300x450
- Images : 270x220 → 280x220

---

#### 2. **Centrage de la Grille**

**Problème** : `GridWrap` ne se centre pas automatiquement.

```go
// ❌ Avant : Grille collée à gauche
cards := container.NewGridWrap(fyne.NewSize(300, 450))
return container.NewVScroll(cards)
```

**Solution** :
```go
// ✅ Après : Grille centrée
cards := container.NewGridWithColumns(3)  // 3 colonnes fixes
return container.NewVScroll(container.NewCenter(cards))
```

---

#### 3. **Thème Personnalisé**

**Problème** : Créer un thème sombre cohérent.

**Solution** : Implémentation complète de l'interface `fyne.Theme`
```go
type DarkPurpleTheme struct{}

func (t DarkPurpleTheme) Color(name fyne.ThemeColorName, ...) color.Color {
    switch name {
    case theme.ColorNameForeground:
        return color.RGBA{R: 255, G: 255, B: 255, A: 255}
    case theme.ColorNameBackground:
        return color.RGBA{R: 25, G: 20, B: 35, A: 255}
    // ... 20+ couleurs à définir
    }
}
```

**Temps** : ~2 heures pour tester toutes les couleurs.

---

### 🌐 **Backend - Géolocalisation**

#### Problème 1 : **Rate Limiting de Nominatim**

**Problème** : L'API Nominatim (OpenStreetMap) limite à **1 requête/seconde**.

```
❌ Erreur: 429 Too Many Requests
```

**Solution** :
```go
// Système de cache persistant
cache := make(map[string]*Coordinates)

// Workers avec délai
time.Sleep(1 * time.Second)  // Respect du rate limit

// Préchargement intelligent
go preloadCoordinates()  // En arrière-plan
```

**Résultat** : 189 lieux → ~3 minutes de préchargement initial, puis cache permanent.

---

#### Problème 2 : **Coordonnées Incorrectes**

**Problème** : Certaines villes retournent des coordonnées fausses.

Exemple : "Salem, Germany" → Salem, Oregon, USA

**Solution** :
```go
// Vérification du pays dans la réponse
if coords.Country != expectedCountry {
    // Affiner la recherche
    query = fmt.Sprintf("%s, %s, %s", city, region, country)
}
```

---

#### Problème 3 : **Parsing des Lieux**

**Problème** : Format incohérent des lieux.

```
"new-york_usa"
"las_vegas-usa"
"san francisco, usa"
```

**Solution** :
```go
func ParseLocation(location string) (city, country string) {
    // Nettoyer les underscores et tirets
    clean := strings.ReplaceAll(location, "-", " ")
    clean = strings.ReplaceAll(clean, "_", " ")
    
    // Séparer ville/pays
    parts := strings.Split(clean, ",")
    // ...
}
```

---

### 📊 **Performance et Optimisation**

#### Problème : **Chargement Initial Lent**

**Causes** :
1. Téléchargement de 52 images d'artistes (5-10 MB total)
2. Géocodage de 189 lieux uniques
3. Parsing de 4 endpoints API différents

**Solutions Appliquées** :

1. **Cache Images**
   ```go
   type ImageCache struct {
       cache map[int]image.Image
       mu    sync.RWMutex
   }
   
   // Préchargement asynchrone
   go preloadImages()
   ```

2. **Lazy Loading pour Géolocalisation**
   ```go
   // ❌ Avant : Tout au démarrage (3 min)
   preloadAllCoordinates()
   
   // ✅ Après : À la demande par artiste (5-10 sec)
   go loadCoordinatesOnDemand(artistID)
   ```

3. **Écran d'Accueil**
   ```go
   // Masquer le chargement avec un splash screen
   splash := NewSplashScreen(onStart)
   ```

**Amélioration** : Démarrage 3 min → 2-3 secondes

---

### 🔄 **Concurrence et Thread Safety**

**Problème** : Accès concurrents au cache.

```go
// ❌ Race condition
cache[key] = value  // Plusieurs goroutines en même temps
```

**Solution** :
```go
type SafeCache struct {
    cache map[string]*Data
    mu    sync.RWMutex
}

func (c *SafeCache) Set(key string, value *Data) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[key] = value
}
```

**Tests** : `go test -race` pour détecter les races.

---

### 📱 **Compatibilité Multi-Plateforme**

#### Problèmes Rencontrés

1. **Windows** : Chemins avec `\` au lieu de `/`
   ```go
   // Solution
   filepath.Join("dir", "file")  // Au lieu de "dir/file"
   ```

2. **macOS** : Permissions de fichiers
   ```go
   os.MkdirAll(dir, 0755)  // Au lieu de 0777
   ```

3. **Linux** : Dépendances système pour Fyne
   ```bash
   # Ubuntu/Debian
   sudo apt install libgl1-mesa-dev xorg-dev
   ```

---

### 🐛 **Bugs Subtils Résolus**

#### 1. **Closure dans les Boucles**

```go
// ❌ Bug classique Go
for _, artist := range artists {
    button.OnTapped = func() {
        showDetails(artist.ID)  // Toujours le dernier !
    }
}

// ✅ Solution
for i := range artists {
    artist := artists[i]  // Copie locale
    button.OnTapped = func() {
        showDetails(artist.ID)
    }
}
```

#### 2. **Rafraîchissement de l'UI**

```go
// ❌ Oublier de rafraîchir
mv.mapContainer.Objects = []fyne.CanvasObject{newMap}
// La carte ne s'affiche pas !

// ✅ Toujours rafraîchir
mv.mapContainer.Objects = []fyne.CanvasObject{newMap}
mv.mapContainer.Refresh()  // Crucial !
```

---

### 📈 **Récapitulatif des Temps Investis**

| Problème | Temps Débogage | Solution |
|----------|----------------|----------|
| Widget Map Fyne | ~6h | Implémentation custom avec tuiles OSM |
| Documents Office | ~10h | Création de skills + abstraction |
| Géolocalisation | ~4h | Cache + workers avec rate limiting |
| UI Fyne (boutons, centrage) | ~3h | Padding + containers spécifiques |
| Thème personnalisé | ~2h | Implémentation complète de fyne.Theme |
| Performance | ~3h | Lazy loading + préchargement async |
| Thread safety | ~2h | Mutexes + tests -race |
| **TOTAL** | **~30h** | Application fonctionnelle et optimisée |

---

### 💡 **Leçons Apprises**

1. **Fyne** : Framework puissant mais documentation parfois obsolète
   - Toujours vérifier le code source
   - Ne pas faire confiance aux exemples sur Internet

2. **Go** : Excellent pour la concurrence
   - Mutexes indispensables pour les caches
   - `go test -race` est votre ami

3. **APIs Externes** : Toujours prévoir des limites
   - Rate limiting
   - Timeouts
   - Fallbacks

4. **Performance** : Le lazy loading change tout
   - Ne jamais tout charger au démarrage
   - Cache intelligent > Force brute

---

## 👨‍💻 Auteur

Développé par Florence Kore-Belle, Theo Bouaziz, Sasha Domin, Mariam Keita

---

**Bon tracking ! 🎸🎵🗺️**


