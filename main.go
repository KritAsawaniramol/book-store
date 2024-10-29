package main

import (
	"log"
	"os"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/pkg/database"
	"github.com/kritAsawaniramol/book-store/server"
)

func main() {
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	db := database.NewPostgresDatabase(cfg)
	
	server.NewGinServer(cfg, db.GetDb()).Start()
}
