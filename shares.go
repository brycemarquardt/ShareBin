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
    // Query non-expired shares; expire is number of minutes from upload time
    rows, err := db.Query("SELECT id, type, expire FROM data")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var share Share
        var expireMinutesStr string
        if err := rows.Scan(&share.ID, &share.Type, &expireMinutesStr); err != nil {
            continue
        }
        // Convert expire (TEXT, number of minutes) to integer
        expireMinutes, err := strconv.ParseInt(expireMinutesStr, 10, 64)
        if err != nil {
            continue
        }
        // Calculate expiration as current time plus minutes
        share.Expiration = time.Now().Add(time.Duration(expireMinutes) * time.Minute)
        // Size isn’t in 'data'; set to 0 or calculate from filePath if needed
        share.Size = 0 // Placeholder
        // Only include non-expired shares (expiration > now)
        if share.Expiration.After(time.Now()) {
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
