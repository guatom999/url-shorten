package main

import (
	"log"
	"shorten-url/configs"
	"shorten-url/internal/databases"
	"shorten-url/internal/server"
)

// Main application logic

func main() {

	//test
	log.Println("Starting application...")

	cfg, err := configs.LoadConfig("../../.env")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db := databases.DatabaseConnect(cfg)

	log.Printf("Loaded config: Server will run on %s:%s", cfg.Server.Host, cfg.Server.Port)

	server.NewEchoServer(cfg, db).Start()
}
