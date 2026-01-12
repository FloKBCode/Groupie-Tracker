package models

type ConcertDate struct {
    ID    int      `json:"id"`
    Dates []string `json:"dates"`
}
