package services

import (
    "testing"
    "time"
)

func TestParseDate(t *testing.T) {
    input := "*23-08-2019"
    expected := time.Date(2019, 8, 23, 0, 0, 0, 0, time.UTC)

    got, err := ParseDate(input)
    if err != nil {
        t.Errorf("ParseDate error: %v", err)
    }
    if !got.Equal(expected) {
        t.Errorf("ParseDate(%s) = %v; want %v", input, got, expected)
    }
}

func TestParseLocation(t *testing.T) {
    city, country := ParseLocation("los_angeles-usa")
    if city != "los angeles" || country != "usa" {
        t.Errorf("ParseLocation = (%s, %s); want (los angeles, usa)", city, country)
    }
}
