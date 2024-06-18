package main

import (
	"Backend/api/songs"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /songs/next", songs.NextSongHandler)
	router.HandleFunc("POST /songs/add", songs.AddSongHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server listening on port :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}
