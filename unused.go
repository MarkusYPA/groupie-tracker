package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type locations struct {
	Id      int      `json:"id"`
	Locales []string `json:"locations"`
}

type locIndex struct {
	Index []locations `json:"index"`
}

type dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type dtIndex struct {
	Index []dates `json:"index"`
}

//dateIndex := fetchDates(apiData.DatesUrl)
//locationIndex := fetchLocations(apiData.LocationsUrl)

func fetchDates(dateURL string) dtIndex {
	resp, err := http.Get(dateURL)
	if err != nil {
		log.Println("Failed to fetch dates:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading dates response:", err)
		os.Exit(1)
	}

	// Parse JSON into Go struct
	var dates dtIndex
	err = json.Unmarshal(body, &dates)
	if err != nil {
		panic(err.Error())
	}

	return dates
}

func fetchLocations(locURL string) locIndex {
	resp, err := http.Get(locURL)
	if err != nil {
		log.Println("Failed to fetch locations:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading locations response:", err)
		os.Exit(1)
	}

	// Parse JSON into Go struct
	var locs locIndex
	err = json.Unmarshal(body, &locs)
	if err != nil {
		panic(err.Error())
	}

	return locs
}
