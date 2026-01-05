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
    parts := strings.Split(s, "-")
    city = strings.ReplaceAll(parts[0], "_", " ")
    if len(parts) > 1 {
        country = strings.ReplaceAll(parts[1], "_", " ")
    } else {
        country = ""
    }
    return
}
