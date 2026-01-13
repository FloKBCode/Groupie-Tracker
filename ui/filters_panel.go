package ui

import (
	"fmt"
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// FiltersPanel - Panneau de filtres amélioré avec fenêtre séparée
type FiltersPanel struct {
	window   fyne.Window
	criteria *services.FilterCriteria
	onApply  func(*services.FilterCriteria)

	// Widgets pour Année de Création
	creationCheck     *widget.Check
	creationMinSlider *widget.Slider
	creationMaxSlider *widget.Slider
	creationMinLabel  *widget.Label
	creationMaxLabel  *widget.Label

	// Widgets pour Premier Album
	albumCheck     *widget.Check
	albumMinSlider *widget.Slider
	albumMaxSlider *widget.Slider
	albumMinLabel  *widget.Label
	albumMaxLabel  *widget.Label

	// Widgets pour Membres
	membersCheck     *widget.Check
	membersMinSlider *widget.Slider
	membersMaxSlider *widget.Slider
	membersMinLabel  *widget.Label
	membersMaxLabel  *widget.Label

	// Widgets pour Locations
	locationCheck  *widget.Check
	locationSelect *widget.Select

	// Boutons
	applyButton *widget.Button
	resetButton *widget.Button
	closeButton *widget.Button
}

// NewFiltersPanel crée un nouveau panneau de filtres amélioré
func NewFiltersPanel(onApply func(*services.FilterCriteria)) *FiltersPanel {
	fp := &FiltersPanel{
		criteria: services.NewFilterCriteria(),
		onApply:  onApply,
	}

	fp.buildWindow()
	return fp
}

// buildWindow crée la fenêtre séparée (déplaçable)
func (fp *FiltersPanel) buildWindow() {
	// Créer une nouvelle fenêtre qui sera déplaçable
	fp.window = fyne.CurrentApp().NewWindow("🔧 Filtres Avancés")

	// Construire le contenu
	content := fp.buildContent()

	// Wrapper dans un scroll pour gérer le contenu long
	scrollContent := container.NewVScroll(content)

	fp.window.SetContent(scrollContent)
	fp.window.Resize(fyne.NewSize(500, 650))
	fp.window.CenterOnScreen()
}

// buildContent construit le contenu complet du panneau
func (fp *FiltersPanel) buildContent() fyne.CanvasObject {
	// Titre
	title := widget.NewLabelWithStyle(
		"🔧 Filtres Avancés",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Sections
	creationSection := fp.buildCreationDateSection()
	albumSection := fp.buildFirstAlbumSection()
	membersSection := fp.buildMembersSection()
	locationSection := fp.buildLocationSection()
	buttonsSection := fp.buildButtons()

	// Assemblage
	return container.NewVBox(
		title,
		widget.NewSeparator(),
		creationSection,
		widget.NewSeparator(),
		albumSection,
		widget.NewSeparator(),
		membersSection,
		widget.NewSeparator(),
		locationSection,
		widget.NewSeparator(),
		buttonsSection,
	)
}

// buildCreationDateSection crée la section filtre par année de création
func (fp *FiltersPanel) buildCreationDateSection() fyne.CanvasObject {
	// Checkbox pour activer/désactiver le filtre
	fp.creationCheck = widget.NewCheck("Activer ce filtre", func(checked bool) {
		fp.criteria.EnableCreationDateFilter = checked
		if checked {
			fp.creationMinSlider.Enable()
			fp.creationMaxSlider.Enable()
		} else {
			fp.creationMinSlider.Disable()
			fp.creationMaxSlider.Disable()
		}
	})

	// Labels pour afficher les valeurs
	fp.creationMinLabel = widget.NewLabel("Min: 1950")
	fp.creationMaxLabel = widget.NewLabel("Max: 2025")

	// Slider minimum
	fp.creationMinSlider = widget.NewSlider(1950, 2025)
	fp.creationMinSlider.SetValue(1950)
	fp.creationMinSlider.Step = 1
	fp.creationMinSlider.OnChanged = func(val float64) {
		fp.criteria.CreationDateMin = int(val)
		fp.creationMinLabel.SetText(fmt.Sprintf("Min: %d", int(val)))
		
		// S'assurer que min <= max
		if val > fp.creationMaxSlider.Value {
			fp.creationMaxSlider.SetValue(val)
		}
	}
	fp.creationMinSlider.Disable()

	// Slider maximum
	fp.creationMaxSlider = widget.NewSlider(1950, 2025)
	fp.creationMaxSlider.SetValue(2025)
	fp.creationMaxSlider.Step = 1
	fp.creationMaxSlider.OnChanged = func(val float64) {
		fp.criteria.CreationDateMax = int(val)
		fp.creationMaxLabel.SetText(fmt.Sprintf("Max: %d", int(val)))
		
		// S'assurer que max >= min
		if val < fp.creationMinSlider.Value {
			fp.creationMinSlider.SetValue(val)
		}
	}
	fp.creationMaxSlider.Disable()

	// Container pour les sliders avec labels
	slidersContent := container.NewVBox(
		fp.creationCheck,
		widget.NewLabel(""),
		container.NewHBox(fp.creationMinLabel, widget.NewLabel("")),
		fp.creationMinSlider,
		widget.NewLabel(""),
		container.NewHBox(fp.creationMaxLabel, widget.NewLabel("")),
		fp.creationMaxSlider,
	)

	return widget.NewCard(
		"📅 Année de Création",
		"Filtrer par année de formation du groupe",
		slidersContent,
	)
}

// buildFirstAlbumSection crée la section filtre par premier album
func (fp *FiltersPanel) buildFirstAlbumSection() fyne.CanvasObject {
	fp.albumCheck = widget.NewCheck("Activer ce filtre", func(checked bool) {
		fp.criteria.EnableFirstAlbumFilter = checked
		if checked {
			fp.albumMinSlider.Enable()
			fp.albumMaxSlider.Enable()
		} else {
			fp.albumMinSlider.Disable()
			fp.albumMaxSlider.Disable()
		}
	})

	fp.albumMinLabel = widget.NewLabel("Min: 1950")
	fp.albumMaxLabel = widget.NewLabel("Max: 2025")

	fp.albumMinSlider = widget.NewSlider(1950, 2025)
	fp.albumMinSlider.SetValue(1950)
	fp.albumMinSlider.Step = 1
	fp.albumMinSlider.OnChanged = func(val float64) {
		fp.criteria.FirstAlbumYearMin = int(val)
		fp.albumMinLabel.SetText(fmt.Sprintf("Min: %d", int(val)))
		
		if val > fp.albumMaxSlider.Value {
			fp.albumMaxSlider.SetValue(val)
		}
	}
	fp.albumMinSlider.Disable()

	fp.albumMaxSlider = widget.NewSlider(1950, 2025)
	fp.albumMaxSlider.SetValue(2025)
	fp.albumMaxSlider.Step = 1
	fp.albumMaxSlider.OnChanged = func(val float64) {
		fp.criteria.FirstAlbumYearMax = int(val)
		fp.albumMaxLabel.SetText(fmt.Sprintf("Max: %d", int(val)))
		
		if val < fp.albumMinSlider.Value {
			fp.albumMinSlider.SetValue(val)
		}
	}
	fp.albumMaxSlider.Disable()

	slidersContent := container.NewVBox(
		fp.albumCheck,
		widget.NewLabel(""),
		container.NewHBox(fp.albumMinLabel, widget.NewLabel("")),
		fp.albumMinSlider,
		widget.NewLabel(""),
		container.NewHBox(fp.albumMaxLabel, widget.NewLabel("")),
		fp.albumMaxSlider,
	)

	return widget.NewCard(
		"💿 Premier Album",
		"Filtrer par année du premier album",
		slidersContent,
	)
}

// buildMembersSection crée la section filtre par nombre de membres
func (fp *FiltersPanel) buildMembersSection() fyne.CanvasObject {
	fp.membersCheck = widget.NewCheck("Activer ce filtre", func(checked bool) {
		fp.criteria.EnableMembersFilter = checked
		if checked {
			fp.membersMinSlider.Enable()
			fp.membersMaxSlider.Enable()
		} else {
			fp.membersMinSlider.Disable()
			fp.membersMaxSlider.Disable()
		}
	})

	fp.membersMinLabel = widget.NewLabel("Min: 1 membre")
	fp.membersMaxLabel = widget.NewLabel("Max: 10 membres")

	fp.membersMinSlider = widget.NewSlider(1, 10)
	fp.membersMinSlider.SetValue(1)
	fp.membersMinSlider.Step = 1
	fp.membersMinSlider.OnChanged = func(val float64) {
		fp.criteria.MembersMin = int(val)
		fp.membersMinLabel.SetText(fmt.Sprintf("Min: %d membre(s)", int(val)))
		
		if val > fp.membersMaxSlider.Value {
			fp.membersMaxSlider.SetValue(val)
		}
	}
	fp.membersMinSlider.Disable()

	fp.membersMaxSlider = widget.NewSlider(1, 10)
	fp.membersMaxSlider.SetValue(10)
	fp.membersMaxSlider.Step = 1
	fp.membersMaxSlider.OnChanged = func(val float64) {
		fp.criteria.MembersMax = int(val)
		fp.membersMaxLabel.SetText(fmt.Sprintf("Max: %d membre(s)", int(val)))
		
		if val < fp.membersMinSlider.Value {
			fp.membersMinSlider.SetValue(val)
		}
	}
	fp.membersMaxSlider.Disable()

	slidersContent := container.NewVBox(
		fp.membersCheck,
		widget.NewLabel(""),
		container.NewHBox(fp.membersMinLabel, widget.NewLabel("")),
		fp.membersMinSlider,
		widget.NewLabel(""),
		container.NewHBox(fp.membersMaxLabel, widget.NewLabel("")),
		fp.membersMaxSlider,
	)

	return widget.NewCard(
		"👥 Nombre de Membres",
		"Filtrer par taille du groupe",
		slidersContent,
	)
}

// buildLocationSection crée la section filtre par pays
func (fp *FiltersPanel) buildLocationSection() fyne.CanvasObject {
	fp.locationCheck = widget.NewCheck("Activer ce filtre", func(checked bool) {
		fp.criteria.EnableLocationsFilter = checked
		if checked {
			fp.locationSelect.Enable()
		} else {
			fp.locationSelect.Disable()
		}
	})

	// Liste de pays (à remplir dynamiquement avec GetAvailableLocations())
	countries := []string{
		"USA",
		"FRANCE",
		"UK",
		"GERMANY",
		"JAPAN",
		"CANADA",
		"SPAIN",
		"ITALY",
	}

	fp.locationSelect = widget.NewSelect(
		countries,
		func(selected string) {
			if selected != "" {
				fp.criteria.Locations = []string{selected}
			} else {
				fp.criteria.Locations = []string{}
			}
		},
	)
	fp.locationSelect.PlaceHolder = "Sélectionner un pays..."
	fp.locationSelect.Disable()

	content := container.NewVBox(
		fp.locationCheck,
		widget.NewLabel(""),
		fp.locationSelect,
	)

	return widget.NewCard(
		"🌍 Pays de Concert",
		"Filtrer par lieu de concert",
		content,
	)
}

// buildButtons crée les boutons d'action
func (fp *FiltersPanel) buildButtons() fyne.CanvasObject {
	fp.applyButton = widget.NewButton("✅ Appliquer les Filtres", func() {
		if fp.onApply != nil {
			fp.onApply(fp.criteria)
		}
		fp.window.Hide()
	})

	fp.resetButton = widget.NewButton("🔄 Réinitialiser", func() {
		fp.resetFilters()
	})

	fp.closeButton = widget.NewButton("❌ Fermer", func() {
		fp.window.Hide()
	})

	// Disposition horizontale des boutons
	return container.NewGridWithColumns(
		3,
		fp.applyButton,
		fp.resetButton,
		fp.closeButton,
	)
}

// resetFilters réinitialise tous les filtres aux valeurs par défaut
func (fp *FiltersPanel) resetFilters() {
	// Désactiver toutes les checkboxes
	fp.creationCheck.SetChecked(false)
	fp.albumCheck.SetChecked(false)
	fp.membersCheck.SetChecked(false)
	fp.locationCheck.SetChecked(false)

	// Réinitialiser les sliders
	fp.creationMinSlider.SetValue(1950)
	fp.creationMaxSlider.SetValue(2025)
	fp.albumMinSlider.SetValue(1950)
	fp.albumMaxSlider.SetValue(2025)
	fp.membersMinSlider.SetValue(1)
	fp.membersMaxSlider.SetValue(10)

	// Réinitialiser la sélection
	fp.locationSelect.SetSelected("")

	// Réinitialiser les critères
	fp.criteria = services.NewFilterCriteria()

	fmt.Println("🔄 Filtres réinitialisés")
}

// Show affiche la fenêtre de filtres
func (fp *FiltersPanel) Show() {
	fp.window.Show()
}

// Hide cache la fenêtre de filtres
func (fp *FiltersPanel) Hide() {
	fp.window.Hide()
}

// LoadAvailableLocations charge les pays disponibles depuis le FilterEngine
func (fp *FiltersPanel) LoadAvailableLocations(filterEngine *services.FilterEngine) {
	locations := filterEngine.GetAvailableLocations()
	if len(locations) > 0 {
		fp.locationSelect.Options = locations
		fp.locationSelect.Refresh()
		fmt.Printf("✅ %d pays chargés dans le filtre\n", len(locations))
	}
}