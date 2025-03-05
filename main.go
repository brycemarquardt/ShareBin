/*
This file is part of ShareBin.

ShareBin is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

ShareBin is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with ShareBin. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
    "embed"
    "fmt"
    "html/template"
    "log"
    "mime"
    _ "modernc.org/sqlite"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "strconv"
    "strings"
    "syscall"
    "github.com/skip2/go-qrcode"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed data/settings.json
var settingsFile string

// ServeFile serves static files or passes to the next handler
func serveFile(w http.ResponseWriter, r *http.Request, next func(w2 http.ResponseWriter, r2 *http.Request)) {
    path := strings.TrimPrefix(r.URL.Path, "/")
    if path == "" {
        http.Redirect(w, r, "/index.html", http.StatusFound)
        return
    }
    file, err := staticFiles.ReadFile("static/" + path)
    if err != nil {
        next(w, r)
        return
    }
    ext := filepath.Ext(path)
    contentType := mime.TypeByExtension(ext)
    if contentType == "" {
        contentType = "application/octet-stream
