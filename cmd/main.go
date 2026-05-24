package main

import (
	"log"
	"os"

	"github.com/DanilaNova/go-6-sprint-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "Http Server ", log.Default().Flags())
	httpServer := server.Create(logger)

	err := httpServer.Server.ListenAndServe()
	if err != nil {
		logger.Panicln(err)
	}
}
