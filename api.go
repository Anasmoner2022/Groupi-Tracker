package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var apiClient = &http.Client{Timeout: 10 * time.Second}

var (
	ArtistCache   []Artist
	RelationCache map[int]RelationItem
)

func InitCache() error {
	// 1. Fetch Artists
	artists, err := getArtists()
	if err != nil {
		return err
	}
	ArtistCache = artists

	// 2. Fetch all Relations at once
	relData, err := getRelations()
	if err != nil {
		return err
	}

	// 3. Convert Relations slice into a Map for instant access by ID
	RelationCache = make(map[int]RelationItem)
	for _, item := range relData.Index {
		RelationCache[item.ID] = item
	}

	return nil
}

func getArtists() ([]Artist, error) {
	resp, err := apiClient.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected Status: %d", resp.StatusCode)
	}
	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return artists, nil
}

func getRelations() (*Relations, error) {
	resp, err := apiClient.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected Status: %d", resp.StatusCode)
	}
	var relations Relations
	if err := json.NewDecoder(resp.Body).Decode(&relations); err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return &relations, nil
}
