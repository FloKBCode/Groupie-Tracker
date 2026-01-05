package models

type Location struct {
    ID        int      `json:"id"`
    Locations []string `json:"locations"`
    DatesURL  string   `json:"dates"` // lien vers les dates pour cette location
}
