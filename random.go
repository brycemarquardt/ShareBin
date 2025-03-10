/*
This file is part of ShareBin.

ShareBin is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

ShareBin is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with ShareBin. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
    "database/sql"
    "math/rand"
    "os"
    "strconv"
    "time"
)

// We want to avoid ambiguous characters like i, I, l, 1, etc.
const charset = "abcdefghkmnpqrstwxyzABCDEFGHJKLMNPQRSTWXYZ23456789"

var seed rand.Source
var random *rand.Rand

func init() {
    seed = rand.NewSource(time.Now().UnixNano())
    random = rand.New(seed)
}

func GenRandFileName(basePath string, extension string) string {
    for {
        fileName := strconv.FormatInt(time.Now().UnixMilli(), 10) + genRandString(5) + extension
        filePath := basePath + fileName

        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            return filePath
        }
    }
}

func GenRandPath(length int, db *sql.DB) string {
    for {
        randPath := genRandString(6)
        var id string

        db.QueryRow("SELECT id FROM data WHERE id = ?", randPath).Scan(&id)

        if id == "" {
            return randPath
        }
    }
}

func genRandString(length int) string {
    result := make([]byte, length)
    for i := range result {
        result[i] = charset[random.Intn(len(charset))]
    }
    return string(result)
}
