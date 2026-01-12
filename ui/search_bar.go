package ui

import (
	"fmt"
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SearchBar représente la barre de recherche avec autocomplétion
type SearchBar struct {
	Container      fyne.CanvasObject
	searchEngine   *services.SearchEngine
	entry          *widget.Entry
	suggestionList *widget.List
	suggestions    []services.SearchResult
	onSelect       func(int) // Callback quand on sélectionne un artiste
}

// NewSearchBar crée une nouvelle barre de recherche
func NewSearchBar(searchEngine *services.SearchEngine, onSelect func(int)) *SearchBar {
	sb := &SearchBar{
		searchEngine: searchEngine,
		suggestions:  []services.SearchResult{},
		onSelect:     onSelect,
	}

	// Entry de recherche
	sb.entry = widget.NewEntry()
	sb.entry.SetPlaceHolder("🔍 Rechercher un artiste, membre, lieu ou date...")

	// Liste de suggestions (cachée par défaut)
	sb.suggestionList = widget.NewList(
		func() int {
			return len(sb.suggestions)
		},
		func() fyne.CanvasObject {
			typeLabel := widget.NewLabel("")
			typeLabel.TextStyle = fyne.TextStyle{Italic: true}
			
			matchLabel := widget.NewLabel("")
			matchLabel.TextStyle = fyne.TextStyle{Bold: true}
			
			artistLabel := widget.NewLabel("")
			
			return container.NewVBox(
				container.NewHBox(matchLabel, typeLabel),
				artistLabel,
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			suggestion := sb.suggestions[id]
			
			vbox := obj.(*fyne.Container)
			topRow := vbox.Objects[0].(*fyne.Container)
			matchLabel := topRow.Objects[0].(*widget.Label)
			typeLabel := topRow.Objects[1].(*widget.Label)
			artistLabel := vbox.Objects[1].(*widget.Label)
			
			// Type avec emoji
			var typeStr string
			switch suggestion.Type {
			case services.SearchTypeArtist:
				typeStr = " 🎵 Artiste"
			case services.SearchTypeMember:
				typeStr = " 👤 Membre"
			case services.SearchTypeLocation:
				typeStr = " 📍 Lieu"
			case services.SearchTypeDate:
				typeStr = " 📅 Date"
			}
			
			matchLabel.SetText(suggestion.MatchedText)
			typeLabel.SetText(typeStr)
			artistLabel.SetText("→ " + suggestion.ArtistName)
		},
	)

	// Gestion de la sélection d'une suggestion
	sb.suggestionList.OnSelected = func(id widget.ListItemID) {
		if id < len(sb.suggestions) {
			selectedArtistID := sb.suggestions[id].ArtistID
			
			// Nettoyer l'interface
			sb.entry.SetText("")
			sb.suggestions = []services.SearchResult{}
			sb.suggestionList.Refresh()
			
			// Callback vers la vue détails
			if sb.onSelect != nil {
				sb.onSelect(selectedArtistID)
			}
		}
	}

	// Événement de changement de texte
	sb.entry.OnChanged = func(query string) {
		sb.updateSuggestions(query)
	}

	// Layout avec suggestions qui apparaissent en dessous
	searchContainer := container.NewBorder(
		sb.entry,
		nil,
		nil,
		nil,
		sb.suggestionList,
	)

	sb.Container = searchContainer
	return sb
}

// updateSuggestions met à jour les suggestions basées sur la query
func (sb *SearchBar) updateSuggestions(query string) {
	if query == "" {
		sb.suggestions = []services.SearchResult{}
		sb.suggestionList.Refresh()
		return
	}

	// Limiter à 10 suggestions
	sb.suggestions = sb.searchEngine.GetSuggestions(query, 10)
	sb.suggestionList.Refresh()
	
	fmt.Printf("🔍 Recherche: '%s' -> %d résultats\n", query, len(sb.suggestions))
}

// Clear vide la barre de recherche
func (sb *SearchBar) Clear() {
	sb.entry.SetText("")
	sb.suggestions = []services.SearchResult{}
	sb.suggestionList.Refresh()
}