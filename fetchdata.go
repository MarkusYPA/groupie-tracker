package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"
)

func fillAllCountries(ais []artistInfo) {
	// Place all country names on slice
	for _, ai := range ais {
		for _, g := range ai.Gigs {
			found := false
			for _, c := range allCountries {
				if c == g.Country {
					found = true
				}
			}
			if !found {
				allCountries = append(allCountries, g.Country)
			}
		}
	}
	// Sort slice alphabetically
	for i := 0; i < len(allCountries)-1; i++ {
		for j := i + 1; j < len(allCountries); j++ {
			if allCountries[i] > allCountries[j] {
				allCountries[i], allCountries[j] = allCountries[j], allCountries[i]
			}
		}
	}
}

// beautifyLocation returns the location and country of a concert, written all nicely
func beautifyLocation(s string) (string, string) {
	name := ""
	// separate location and country
	placeCountry := strings.Split(s, "-")
	for iWd, wd := range placeCountry {
		if wd == "usa" || wd == "uk" {
			name += strings.ToUpper(wd)
			continue
		}
		for i := 0; i < len(wd); i++ {
			r := rune(wd[i])

			if unicode.IsLetter(r) {
				// Don't capitalize a middle word "del"
				if i != 0 && i < len(wd)-4 && wd[i-1:i+4] == "_del_" {
					name += "del "
					i += 3
					continue
				}
				// Don't capitalize a middle word "on"
				if i != 0 && i < len(wd)-3 && wd[i-1:i+3] == "_on_" {
					name += "on "
					i += 2
					continue
				}
				// Don't capitalize a middle word "de"
				if i != 0 && i < len(wd)-3 && wd[i-1:i+3] == "_de_" {
					name += "de "
					i += 2
					continue
				}
				if i == 0 || (i > 0 && wd[i-1] == '_') {
					name += strings.ToUpper(string(r))
				} else {
					name += string(r)
				}
			} else {
				name += " "
			}
		}
		if iWd == 0 {
			name += ","
		}
	}
	return strings.Split(name, ",")[0], strings.Split(name, ",")[1]
}

// dateAndGig writes parsed dates, formatted dates and nicely spelled countries and locations to a slice of structs
func dateAndGig(rels map[string][]string) (dateGig []dateWithGig) {
	// parse time from string and combine with location
	for place, sli := range rels {
		for _, dateRaw := range sli {
			dat, err := time.Parse("02-01-2006", dateRaw)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				continue
			}
			loc, cou := beautifyLocation(place)
			dateGig = append(dateGig, dateWithGig{Date: dat, DateStr: dat.Format("Jan. 2, 2006"), Locale: loc, Country: cou})
		}
	}

	// Put most recent gigs first
	for i := 0; i < len(dateGig)-1; i++ {
		for j := i + 1; j < len(dateGig); j++ {
			if dateGig[i].Date.Before(dateGig[j].Date) {
				dateGig[i], dateGig[j] = dateGig[j], dateGig[i]
			}
		}
	}

	return
}

// getArtisInfo puts all the API info about an artist to a struct
func getArtisInfo(art artist, index int, ri relIndex) artistInfo {
	ai := artistInfo{}
	ai.Id, ai.Name, ai.Image = art.Id, art.Name, art.Image
	ai.Members, ai.CreDate = art.Members, art.CreDate
	albumDate, err := time.Parse("02-01-2006", art.FirstAlbum)
	if err != nil {
		log.Fatalln("Error parsing date:", err)
	}
	ai.FirstAlbum = albumDate
	ai.FAString = albumDate.Format("January 2, 2006")
	ai.Gigs = dateAndGig(ri.Index[index].DatesLocations)

	return ai
}

// fetchAPI parses JSON into a Go struct to extract URLs
func fetchAPI(body []byte) APIResponse {
	var apiData APIResponse
	err := json.Unmarshal(body, &apiData)
	if err != nil {
		log.Fatalln("Error parsing API JSON", err)
	}

	return apiData
}

// Function to fetch data from the "artists" API endpoint
func fetchArtists(artistsURL string) []artist {
	resp, err := http.Get(artistsURL)
	if err != nil {
		log.Fatalln("Failed to fetch artists:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading artists response:", err)
	}

	// Parse JSON into Go struct
	var artists []artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		panic(err.Error())
	}

	return artists
}

// Function to fetch data from the "realations" API endpoint
func fetchRelations(relURL string) relIndex {
	resp, err := http.Get(relURL)
	if err != nil {
		log.Fatalln("Failed to fetch relations:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading relations response:", err)
	}

	// Parse JSON into Go struct
	var rels relIndex
	err = json.Unmarshal(body, &rels)
	if err != nil {
		panic(err.Error())
	}

	return rels
}
