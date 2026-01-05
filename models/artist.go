package models

type Artist struct {
    ID           int      `json:"id"`
    Image        string   `json:"image"`
    Name         string   `json:"name"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`  // année de création
    FirstAlbum   string   `json:"firstAlbum"`    // format date "14-12-1973"
    LocationsURL string   `json:"locations"`     // lien vers l’API locations
    DatesURL     string   `json:"concertDates"`  // lien vers l’API dates
    RelationsURL string   `json:"relations"`     // lien vers l’API relations
}
