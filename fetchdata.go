package main

import (
	"encoding/json"
	"fmt"
	"io"
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

// Function to fetch data from an artists' locations url at the "locations" API endpoint
func fetchLocation(relURL string) (locations, error) {
	var loc locations
	resp, err := http.Get(relURL)
	if err != nil {
		return loc, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return loc, err
	}

	// Parse JSON into Go struct
	err = json.Unmarshal(body, &loc)
	if err != nil {
		return loc, err
	}

	return loc, nil
}

// Function to fetch data from an artists' dates url at the "dates" API endpoint
func fetchDate(relURL string) (dates, error) {
	var dat dates
	resp, err := http.Get(relURL)
	if err != nil {
		return dat, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dat, err
	}

	// Parse JSON into Go struct
	err = json.Unmarshal(body, &dat)
	if err != nil {
		return dat, err
	}

	return dat, nil
}

func getGigs(artist artist) ([][2]string, error) {
	gigs := [][2]string{}

	loc, e2 := fetchLocation(artist.Locations)
	if e2 != nil {
		return gigs, e2
	}
	dat, e3 := fetchDate(artist.ConcertDates)
	if e3 != nil {
		return gigs, e3
	}

	localeIndex := -1
	for _, day := range dat.Dates {
		daymod := day
		if day[0] == '*' {
			localeIndex++
			daymod = day[1:]
		}

		dat, err := time.Parse("02-01-2006", daymod)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}
		locale, cou := beautifyLocation(loc.Locales[localeIndex])
		dateStr := dat.Format("Jan. 2, 2006")

		gigs = append(gigs, [2]string{dateStr, locale + ", " + cou})
	}

	return gigs, nil
}

// getArtisInfo puts all the API info about an artist to a struct
func getArtisInfo(art artist, index int, ri relIndex) (artistInfo, error) {
	ai := artistInfo{}
	ai.Id, ai.Name, ai.Image = art.Id, art.Name, art.Image
	ai.Members, ai.CreDate = art.Members, art.CreDate
	albumDate, err := time.Parse("02-01-2006", art.FirstAlbum)
	if err != nil {
		fmt.Println("Error parsing date:", err)
	}
	ai.FirstAlbum = albumDate
	ai.FAString = albumDate.Format("January 2, 2006")

	ai.Gigs = dateAndGig(ri.Index[index].DatesLocations)

	return ai, nil
}

// artistInformation combines the API information from artists and relations
func artistInformation(artists []artist, rI relIndex) ([]artistInfo, error) {
	artInfos := []artistInfo{}
	for i := 0; i < len(artists); i++ {
		info, err := getArtisInfo(artists[i], i, rI)
		if err != nil {
			return artInfos, err
		}
		artInfos = append(artInfos, info)

	}
	return artInfos, nil
}

// Function to fetch data from the "artists" API endpoint
func fetchArtists(artistsURL string) ([]artist, error) {
	resp, err := http.Get(artistsURL)
	if err != nil {
		return artists, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return artists, err
	}

	// Parse JSON into Go struct
	var artists []artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return artists, err
	}

	return artists, nil
}

// Function to fetch data from the "relations" API endpoint
func fetchRelations(relURL string) (relIndex, error) {
	var rels relIndex
	resp, err := http.Get(relURL)
	if err != nil {
		return rels, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rels, err
	}

	// Parse JSON into Go struct
	err = json.Unmarshal(body, &rels)
	if err != nil {
		return rels, err
	}

	return rels, nil
}

// fetchAPI parses JSON into a Go struct to extract URLs
func fetchAPI(body []byte) (APIResponse, error) {
	var apiData APIResponse
	err := json.Unmarshal(body, &apiData)
	return apiData, err
}

// goToErrorPage handles errors by loading an error page to the user
func goToErrorPage(errorN int, m1 string, m2 string, w http.ResponseWriter) {
	w.WriteHeader(errorN)
	epd := ErrorPageData{uint(errorN), m1, m2}
	errorTemplate.Execute(w, epd)
}

// readAPI gets the data from the given API and stores it into some global variables
func readAPI(w http.ResponseWriter) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Failed to fetch data from API", w)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Error reading response", w)
		return
	}

	apiData, err = fetchAPI(body)
	if err != nil {
		goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Error parsing API JSON: "+err.Error(), w)
		return
	}
	artists, err = fetchArtists(apiData.ArtistsUrl)
	if err != nil {
		goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Error reading artist API: "+err.Error(), w)
		return
	}
	relationIndex, err = fetchRelations(apiData.RelationUrl)
	if err != nil {
		goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Error reading relations API: "+err.Error(), w)
		return
	}
	artInfos, err = artistInformation(artists, relationIndex)
	if err != nil {
		goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Failed to fetch data from API", w)
		return
	}
	fillAllCountries(artInfos)
}
