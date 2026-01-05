package api

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
)

// FetchJSON récupère l'URL et la décode dans target
func FetchJSON(url string, target interface{}) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    return json.Unmarshal(body, target)
}
