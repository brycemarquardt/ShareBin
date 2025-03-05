/*
This file is part of ShareBin.

ShareBin is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

ShareBin is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with ShareBin. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
    "archive/zip"
    "io"
    "mime/multipart"
    "os"
    "time"
    "log"
)

func MultipleFileWriter(files []*multipart.FileHeader, path string, aesKey []byte, callback func()) {
    outFile, err := os.Create(path)
    if err != nil {
        log.Printf("ShareBin: Failed to create multiple file output for path %s: %v", path, err)
        return
    }
    defer outFile.Close()

    // Make new zip file
    zipWriter := zip.NewWriter(outFile)
    defer zipWriter.Close()

    for _, fileHeader := range files {
        file, err := fileHeader.Open()
        if err != nil {
            log.Printf("ShareBin: Failed to open multiple file %s: %v", fileHeader.Filename, err)
            return
        }
        defer file.Close()

        // Create file inside the zip
        writer, err := zipWriter.Create(fileHeader.Filename)
        if err != nil {
            log.Printf("ShareBin: Failed to create zip entry for file %s: %v", fileHeader.Filename, err)
            return
        }

        buffer := make([]byte, 1024*Global.StreamSizeLimit)
        for {
            n, err := file.Read(buffer)
            if err != nil && err != io.EOF {
                log.Printf("ShareBin: Error reading multiple file chunk for file %s: %v", fileHeader.Filename, err)
                return
            }
            if n == 0 {
                break
            }

            // Write the chunk to the ZIP file
            _, err = writer.Write(buffer[:n])
            if err != nil {
                log.Printf("ShareBin: Failed to write to zip for file %s: %v", fileHeader.Filename, err)
                return
            }

            // Need to add check for > 0 because Sleep(0) will just trigger context switch
            if Global.StreamThrottle > 0 {
                time.Sleep(time.Duration(Global.StreamThrottle) * time.Millisecond)
            }
        }
    }

    zipWriter.Close()
    outFile.Close()

    if aesKey != nil {
        err = EncryptFile(path, aesKey)
        if err != nil {
            log.Printf("ShareBin: Failed to encrypt multiple file ZIP at path %s: %v", path, err)
        }
    }

    callback()
}

func SingleFileWriter(files []*multipart.FileHeader, path string, aesKey []byte, callback func()) {
    outFile, err := os.Create(path)
    if err != nil {
        log.Printf("ShareBin: Failed to create single file output for path %s: %v", path, err)
        return
    }
    defer outFile.Close()

    file, err := files[0].Open()
    if err != nil {
        log.Printf("ShareBin: Failed to open single file %s: %v", files[0].Filename, err)
        return
    }
    defer file.Close()

    // Use a buffered reader to read the file in chunks
    buffer := make([]byte, 1024*Global.StreamSizeLimit)
    for {
        n, err := file.Read(buffer)
        if err != nil && err != io.EOF {
            log.Printf("ShareBin: Error reading single file chunk for file %s: %v", files[0].Filename, err)
            return
        }
        if n == 0 {
            break
        }

        // Write the chunk to the output file
        _, err = outFile.Write(buffer[:n])
        if err != nil {
            log.Printf("ShareBin: Failed to write single file chunk for file %s: %v", files[0].Filename, err)
            return
        }

        // Need to add check for > 0 because Sleep(0) will just trigger context switch
        if Global.StreamThrottle > 0 {
            time.Sleep(time.Duration(Global.StreamThrottle) * time.Millisecond)
        }
    }

    outFile.Close()
    if aesKey != nil {
        err = EncryptFile(path, aesKey)
        if err != nil {
            log.Printf("ShareBin: Failed to encrypt single file at path %s: %v", path, err)
        }
    }

    callback()
}
