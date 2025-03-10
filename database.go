/*
This file is part of ShareBin.

ShareBin is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

ShareBin is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with ShareBin. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
    "database/sql"
    "fmt"
    _ "modernc.org/sqlite"
)

func InitDatabase() *sql.DB {
    db, err := sql.Open("sqlite", "file:./data/sharebin.db?cache=shared")
    if err != nil {
        fmt.Println("Failed to open database:", err)
        return nil
    }

    db.SetMaxOpenConns(1)

    // Create the existing 'data' table
    createDataTableSQL := `CREATE TABLE IF NOT EXISTS data (
        id TEXT NOT NULL,
        type TEXT NOT NULL,
        fileName TEXT NOT NULL,
        filePath TEXT NOT NULL,
        burn TEXT NOT NULL,
        expire TEXT NOT NULL,
        passwordHash TEXT NOT NULL,
        passwordSalt TEXT NOT NULL,
        encryptSalt TEXT NOT NULL
    );`
    _, err = db.Exec(createDataTableSQL)
    if err != nil {
        fmt.Println("Failed to create data table:", err)
        return nil
    }

    return db
}
