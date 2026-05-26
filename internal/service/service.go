package service

import (
	"strings"

	"github.com/DanilaNova/go-6-sprint-final/pkg/morse"
)

func Convert(input string) string {
	if strings.ContainsFunc(input, isntMorse) {
		return morse.ToMorse(string(input))
	} else {
		return morse.ToText(string(input))
	}
}

func isntMorse(input rune) bool {
	switch input {
	case '.', '-', ' ':
		return false
	}
	return true
}
