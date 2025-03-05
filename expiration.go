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
    "os"
    "strconv"
    "time"
)

func CheckExpiration(db *sql.DB) {
    for {
        tx, err := db.Begin()
        if err != nil {
            fmt.Println("ShareBin: Failed to begin transaction:", err)
        }

        // SQLite doesn't support SELECT FOR UPDATE; the transaction locks the entire DB file
        rows, err := tx.Query("SELECT id, filePath FROM data WHERE expire <= ?", strconv.FormatInt(time.Now().Unix(), 10))
        if err != nil {
            tx.Rollback()
            fmt.Println("ShareBin: Error querying expired shares:", err)
            return
        }

        var toDelete = []struct {
            Id       string
            FilePath string
        }{}

        for rows.Next() {
            var id string
            var filePath string

            err = rows.Scan(&id, &filePath)
            if err != nil {
                tx.Rollback()
                fmt.Println("ShareBin: Error scanning row:", err)
                rows.Close()
                return
            }

            // Avoid cluttering with a type struct; collect IDs and file paths for deletion
            // We can't remove rows directly due to the lock; store them in an array
            toDelete = append(toDelete, struct {
                Id       string
                FilePath string
            }{Id: id, FilePath: filePath})
        }

        rows.Close()

        for _, v := range toDelete {
            _, err = tx.Exec("DELETE FROM data WHERE id = ?", v.Id)
            if err != nil {
                tx.Rollback()
                fmt.Println("ShareBin: Failed to delete share from database:", err)
                return
            }

            err = os.Remove(v.FilePath)
            if err != nil {
