package songs

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
)

type Song struct {
	Title      string  `json:"title"`
	Artist     string  `json:"artist"`
	Album      string  `json:"album"`
	SpotifyURI url.URL `json:"spotify_uri"`
	TidalURI   url.URL `json:"tidal_uri"`
}

type SongQueue struct {
	songs []Song
	mutex sync.Mutex
}

var queue = &SongQueue{}

func (q *SongQueue) GetNextSong() (Song, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.songs) == 0 {
		return Song{}, false
	}
	nextSong := q.songs[0]
	q.songs = q.songs[1:]
	return nextSong, true
}

func (q *SongQueue) AddSong(song Song) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.songs = append(q.songs, song)
}

func (q *SongQueue) ReadSongQueue() []Song {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.songs
}

func ReadQueueHandler(w http.ResponseWriter, r *http.Request) {
	songs := queue.ReadSongQueue()
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddSongHandler(w http.ResponseWriter, r *http.Request) {
	var song Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queue.AddSong(song)
	w.WriteHeader(http.StatusCreated)
}

func NextSongHandler(w http.ResponseWriter, r *http.Request) {
	song, ok := queue.GetNextSong()
	if !ok {
		http.Error(w, "No songs in queue", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(song); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
