package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//go:embed templates/*.html
var templateFiles embed.FS

var (
	HomeTemplate   *template.Template
	ArtistTemplate *template.Template
)

func init() {
	HomeTemplate = template.Must(template.ParseFS(templateFiles, "templates/index.html"))
	ArtistTemplate = template.Must(template.ParseFS(templateFiles, "templates/artist.html"))
}

func HomeHandler(resp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(resp, "404 Not Found", http.StatusNotFound)
		return
	}

	if req.Method != "GET" {
		http.Error(resp, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := HomeTemplate.Execute(resp, ArtistCache); err != nil {
		log.Println("Template error:", err)
	}
}

func ArtistHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(resp, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil {
		http.Error(resp, "Bad Request", http.StatusBadRequest)
		return
	}

	if id < 1 || id > 52 {
		http.Error(resp, "Bad Request", http.StatusBadRequest)
		return
	}

	var selectedArtist *Artist
	for _, a := range ArtistCache {
		if a.ID == id {
			selectedArtist = &a
			break
		}
	}

	if selectedArtist == nil {
		http.Error(resp, "Artist Not Found", http.StatusNotFound)
		return
	}

	relation, exists := RelationCache[id]
	if !exists {
		log.Printf("Warning: No relation found for artist %d", id)
	}

	// 3. Prepare data for the template
	data := struct {
		Artist   *Artist
		Relation RelationItem
	}{
		Artist:   selectedArtist,
		Relation: relation,
	}

	if err := ArtistTemplate.Execute(resp, data); err != nil {
		log.Println("Template error:", err)
		http.Error(resp, "Internal Server Error", 500)
	}
}

func SearchHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(resp, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.ToLower(req.URL.Query().Get("q"))

	var results []Artist
	for _, a := range ArtistCache {
		// Check Artist Name
		nameMatch := strings.Contains(strings.ToLower(a.Name), query)

		// Check all Members
		memberMatch := false
		for _, m := range a.Members {
			if strings.Contains(strings.ToLower(m), query) {
				memberMatch = true
				break
			}
		}

		// Check Creation Date
		dateMatch := strings.Contains(strconv.Itoa(a.CreationDate), query)

		if nameMatch || memberMatch || dateMatch {
			results = append(results, a)
		}
	}
	if results == nil {
		results = []Artist{}
	}
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(results)
}
