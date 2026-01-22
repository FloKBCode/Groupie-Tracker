package ui

import (
	"fmt"
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SearchBar représente une barre de recherche améliorée avec tous les moteurs
type SearchBar struct {
	Container      fyne.CanvasObject
	
	// Moteurs de recherche
	searchEngine   *services.SearchEngine
	fuzzyEngine    *services.FuzzySearchEngine
	initialsEngine *services.InitialsSearchEngine
	searchHistory  *services.SearchHistory
	
	// Widgets
	entry          *widget.Entry
	suggestionList *widget.List
	suggestions    []services.SearchResult
	onSelect       func(int) // Callback quand on sélectionne un artiste
	
	// État
	suggestionsVisible bool
}

// NewSearchBar crée une nouvelle barre de recherche complète
func NewSearchBar(searchEngine *services.SearchEngine, onSelect func(int)) *SearchBar {
	sb := &SearchBar{
		searchEngine:       searchEngine,
		fuzzyEngine:        services.NewFuzzySearchEngine(searchEngine),
		initialsEngine:     services.NewInitialsSearchEngine(searchEngine),
		searchHistory:      services.NewSearchHistory(50),
		suggestions:        []services.SearchResult{},
		onSelect:           onSelect,
		suggestionsVisible: false,
	}

	// Entry de recherche
	sb.entry = widget.NewEntry()
	sb.entry.SetPlaceHolder("🔍 Rechercher (essayez 'fm' pour Freddie Mercury ou 'qeen' pour Queen)...")
	
	// Liste de suggestions
	sb.suggestionList = widget.NewList(
		func() int {
			return len(sb.suggestions)
		},
		func() fyne.CanvasObject {
			return sb.createSuggestionTemplate()
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			sb.updateSuggestionItem(id, obj)
		},
	)

	sb.suggestionList.Resize(fyne.NewSize(0, 300))

	// Gestion de la sélection
	sb.suggestionList.OnSelected = func(id widget.ListItemID) {
		if id < len(sb.suggestions) {
			selected := sb.suggestions[id]
			
			// Sauvegarder dans l'historique
			sb.searchHistory.Add(sb.entry.Text, selected.ArtistID)
			
			// Nettoyer l'interface
			sb.entry.SetText("")
			sb.hideSuggestions()
			
			// Callback
			if sb.onSelect != nil {
				sb.onSelect(selected.ArtistID)
			}
		}
	}

	// Événement de changement de texte
	sb.entry.OnChanged = func(query string) {
		sb.updateSuggestionsAdvanced(query)
	}

	// Enter pour sélectionner
	sb.entry.OnSubmitted = func(query string) {
		if len(sb.suggestions) > 0 {
			sb.suggestionList.Select(0)
		}
	}

	// Layout
	searchContainer := container.NewBorder(
		sb.entry,
		nil,
		nil,
		nil,
		sb.suggestionList,
	)

	sb.suggestionList.Hide()

	sb.Container = searchContainer
	return sb
}

// createSuggestionTemplate crée le template pour une suggestion
func (sb *SearchBar) createSuggestionTemplate() fyne.CanvasObject {
	typeIcon := widget.NewLabel("")
	matchedLabel := widget.NewRichTextFromMarkdown("")
	artistLabel := widget.NewLabel("")
	artistLabel.TextStyle = fyne.TextStyle{Italic: true}
	
	scoreLabel := widget.NewLabel("")
	scoreLabel.TextStyle = fyne.TextStyle{Monospace: true}
	scoreLabel.Hide() // Caché par défaut
	
	return container.NewVBox(
		container.NewHBox(typeIcon, matchedLabel),
		artistLabel,
		scoreLabel,
	)
}

// updateSuggestionItem met à jour une suggestion
func (sb *SearchBar) updateSuggestionItem(id widget.ListItemID, obj fyne.CanvasObject) {
	suggestion := sb.suggestions[id]
	
	vbox := obj.(*fyne.Container)
	topRow := vbox.Objects[0].(*fyne.Container)
	typeIcon := topRow.Objects[0].(*widget.Label)
	matchedLabel := topRow.Objects[1].(*widget.RichText)
	artistLabel := vbox.Objects[1].(*widget.Label)
	scoreLabel := vbox.Objects[2].(*widget.Label)
	
	// Icône selon le type
	var icon string
	switch suggestion.Type {
	case services.SearchTypeArtist:
		icon = "🎵"
	case services.SearchTypeMember:
		icon = "👤"
	case services.SearchTypeLocation:
		icon = "📍"
	case services.SearchTypeDate:
		icon = "📅"
	}
	typeIcon.SetText(icon)
	
	// Highlighting avec RichText (pas d'espaces)
	before, match, after := sb.searchEngine.HighlightMatch(suggestion)
	markdownText := before + "**" + match + "**" + after
	matchedLabel.ParseMarkdown(markdownText)
	
	// Nom de l'artiste
	if suggestion.Type != services.SearchTypeArtist {
		artistLabel.SetText("→ " + suggestion.ArtistName)
		artistLabel.Show()
	} else {
		artistLabel.Hide()
	}
	
	// Score (pour debug, peut être affiché)
	scoreLabel.SetText(fmt.Sprintf("[%d pts]", suggestion.Score))
}

// updateSuggestionsAdvanced - Version avancée avec tous les moteurs
func (sb *SearchBar) updateSuggestionsAdvanced(query string) {
	if query == "" {
		// Afficher l'historique récent quand la recherche est vide
		sb.showHistorySuggestions()
		return
	}

	// Combiner tous les types de recherche
	allResults := []services.SearchResult{}
	seen := make(map[string]bool) // Pour éviter les doublons
	
	// 1. Recherche normale (score le plus élevé)
	normalResults := sb.searchEngine.GetSuggestions(query, 8)
	for _, r := range normalResults {
		key := fmt.Sprintf("%d-%s-%s", r.ArtistID, r.MatchedText, r.Type)
		if !seen[key] {
			allResults = append(allResults, r)
			seen[key] = true
		}
	}
	
	// 2. Recherche par initiales (si query courte)
	if len(query) >= 2 && len(query) <= 5 {
		initialsResults := sb.initialsEngine.SearchByInitials(query)
		for _, r := range initialsResults {
			key := fmt.Sprintf("%d-%s-%s", r.ArtistID, r.MatchedText, r.Type)
			if !seen[key] {
				allResults = append(allResults, r)
				seen[key] = true
			}
		}
	}
	
	// 3. Recherche floue (si peu de résultats normaux)
	if len(normalResults) < 3 {
		fuzzyResults := sb.fuzzyEngine.FuzzySearch(query, 2)
		for _, r := range fuzzyResults {
			key := fmt.Sprintf("%d-%s-%s", r.ArtistID, r.MatchedText, r.Type)
			if !seen[key] {
				allResults = append(allResults, r)
				seen[key] = true
			}
		}
	}
	
	// Limiter à 10 suggestions
	if len(allResults) > 10 {
		allResults = allResults[:10]
	}
	
	sb.suggestions = allResults
	
	if len(sb.suggestions) > 0 {
		sb.showSuggestions()
	} else {
		sb.hideSuggestions()
	}
	
	sb.suggestionList.Refresh()
	
	fmt.Printf("🔍 Recherche avancée: '%s' -> %d résultats (normal: %d, initiales: vérifiées, fuzzy: vérifiée)\n", 
		query, len(sb.suggestions), len(normalResults))
}

// showHistorySuggestions affiche les suggestions de l'historique
func (sb *SearchBar) showHistorySuggestions() {
	recent := sb.searchHistory.GetRecent(5)
	
	if len(recent) == 0 {
		sb.hideSuggestions()
		return
	}
	
	// Convertir l'historique en suggestions
	sb.suggestions = []services.SearchResult{}
	
	for _, entry := range recent {
		sb.suggestions = append(sb.suggestions, services.SearchResult{
			ArtistID:    entry.ResultID,
			ArtistName:  "",
			MatchedText: entry.Query + " (récent)",
			Type:        services.SearchTypeArtist,
			Score:       100,
		})
	}
	
	sb.showSuggestions()
	sb.suggestionList.Refresh()
}

// showSuggestions affiche la liste
func (sb *SearchBar) showSuggestions() {
	if !sb.suggestionsVisible {
		sb.suggestionList.Show()
		sb.suggestionsVisible = true
	}
}

// hideSuggestions masque la liste
func (sb *SearchBar) hideSuggestions() {
	if sb.suggestionsVisible {
		sb.suggestionList.Hide()
		sb.suggestionsVisible = false
	}
	sb.suggestions = []services.SearchResult{}
	sb.suggestionList.Refresh()
}

// Clear vide la barre de recherche
func (sb *SearchBar) Clear() {
	sb.entry.SetText("")
	sb.hideSuggestions()
}

// Focus met le focus sur la barre
func (sb *SearchBar) Focus() {
	sb.entry.FocusGained()
}

// SetPlaceholder change le placeholder
func (sb *SearchBar) SetPlaceholder(text string) {
	sb.entry.SetPlaceHolder(text)
}

// GetSearchHistory retourne l'historique (pour affichage externe)
func (sb *SearchBar) GetSearchHistory() *services.SearchHistory {
	return sb.searchHistory
}