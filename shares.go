package main

import (
    "database/sql"
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

// ListShares retrieves all shares from the SQLite database
func ListShares(db *sql.DB) ([]Share, error) {
    var shares []Share
    rows, err := db.Query("SELECT id, type, size, expiration FROM shares WHERE expiration > ?", time.Now().Unix())
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var share Share
        var expiration int64
        if err := rows.Scan(&share.ID, &share.Type, &share.Size, &expiration); err != nil {
            continue
        }
        share.Expiration = time.Unix(expiration, 0)
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
