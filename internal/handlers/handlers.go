package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DanilaNova/go-6-sprint-final/internal/service"
)

const (
	// 10 MB
	MultipartParseMaxMemory = 10 * 1024 * 1024
	StatusContentTooLarge   = 413
)

var (
	errCannotCreateDumpFolder = errors.New("cannot create result dump folder")

	errCannotCreateDumpFile = errors.New("cannot create result dump file")
	errStat                 = errors.New("error in os.Stat")
)

// "GET /" pattern handler generator
//
// # Request
//
// irrelevant
//
// # Response
//
// Content-Type: text/html
func Root(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "index.html")
	}
}

// "POST /upload" pattern handler generator
//
// # Request
//
// Content-Type: multipart/form-data
//
// # Response
//
// Content-Type: text/plain;charset=utf-8
func Upload(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength > MultipartParseMaxMemory {
			http.Error(w, fmt.Sprintf("Content Too Large (max %d)", MultipartParseMaxMemory), http.StatusInternalServerError)
			logger.Println("Upload error: content too large: ", req.ContentLength)
			return
		}

		err := req.ParseMultipartForm(req.ContentLength)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Println("Upload error: multipart form parsing error: ", err)
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

		data, readErr := io.ReadAll(file)
		closeErr := file.Close()
		if err = errors.Join(readErr, closeErr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Println("Upload error: file header reading error: ", err)
			return
		}

		converted := service.Convert(string(data))

		logger.Printf("Got data: %q.\nConverted to: %q.", string(data), converted)
		data = []byte(converted)

		err = createDump(data)
		if err != nil {
			logger.Printf("Error in creating dump: %v", err)
		}

		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Creates dump file with current time as a name and writes data into it
//
// # Errors:
//
// [errCannotCreateDumpFolder] |
// [errCannotCreateDumpFile] |
// [errStat]
func createDump(data []byte) error {
	_, err := os.Stat("dumps")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("dumps", 0755)
			if err != nil {
				return fmt.Errorf("%w: %w", errCannotCreateDumpFolder, err)
			}
		} else {
			return fmt.Errorf("%w: %w", errStat, err)
		}
	}

	file, err := os.Create("dumps/" + time.Now().UTC().String() + ".txt")
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotCreateDumpFile, err)
	}

	_, writeErr := file.Write(data)
	closeErr := file.Close()
	return errors.Join(writeErr, closeErr)
}
