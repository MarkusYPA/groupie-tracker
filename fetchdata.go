package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
)

// URLs from the initial API response
type apiResponse struct {
	ArtistsUrl  string `json:"artists"`
	RelationUrl string `json:"relation"`
}

// Raw artist data from API
type artist struct {
	Id              int      `json:"id"`
	Image           string   `json:"image"`
	Name            string   `json:"name"`
	Members         []string `json:"members"`
	StartDate       int      `json:"creationDate"`
	FirstAlbum      string   `json:"firstAlbum"`
	LocationsUrl    string   `json:"locations"`
	ConcertDatesUrl string   `json:"concertDates"`
}

// Raw relation data from API
type relIndex struct {
	Index []relations `json:"index"`
}

// Stores data for relIndex, also straight from API
type relations struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Stores data for locIndex, also straight from API
type locations struct {
	Id      int      `json:"id"`
	Locales []string `json:"locations"`
}

// Stores data for dtIndex, also straight from API
type dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Parsed dates, formatted dates, and nicely spelled locations and countries
type dateWithGig struct {
	Date    time.Time
	Country string
	Locale  string
}

// Combination of info from artist and relations with nice dates
type artistInfo struct {
	Id              int
	Name            string
	Image           string
	Members         []string
	StartDate       int
	FirstAlbum      time.Time
	FAString        string
	Gigs            []dateWithGig
	LocationsUrl    string
	ConcertDatesUrl string
}

var (
	allCountries          []string
	allLocales            []string
	allCountryLocalePairs [][2]string
	apiData               apiResponse
	artInfos              []artistInfo
	artistsApi            []artist
	relationIndex         relIndex
)

// beautifyLocation returns the location and country of a concert, written all nicely
func beautifyLocation(s string) (string, string) {
	name := ""
	// separate location and country
	placeCountry := strings.Split(s, "-")
	for indexWd, wd := range placeCountry {
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
		if indexWd == 0 {
			name += ","
		}
	}
	return strings.Split(name, ",")[0], strings.Split(name, ",")[1]
}

// Function to fetch data from different API endpoints
func fetchFromAPI(relURL string, dataReciever interface{}) (int, string) {
	resp, err := http.Get(relURL)
	if err != nil {
		return http.StatusBadGateway, "Bad Gateway"
	}
	defer resp.Body.Close()

	// Parse JSON directly from the response body into the Go struct
	err = json.NewDecoder(resp.Body).Decode(&dataReciever)
	if err != nil {
		return http.StatusInternalServerError, "Internal Server Error"
	}

	return resp.StatusCode, ""
}

// getGigs retrieves and parses the dates, locations and countries for an artist's concerts
func getGigs(artistI artistInfo) ([][2]string, int, string) {
	gigs := [][2]string{}
	gigDates := []time.Time{}
	errorMessage := ""

	var loc locations
	status, errorMessage := fetchFromAPI(artistI.LocationsUrl, &loc)
	if status != http.StatusOK {
		return gigs, status, errorMessage
	}
	var dat dates
	status, errorMessage = fetchFromAPI(artistI.ConcertDatesUrl, &dat)
	if status != http.StatusOK {
		return gigs, status, errorMessage
	}

	localeIndex := -1
	for _, day := range dat.Dates {
		if day[0] == '*' {
			localeIndex++
			day = day[1:]
		}

		dat, err := time.Parse("02-01-2006", day)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}
		locale, cou := beautifyLocation(loc.Locales[localeIndex])
		dateStr := dat.Format("Jan. 2, 2006")
		if len(dateStr) > 4 && dateStr[6] == ',' {
			dateStr = dateStr[:4] + " " + dateStr[4:] // add space to single digit days so row lengths match
		}

		gigs = append(gigs, [2]string{dateStr, locale + ", " + cou})
		gigDates = append(gigDates, dat)
	}

	// Sort gigs from newest to oldest
	if len(gigs) == len(gigDates) {
		for i := 0; i < len(gigs)-1; i++ {
			for j := i + 1; j < len(gigs); j++ {
				if gigDates[i].Before(gigDates[j]) {
					gigs[i], gigs[j] = gigs[j], gigs[i]
					gigDates[i], gigDates[j] = gigDates[j], gigDates[i]
				}
			}
		}
	}

	return gigs, http.StatusOK, ""
}

// dateAndGig writes parsed dates, formatted dates and nicely spelled countries and locations to a slice of structs
func dateAndGig(rels map[string][]string) (dateGig []dateWithGig) {
	// parse time from string and combine with location
	for place, dates := range rels {
		for _, dateRaw := range dates {
			dat, err := time.Parse("02-01-2006", dateRaw)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				continue
			}
			loc, cou := beautifyLocation(place)
			dateGig = append(dateGig, dateWithGig{Date: dat, Country: cou, Locale: loc})
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
func getArtisInfo(art *artist, index int, ri *relIndex) (artistInfo, error) {
	ai := artistInfo{}
	ai.Id, ai.Name, ai.Image = art.Id, art.Name, art.Image
	ai.Members, ai.StartDate = art.Members, art.StartDate
	ai.LocationsUrl, ai.ConcertDatesUrl = art.LocationsUrl, art.ConcertDatesUrl

	albumDate, err := time.Parse("02-01-2006", art.FirstAlbum)
	if err != nil {
		fmt.Println("Error parsing date:", err)
	}
	ai.FirstAlbum = albumDate
	ai.FAString = albumDate.Format("January 2, 2006")

	ai.Gigs = dateAndGig(ri.Index[index].DatesLocations)

	return ai, err
}

// artistInformation combines the API information from artists and relations
func artistInformation(artistsApi *[]artist, relationsInd *relIndex) ([]artistInfo, error) {
	artInfos := []artistInfo{}
	for i := 0; i < len(*artistsApi); i++ {
		info, err := getArtisInfo(&(*artistsApi)[i], i, relationsInd)
		if err != nil {
			return artInfos, err
		}
		artInfos = append(artInfos, info)
	}
	return artInfos, nil
}

// fillAllCountries places all visited countries' and locales' names and their pairs on slices
func fillAllCountries(ais *[]artistInfo) {
	for _, ai := range *ais {
		for _, g := range ai.Gigs {
			foundC := false
			for _, c := range allCountries {
				if c == g.Country {
					foundC = true
				}
			}
			if !foundC {
				allCountries = append(allCountries, g.Country)
			}

			foundL := false
			for _, l := range allLocales {
				if l == g.Locale {
					foundL = true
				}
			}
			if !foundL {
				allLocales = append(allLocales, g.Locale)
				allCountryLocalePairs = append(allCountryLocalePairs, [2]string{g.Country, g.Locale})
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

	// Sort locales slice alphabetically
	for i := 0; i < len(allLocales)-1; i++ {
		for j := i + 1; j < len(allLocales); j++ {
			if allLocales[i] > allLocales[j] {
				allLocales[i], allLocales[j] = allLocales[j], allLocales[i]
			}
		}
	}
}

// readAPI gets the data from the given API and stores it into some global variables
func readAPI(w http.ResponseWriter) error {
	var status int
	var errorMessage string
	var err error

	status, errorMessage = fetchFromAPI("https://groupietrackers.herokuapp.com/api", &apiData)
	if status != http.StatusOK {
		goToErrorPage(status, errorMessage, "Error reading API", w)
		return fmt.Errorf("error parsing API JSON")
	}

	status, errorMessage = fetchFromAPI(apiData.ArtistsUrl, &artistsApi)
	if status != http.StatusOK {
		goToErrorPage(status, errorMessage, "Error reading artist API", w)
		return fmt.Errorf("error reading artist API")
	}

	status, errorMessage = fetchFromAPI(apiData.RelationUrl, &relationIndex)
	if status != http.StatusOK {
		goToErrorPage(status, errorMessage, "Error reading relations API", w)
		return fmt.Errorf("error reading relations API")
	}

	artInfos, err = artistInformation(&artistsApi, &relationIndex)
	if err != nil { // Error parsing date
		goToErrorPage(http.StatusInternalServerError, "Internal Server Error", err.Error(), w)
		return fmt.Errorf("internal Server Error")
	}
	fillAllCountries(&artInfos)
	return nil
}
