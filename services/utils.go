package services

import (
    "strings"
    "time"
)

// ParseDate convertit une date au format "*dd-mm-yyyy" en time.Time
func ParseDate(s string) (time.Time, error) {
    s = strings.TrimPrefix(s, "*") // supprime l'éventuel *
    return time.Parse("02-01-2006", s)
}

// ParseLocation sépare une location de type "city-country" en ville et pays
func ParseLocation(s string) (city, country string) {
	if s == "" {
		return "", ""
	}

	parts := strings.Split(s, "-")
	if len(parts) < 2 {
		return "", ""
	}

	city = strings.ReplaceAll(parts[0], "_", " ")
	country = strings.ReplaceAll(parts[1], "_", " ")

	city = strings.TrimSpace(city)
	country = strings.TrimSpace(country)

	return city, country
}


// FormatDate convertit une date du format "DD-MM-YYYY" ou "*DD-MM-YYYY" au format "JJ/MM/AAAA"
func FormatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	
	// Retirer l'éventuel astérisque
	dateStr = strings.TrimPrefix(dateStr, "*")
	
	// Parser la date
	t, err := ParseDate(dateStr)
	if err != nil {
		return dateStr // Retourner tel quel si format invalide
	}
	
	// Reformater en JJ/MM/AAAA
	return t.Format("02/01/2006")
}

// FormatDateList formate une liste de dates
func FormatDateList(dates []string) []string {
	formatted := make([]string, len(dates))
	for i, date := range dates {
		formatted[i] = FormatDate(date)
	}
	return formatted
}
