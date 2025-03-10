package main

import (
    "database/sql"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "time"
)

// Share represents metadata for an uploaded share
type Share struct {
    ID         string    `json:"id"`
    Type       string    `json:"type"`       // "file" or "text"
    Size       int       `json:"size"`       // Size in bytes
    Expiration time.Time `json:"expiration"` // Expiration timestamp
    Host       string    `json:"host"`       // Host for URL construction (optional)
}

// ListShares retrieves all shares from the SQLite 'data' table
func ListShares(db *sql.DB) ([]Share, error) {
    var shares []Share
    // Query non-expired shares; expire is a Unix timestamp
    rows, err := db.Query("SELECT id, type, filePath, expire FROM data WHERE expire > ?", time.Now().Unix())
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var share Share
        var expireInt int64
        var filePath string
        if err := rows.Scan(&share.ID, &share.Type, &filePath, &expireInt); err != nil {
            continue
        }
        share.Expiration = time.Unix(expireInt, 0)
        // Calculate size from filePath if it exists
        if filePath != "" {
            fullPath := filepath.Join("uploads", filePath)
            fileInfo, err := os.Stat(fullPath)
            if err == nil {
                share.Size = int(fileInfo.Size())
            }
        }
        shares = append(shares, share)
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
