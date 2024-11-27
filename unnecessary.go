package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// Raw location data from API
type locIndex struct {
	Index []locations `json:"index"`
}

// Stores data for locIndex, also straight from API
type locations struct {
	Id      int      `json:"id"`
	Locales []string `json:"locations"`
}

// Raw date data from API
type dtIndex struct {
	Index []dates `json:"index"`
}

// Stores data for dtIndex, also straight from API
type dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Function to fetch data from the "dates" API endpoint
func fetchDates(dateURL string) (dtIndex, error) {
	var dates dtIndex

	resp, err := http.Get(dateURL)
	if err != nil {
		return dates, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dates, err
	}

	// Parse JSON into Go struct
	err = json.Unmarshal(body, &dates)
	if err != nil {
		return dates, err
	}

	return dates, nil
}

// Function to fetch data from the "locations" API endpoint
func fetchLocations(locURL string) (locIndex, error) {
	var locs locIndex

	resp, err := http.Get(locURL)
	if err != nil {
		return locs, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return locs, err
	}

	// Parse JSON into Go struct
	err = json.Unmarshal(body, &locs)
	if err != nil {
		return locs, err
	}

	return locs, nil
}
