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
        contentType = "application/octet-stream"
    }
    w.Header().Set("Content-Type", contentType)
    w.Header().Set("Content-Length", strconv.Itoa(len(file)))
    w.WriteHeader(http.StatusOK)
    w.Write(file)
}

// DashboardHandler renders the shares dashboard
func dashboardHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if !ValidateSession(w, r) && Global.EnablePassword {
        http.Redirect(w, r, "/auth.html", http.StatusFound)
        return
    }

    shares, err := ListShares(db)
    if err != nil {
        http.Error(w, "Error listing shares", http.StatusInternalServerError)
        return
    }

    if typeFilter := r.URL.Query().Get("type"); typeFilter != "" {
        shares = FilterSharesByType(shares, typeFilter)
    }

    sortBy := r.URL.Query().Get("sort")
    order := r.URL.Query().Get("order")
    if sortBy != "" {
        SortShares(shares, sortBy, order)
    }

    tmpl, err := template.ParseFiles("templates/dashboard.html")
    if err != nil {
        http.Error(w, "Error parsing template", http.StatusInternalServerError)
        return
    }
    data := struct {
        Shares []Share
    }{Shares: shares}
    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}

func main() {
    if _, err := os.Stat("./uploads/"); os.IsNotExist(err) {
        err := os.MkdirAll("./uploads", os.ModePerm)
        if err != nil {
            fmt.Println(err)
        }
    }

    if _, err := os.Stat("./data/"); os.IsNotExist(err) {
        err := os.MkdirAll("./data", os.ModePerm)
        if err != nil {
            fmt.Println(err)
        }
        err = os.WriteFile("./data/settings.json", []byte(settingsFile), 0644)
        if err != nil {
            fmt.Println(err)
            return
        }
    }

    db := InitDatabase()
    InitSettings()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            if (r.URL.Path == "/" || r.URL.Path == "/index.html") && !ValidateSession(w, r) {
                http.Redirect(w, r, "/auth.html", http.StatusFound)
                return
            }
            serveFile(w, r, func(w2 http.ResponseWriter, r2 *http.Request) {
                DownloadHandler(w2, r2, db)
            })
        }

        if r.Method == http.MethodPost {
            r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*Global.FileSizeLimit)
            err := r.ParseMultipartForm(1024 * Global.StreamSizeLimit)
            defer func() {
                if r.MultipartForm != nil {
                    err := r.MultipartForm.RemoveAll()
                    if err != nil {
                        fmt.Println(err)
                    }
                }
            }()
            if err != nil {
                if err.Error() == "http: request body too large" {
                    http.Error(w, "SizeExceeded", http.StatusInternalServerError)
                    return
                }
                fmt.Println(err)
                return
            }
            if !ValidateSession(w, r) {
                if len(r.MultipartForm.Value["auth"]) > 0 {
                    auth := r.MultipartForm.Value["auth"][0]
                    if auth != Global.Password {
                        return
                    }
                } else {
                    return
                }
            }
            FileHandler(w, r, db)
        }
    })

    http.HandleFunc("/postText", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            if !ValidateSession(w, r) {
                return
            }
            r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*Global.TextSizeLimit)
            TextHandler(w, r, db)
        }
    })

    http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            AuthHandler(w, r)
        }
    })

    http.HandleFunc("/deleteSession", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            DeleteSession(w, r)
        }
    })

    http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
        dashboardHandler(w, r, db)
    })

    go CheckExpiration(db)

    server := &http.Server{Addr: ":80"}
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
        fmt.Println("Shutting down ShareBin")
        if err := server.Close(); err != nil {
            fmt.Println(err)
        }
        if err := db.Close(); err != nil {
            fmt.Println(err)
        }
    }()

    log.Println("ShareBin server running")
    log.Fatal(server.ListenAndServe())
}
