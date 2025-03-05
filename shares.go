package main

import (
    "encoding/json"
    "io/fs"
    "os"
    "path/filepath"
    "sort"
    "time"
)

// Share represents metadata for an uploaded share
type Share struct {
    ID         string    `json:"id"`
    Type       string    `json:"type"`       // "file" or "text"
    Size       int       `json:"size"`       // Size in bytes
    Expiration time.Time `json:"expiration"` // Expiration timestamp
}

// ListShares retrieves all shares from the uploads directory
func ListShares() ([]Share, error) {
    var shares []Share
    entries, err := os.ReadDir("uploads")
    if err != nil {
        return nil, err
    }
    for _, entry := range entries {
        if entry.IsDir() {
            metaPath := filepath.Join("uploads", entry.Name(), "meta.json")
            metaFile, err := os.Open(metaPath)
            if err != nil {
                continue // Skip if meta.json is missing
            }
            var share Share
            if err := json.NewDecoder(metaFile).Decode(&share); err != nil {
                metaFile.Close()
                continue
            }
            metaFile.Close()
            share.ID = entry.Name() // Set ID from directory name
            shares = append(shares, share)
        }
    }
    return shares, nil
}

// FilterSharesByType filters shares by type (e.g., "file" or "text")
func FilterSharesByType(shares []Share, shareType string) []Share {
    var filtered []Share
    for _, share := range shares {
        if share.Type == shareType {
            filtered = append(filtered, share)
        }
    }
    return filtered
}

// SortShares sorts shares by a given field and order
func SortShares(shares []Share, sortBy, order string) {
    sort.Slice(shares, func(i, j int) bool {
        if order == "desc" {
            i, j = j, i // Reverse order for descending
        }
        switch sortBy {
        case "size":
            return shares[i].Size < shares[j].Size
        case "expiration":
            return shares[i].Expiration.Before(shares[j].Expiration)
        default: // "id"
            return shares[i].ID < shares[j].ID
        }
    })
}
