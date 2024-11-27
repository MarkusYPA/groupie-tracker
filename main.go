package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// URLs from the initial API response
type APIResponse struct {
	ArtistsUrl   string `json:"artists"`
	LocationsUrl string `json:"locations"`
	DatesUrl     string `json:"dates"`
	RelationUrl  string `json:"relation"`
}

// Raw artist data from API
type artist struct {
	Id         int      `json:"id"`
	Image      string   `json:"image"`
	Name       string   `json:"name"`
	Members    []string `json:"members"`
	CreDate    int      `json:"creationDate"`
	FirstAlbum string   `json:"firstAlbum"`
	Locations  string   `json:"locations"`
	Relations  string   `json:"relations"`
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

// Parsed dates, formatted dates, and nicely spelled locations and countries
type dateWithGig struct {
	Date    time.Time
	DateStr string
	Locale  string
	Country string
}

// Combination of info from artist and relations with nice dates
type artistInfo struct {
	Id         int
	Name       string
	Image      string
	Members    []string
	CreDate    int
	FirstAlbum time.Time
	FAString   string
	Gigs       []dateWithGig
}

// Contains user selections
type filter struct {
	order     string
	created   [2]int
	firstAl   [2]int
	recPerf   [2]int
	band      bool
	solo      bool
	countries []bool
}

type countryInfo struct {
	Name     string
	Selected bool
}

// Filter selections and artistInfos for template to display
type HomePageData struct {
	Order     string
	BandCheck bool
	SoloCheck bool
	CreMin    string
	CreMax    string
	FiAlMin   string
	FiAlMax   string
	PeMin     string
	PeMax     string
	Countries []countryInfo
	Artists   []artistInfo
	MinMax    [6]int
}

// Info that gets displayed on the artist page
type ArtisPageData struct {
	Artist     artistInfo
	Members    []string
	FirstAlbum string
	Locations  []string
	Dates      []string
}

var (
	allCountries  []string
	apiData       APIResponse
	artists       []artist
	relationIndex relIndex
	artInfos      []artistInfo
	firstLoad     bool = true
	flt           filter
	minmaxFirst   [6]int
)

// handler for the homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if firstLoad {
		readAPI(w)
		flt = defaultFilter()
		firstLoad = false
	}

	if r.Method == http.MethodPost && r.FormValue("reset") != "rd" {
		flt = newFilter(r)
	}

	if r.FormValue("reset") == "rd" {
		flt = defaultFilter()
	}

	toDisplay := filterBy(flt, artInfos)
	data := homePageDataValues(flt, toDisplay)
	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, data)
}

// artistHandler serves a site for a specific artist
func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/artist/"):]
	artistID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error parsing id:", id)
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	var dataAP ArtisPageData
	for _, ai := range artInfos {
		if ai.Id == artistID {
			dataAP.Artist = ai
		}
	}

	for _, d := range fetchDates(apiData.DatesUrl).Index {
		if d.Id == artistID {
			dataAP.Dates = d.Dates
		}
	}

	for _, l := range fetchLocations(apiData.LocationsUrl).Index {
		if l.Id == artistID {
			dataAP.Locations = l.Locales
		}
	}

	for _, a := range artists {
		if a.Id == artistID {
			dataAP.Members = a.Members
			dataAP.FirstAlbum = a.FirstAlbum
		}
	}

	t := template.Must(template.ParseFiles("templates/artistpage.html"))
	t.Execute(w, dataAP)
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/static/styles.css", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist/", artistHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
