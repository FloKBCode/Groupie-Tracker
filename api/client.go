package api

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// Client HTTP avec timeout pour éviter les blocages
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
}

// FetchJSON récupère l'URL et la décode dans target
func FetchJSON(url string, target interface{}) error {
    resp, err := httpClient.Get(url)
    if err != nil {
        return fmt.Errorf("erreur lors de la requête GET %s: %w", url, err)
    }
    defer resp.Body.Close()

    // Vérifier le code de statut HTTP
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("statut HTTP %d pour %s", resp.StatusCode, url)
    }

    // io.ReadAll au lieu de ioutil.ReadAll (deprecated)
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("erreur lecture body: %w", err)
    }

    if err := json.Unmarshal(body, target); err != nil {
        return fmt.Errorf("erreur décodage JSON: %w", err)
    }

    return nil
}