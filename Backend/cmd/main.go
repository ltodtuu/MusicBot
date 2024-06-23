package main

import (
	"Backend/api/songs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Port          int         `json:"port"`
	TidalSettings TidalConfig `json:"tidal_settings"`
}

type TidalConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	config, errCfg := LoadConfig("appsettings.json")
	if errCfg != nil {
		log.Fatalf("Failed to load configuration: %s", errCfg)
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /songs/next", songs.NextSongHandler)
	router.HandleFunc("GET /songs/queue", songs.ReadQueueHandler)
	router.HandleFunc("POST /songs/add", songs.AddSongHandler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	log.Println("Server listening on port :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}
