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
    Size       int       `json:"size"`       // Size in bytes (derived or approximated)
    Expiration time.Time `json:"expiration"` // Expiration timestamp
}

// ListShares retrieves all shares from the SQLite 'data' table
func ListShares(db *sql.DB) ([]Share, error) {
    var shares []Share
    rows, err := db.Query("SELECT id, type, expire FROM data WHERE expire > ?", time.Now().Unix())
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var share Share
        var expireStr string
        if err := rows.Scan(&share.ID, &share.Type, &expireStr); err != nil {
            continue
        }
        // Convert expire (TEXT) to time.Time
        expireInt, err := parseExpiration(expireStr)
        if err != nil {
            continue
        }
        share.Expiration = time.Unix(expireInt, 0)
        // Size isn’t in 'data' table; approximate or fetch from file if needed
        share.Size = 0 // Placeholder; update if size is calculable
        shares = append(shares, share)
    }
    return shares, nil
}

// parseExpiration converts the expire TEXT field to Unix timestamp
func parseExpiration(expireStr string) (int64, error) {
    // Assuming expire is stored as Unix timestamp string; adjust if format differs
    return time.Parse("2006-01-02 15:04:05", expireStr).Unix()
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
