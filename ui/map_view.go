package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"math"
	"net/http"

	"groupie-tracker/models"
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// MapView représente une vue carte avec vraies tuiles OpenStreetMap
type MapView struct {
	Container   fyne.CanvasObject
	geocoder    *services.GeocodingService
	artistData  models.ArtistAggregate
	coordinates map[string]*services.Coordinates

	mapContainer    *fyne.Container
	selectedLocation string
}

// NewMapView crée une vue carte avec chargement à la demande
func NewMapView(artistData models.ArtistAggregate, geocoder *services.GeocodingService) *MapView {
	mv := &MapView{
		geocoder:    geocoder,
		artistData:  artistData,
		coordinates: make(map[string]*services.Coordinates),
	}

	mv.buildUI()
	
	// Charger les coordonnées à la demande pour cet artiste
	go mv.loadCoordinatesOnDemand()

	return mv
}

// buildUI construit l'interface avec carte et liste des lieux
func (mv *MapView) buildUI() {
	// Titre simplifié
	title := widget.NewLabelWithStyle(
		fmt.Sprintf("Concerts de %s", mv.artistData.Artist.Name),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Container pour la carte
	mv.mapContainer = container.NewMax()
	mv.showDefaultMap()

	// Statistiques
	statsBox := container.NewVBox(
		widget.NewLabelWithStyle("Statistiques", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("%d lieux de concert", len(mv.artistData.Locations.Locations))),
		widget.NewLabel(fmt.Sprintf("%d dates programmées", len(mv.artistData.Dates.Dates))),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Navigation", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Cliquez sur un lieu pour le centrer"),
		widget.NewLabel("Carte OpenStreetMap dynamique"),
	)

	// Liste des lieux
	locationsList := mv.createLocationsList()

	// Panneau de contrôle
	controlPanel := container.NewVBox(
		statsBox,
		widget.NewSeparator(),
		widget.NewLabel("Lieux de concert"),
		locationsList,
	)

	// Layout principal
	split := container.NewHSplit(
		container.NewBorder(title, nil, nil, nil, mv.mapContainer),
		container.NewScroll(controlPanel),
	)
	split.SetOffset(0.72)

	mv.Container = split
}

// showDefaultMap affiche une carte du monde par défaut
func (mv *MapView) showDefaultMap() {
	centerLat, centerLon, zoom := 20.0, 0.0, 2

	mapImg := mv.renderMap(centerLat, centerLon, zoom, nil)
	
	canvasImg := canvas.NewImageFromImage(mapImg)
	canvasImg.FillMode = canvas.ImageFillContain
	canvasImg.SetMinSize(fyne.NewSize(700, 700))

	infoLabel := widget.NewLabel("Vue mondiale")
	infoLabel.Alignment = fyne.TextAlignCenter
	infoLabel.TextStyle = fyne.TextStyle{Bold: true}

	mv.mapContainer.Objects = []fyne.CanvasObject{
		container.NewBorder(
			container.NewPadded(infoLabel),
			nil, nil, nil,
			container.NewCenter(canvasImg),
		),
	}
}

// showLocationMap affiche la carte centrée sur un lieu
func (mv *MapView) showLocationMap(city, country string, coords *services.Coordinates) {
	zoom := 8
	
	mapImg := mv.renderMap(coords.Latitude, coords.Longitude, zoom, coords)
	
	canvasImg := canvas.NewImageFromImage(mapImg)
	canvasImg.FillMode = canvas.ImageFillContain
	canvasImg.SetMinSize(fyne.NewSize(700, 700))

	// Info simplifiée - seulement le lieu
	infoText := fmt.Sprintf("%s, %s", city, country)

	infoLabel := widget.NewLabel(infoText)
	infoLabel.Alignment = fyne.TextAlignCenter
	infoLabel.TextStyle = fyne.TextStyle{Bold: true}

	mv.mapContainer.Objects = []fyne.CanvasObject{
		container.NewBorder(
			container.NewPadded(infoLabel),
			nil, nil, nil,
			container.NewCenter(canvasImg),
		),
	}
	mv.mapContainer.Refresh()
	
	fmt.Printf("🎯 Carte affichée pour %s, %s (%.4f, %.4f)\n", city, country, coords.Latitude, coords.Longitude)
}

// renderMap génère l'image de la carte avec tuiles OSM
func (mv *MapView) renderMap(centerLat, centerLon float64, zoom int, selectedCoords *services.Coordinates) image.Image {
	centerTileX, centerTileY := latLonToTile(centerLat, centerLon, zoom)

	mapWidth := 768
	mapHeight := 768

	mapImg := image.NewRGBA(image.Rect(0, 0, mapWidth, mapHeight))
	draw.Draw(mapImg, mapImg.Bounds(), &image.Uniform{color.RGBA{200, 220, 240, 255}}, image.Point{}, draw.Src)

	tilesLoaded := 0

	// Télécharger 9 tuiles (3x3)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			tx := centerTileX + dx
			ty := centerTileY + dy

			url := fmt.Sprintf("https://tile.openstreetmap.org/%d/%d/%d.png", zoom, tx, ty)

			tile := downloadImage(url)
			if tile != nil {
				destX := (dx + 1) * 256
				destY := (dy + 1) * 256
				destRect := image.Rect(destX, destY, destX+256, destY+256)

				draw.Draw(mapImg, destRect, tile, image.Point{}, draw.Src)
				tilesLoaded++
			}
		}
	}

	// Dessiner tous les concerts sur la carte
	for _, coords := range mv.coordinates {
		if coords == nil {
			continue
		}
		
		px, py := latLonToPixel(coords.Latitude, coords.Longitude, zoom, centerTileX, centerTileY)
		
		// Vérifier si le point est dans l'image
		if px >= 0 && px < mapWidth && py >= 0 && py < mapHeight {
			drawMarker(mapImg, px, py, false)
		}
	}

	// Marqueur du lieu sélectionné (au centre, plus gros)
	if selectedCoords != nil {
		markerX := mapWidth / 2
		markerY := mapHeight / 2
		drawMarker(mapImg, markerX, markerY, true)
	}

	return mapImg
}

// loadCoordinatesOnDemand charge les coordonnées à la demande pour cet artiste
func (mv *MapView) loadCoordinatesOnDemand() {
	loaded := 0
	total := len(mv.artistData.Locations.Locations)

	fmt.Printf("🌍 Chargement des coordonnées pour %s (%d lieux)...\n", mv.artistData.Artist.Name, total)

	for i, location := range mv.artistData.Locations.Locations {
		// Vérifier si déjà en cache
		if coords, exists := mv.geocoder.GetFromCache(location); exists {
			mv.coordinates[location] = coords
			loaded++
			if i%5 == 0 || i == total-1 {
				fmt.Printf("✓ Cache: %d/%d\n", i+1, total)
			}
		} else {
			// Géocoder à la demande
			coords, err := mv.geocoder.Geocode(location)
			if err == nil && coords != nil {
				mv.coordinates[location] = coords
				loaded++
				fmt.Printf("✅ Géocodé: %s (%.4f, %.4f) [%d/%d]\n", location, coords.Latitude, coords.Longitude, i+1, total)
			} else {
				fmt.Printf("❌ Échec: %s [%d/%d]\n", location, i+1, total)
			}
		}
	}

	if loaded == 0 {
		fmt.Printf("⚠️ Aucune coordonnée disponible pour %s\n", mv.artistData.Artist.Name)
	} else {
		fmt.Printf("✅ %d/%d lieux géocodés pour %s\n", loaded, total, mv.artistData.Artist.Name)
	}
}

// createLocationsList crée la liste des lieux avec boutons
func (mv *MapView) createLocationsList() *widget.List {
	return widget.NewList(
		func() int {
			return len(mv.artistData.Locations.Locations)
		},
		func() fyne.CanvasObject {
			icon := widget.NewLabel("•")
			name := widget.NewLabel("")
			name.TextStyle = fyne.TextStyle{Bold: true}
			status := widget.NewLabel("")
			status.TextStyle = fyne.TextStyle{Italic: true}
			viewBtn := widget.NewButton("Voir", nil)
			viewBtn.Importance = widget.LowImportance

			return container.NewVBox(
				container.NewHBox(icon, name),
				status,
				viewBtn,
				widget.NewSeparator(),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(mv.artistData.Locations.Locations) {
				return
			}

			location := mv.artistData.Locations.Locations[id]
			city, country := services.ParseLocation(location)

			vbox := obj.(*fyne.Container)
			hbox := vbox.Objects[0].(*fyne.Container)

			nameLabel := hbox.Objects[1].(*widget.Label)
			statusLabel := vbox.Objects[1].(*widget.Label)
			viewBtn := vbox.Objects[2].(*widget.Button)

			nameLabel.SetText(fmt.Sprintf("%s, %s", city, country))

			if coords, exists := mv.coordinates[location]; exists {
				statusLabel.SetText(
					fmt.Sprintf("✓ %.4f°, %.4f°", coords.Latitude, coords.Longitude),
				)

				viewBtn.Enable()
				viewBtn.OnTapped = func() {
					mv.showLocationMap(city, country, coords)
				}
			} else {
				statusLabel.SetText("⏳ Chargement en cours...")
				viewBtn.Disable()
			}
		},
	)
}

// drawMarker dessine un marqueur sur la carte
func drawMarker(img *image.RGBA, x, y int, isSelected bool) {
	markerColor := color.RGBA{255, 0, 0, 255}
	outlineColor := color.RGBA{255, 255, 255, 255}
	shadowColor := color.RGBA{0, 0, 0, 100}
	centerColor := color.RGBA{255, 255, 255, 255}

	var shadowRadius, outlineRadius, markerRadius, centerRadius int
	
	if isSelected {
		// Marqueur plus gros pour le lieu sélectionné
		shadowRadius = 18
		outlineRadius = 16
		markerRadius = 14
		centerRadius = 6
	} else {
		// Marqueur normal pour les autres lieux
		shadowRadius = 10
		outlineRadius = 9
		markerRadius = 8
		centerRadius = 3
	}

	// Ombre
	for dy := -shadowRadius; dy <= shadowRadius; dy++ {
		for dx := -shadowRadius; dx <= shadowRadius; dx++ {
			if dx*dx+dy*dy <= shadowRadius*shadowRadius {
				px, py := x+dx+2, y+dy+2
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, shadowColor)
				}
			}
		}
	}

	// Contour blanc
	for dy := -outlineRadius; dy <= outlineRadius; dy++ {
		for dx := -outlineRadius; dx <= outlineRadius; dx++ {
			if dx*dx+dy*dy <= outlineRadius*outlineRadius {
				px, py := x+dx, y+dy
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, outlineColor)
				}
			}
		}
	}

	// Cercle rouge
	for dy := -markerRadius; dy <= markerRadius; dy++ {
		for dx := -markerRadius; dx <= markerRadius; dx++ {
			if dx*dx+dy*dy <= markerRadius*markerRadius {
				px, py := x+dx, y+dy
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, markerColor)
				}
			}
		}
	}

	// Point blanc au milieu
	for dy := -centerRadius; dy <= centerRadius; dy++ {
		for dx := -centerRadius; dx <= centerRadius; dx++ {
			if dx*dx+dy*dy <= centerRadius*centerRadius {
				px, py := x+dx, y+dy
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, centerColor)
				}
			}
		}
	}

	// Croix noire pour le marqueur sélectionné
	if isSelected {
		crossSize := 10
		crossThickness := 2
		crossColor := color.RGBA{0, 0, 0, 255}

		for dx := -crossSize; dx <= crossSize; dx++ {
			for t := -crossThickness; t <= crossThickness; t++ {
				px, py := x+dx, y+t
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, crossColor)
				}
			}
		}

		for dy := -crossSize; dy <= crossSize; dy++ {
			for t := -crossThickness; t <= crossThickness; t++ {
				px, py := x+t, y+dy
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, crossColor)
				}
			}
		}
	}
}

// downloadImage télécharge une image depuis une URL
func downloadImage(url string) image.Image {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "GroupieTracker/1.0")
	req.Header.Set("Referer", "https://www.openstreetmap.org/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil
	}

	return img
}

// latLonToTile transforme GPS en numéro de tuile
func latLonToTile(lat, lon float64, zoom int) (int, int) {
	latRad := lat * math.Pi / 180.0
	n := math.Pow(2.0, float64(zoom))

	x := int((lon + 180.0) / 360.0 * n)
	y := int((1.0 - math.Log(math.Tan(latRad)+1/math.Cos(latRad))/math.Pi) / 2.0 * n)

	maxTile := int(n) - 1
	if x < 0 {
		x = 0
	}
	if x > maxTile {
		x = maxTile
	}
	if y < 0 {
		y = 0
	}
	if y > maxTile {
		y = maxTile
	}

	return x, y
}

// latLonToPixel conversion lat/lon → pixel dans l'image
func latLonToPixel(lat, lon float64, zoom int, centerTileX, centerTileY int) (int, int) {
	tileX, tileY := latLonToTile(lat, lon, zoom)

	dx := tileX - centerTileX
	dy := tileY - centerTileY

	px := (dx+1)*256 + 128
	py := (dy+1)*256 + 128

	return px, py
}
