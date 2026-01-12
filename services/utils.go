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

