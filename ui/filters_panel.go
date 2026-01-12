package ui

import (
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FiltersPanel struct {
	Container fyne.CanvasObject
	criteria  *services.FilterCriteria
	onApply   func(*services.FilterCriteria)
}

func NewFiltersPanel(onApply func(*services.FilterCriteria)) *FiltersPanel {
	fp := &FiltersPanel{
		criteria: services.NewFilterCriteria(),
		onApply:  onApply,
	}

	// Créer les widgets
	fp.buildUI()
	return fp
}

func (fp *FiltersPanel) buildUI() {
	// Slider pour années de création
	creationSlider := widget.NewSlider(1950, 2025)
	creationSlider.OnChanged = func(val float64) {
		fp.criteria.CreationDateMin = int(val)
	}

	// Slider pour membres
	membersSlider := widget.NewSlider(1, 10)
	membersSlider.OnChanged = func(val float64) {
		fp.criteria.MembersMin = int(val)
	}

	// Bouton appliquer
	applyBtn := widget.NewButton("✅ Appliquer", func() {
		if fp.onApply != nil {
			fp.onApply(fp.criteria)
		}
	})

	// Layout
	fp.Container = container.NewVBox(
		widget.NewLabel("🔧 Filtres"),
		widget.NewSeparator(),
		creationSlider,
		membersSlider,
		applyBtn,
	)
}
