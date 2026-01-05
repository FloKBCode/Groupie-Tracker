package models

import "testing"

func TestArtistStruct(t *testing.T) {
    a := Artist{
        ID:           1,
        Name:         "Queen",
        Members:      []string{"Freddie Mercury"},
        CreationDate: 1970,
        FirstAlbum:   "14-12-1973",
        Image:        "queen.jpg",
    }

    if a.ID != 1 || a.Name != "Queen" {
        t.Errorf("Artist struct fields incorrect")
    }
}
