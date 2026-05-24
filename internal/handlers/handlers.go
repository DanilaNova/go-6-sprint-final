package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/DanilaNova/go-6-sprint-final/internal/service"
	"github.com/DanilaNova/go-6-sprint-final/pkg/morse"
)

const (
	// 10 MB
	MULTIPART_PARSE_MAX_MEMORY = 10 * 1024 * 1024
	STATUS_CONTENT_TOO_LARGE   = 413
)

// "GET /" pattern handler generator
func Root(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Println("Root request error: index reading error: ", err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// "POST /upload" pattern handler generator
func Upload(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength > MULTIPART_PARSE_MAX_MEMORY {
			http.Error(w, fmt.Sprintf("Content Too Large (max %d)", MULTIPART_PARSE_MAX_MEMORY), STATUS_CONTENT_TOO_LARGE)
			logger.Println("Upload error: content too large: ", req.ContentLength)
			return
		}

		err := req.ParseMultipartForm(req.ContentLength)

		if errors.Is(err, http.ErrNotMultipart) {
			w.Header().Set("Accept", "multipart/form-data")
			http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
			logger.Println("Upload error: unsupported media type: ", req.Header.Get("Content-Type"))
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Println("Upload error: unknown parsing error: ", err)
			return
		}

		files, ok := req.MultipartForm.File["myFile"]
		if !ok || len(files) < 1 {
			http.Error(w, "Missing 'myFile' file from multipart form", http.StatusInternalServerError)
			logger.Println("Upload error: missing 'myFile' file from multipart form")
			return
		}

		file, err := files[0].Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Println("Upload error: file header opening error: ", err)
			return
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Println("Upload error: file reading error: ", err)
			return
		}

		isMorse := service.IsMorse(data)
		var isMorseStr string
		if isMorse {
			isMorseStr = "morse"
		} else {
			isMorseStr = "text"
		}

		var converted string
		if isMorse {
			converted = morse.ToText(string(data))
		} else {
			converted = morse.ToMorse(string(data))
		}
		logger.Printf("Got %s data: %q.\nConverted to: %q.", isMorseStr, string(data), converted)
		data = []byte(converted)

		err = service.CreateDump(data)
		if err != nil {
			logger.Printf("Error in creating dump: %v", err)
		}

		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
