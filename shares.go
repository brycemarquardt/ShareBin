package main

import (
    "database/sql"
    "sort"
    "strconv"
    "time"
)

// Share represents metadata for an uploaded share
type Share struct {
    ID         string    `json:"id"`
    Type       string    `json:"type"`       // "file" or "text"
    Size       int       `json:"size"`       // Size in bytes (approximated if not stored)
    Expiration time.Time `json:"expiration"` // Expiration timestamp
}

// ListShares retrieves all shares from the SQLite 'data' table
func ListShares(db *sql.DB) ([]Share, error) {
    var shares []Share
    // Query non-expired shares; expire is TEXT, so we compare as string or convert
    rows, err := db.Query("SELECT id, type, expire FROM data WHERE expire > ?", strconv.FormatInt(time.Now().Unix(), 10))
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
        expireInt, err := strconv.ParseInt(expireStr, 10, 64)
        if err != nil {
            // If expire isn’t a Unix timestamp, adjust parsing logic below
            continue
        }
        share.Expiration = time.Unix(expireInt, 0)
        // Size isn’t in 'data'; set to 0 or calculate from filePath if needed
        share.Size = 0 // Placeholder
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
