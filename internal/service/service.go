package service

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	ErrCannotCreateDumpFolder = errors.New("Cannot create result dump folder")

	ErrCannotCreateDumpFile = errors.New("Cannot create result dump file")
	ErrStat                 = errors.New("Error in os.Stat")
)

// Checks whether byte slice is valid morse
func IsMorse(char []byte) bool {
	for _, ch := range char {
		switch ch {
		case '.', '-', ' ':
			continue
		default:
			return false
		}
	}
	return true
}

// Creates dump file with current time as a name and writes data into it
func CreateDump(data []byte) error {
	_, err := os.Stat("dumps")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("dumps", 0755)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrCannotCreateDumpFolder, err)
			}
		} else {
			return fmt.Errorf("%w: %w", ErrStat, err)
		}
	}

	file, err := os.Create("dumps/" + time.Now().UTC().String() + ".txt")
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreateDumpFile, err)
	}

	_, writeErr := file.Write(data)
	closeErr := file.Close()
	return errors.Join(writeErr, closeErr)
}
