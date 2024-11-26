package main

import (
	"fmt"
	"html/template"
	"io"
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

type relIndex struct {
	Index []relations `json:"index"`
}

type relations struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type dateWithGig struct {
	Date    time.Time
	DateStr string
	Locale  string
	Country string
}

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

type filter struct {
	order     string
	created   [2]int
	firstAl   [2]int
	recPerf   [2]int
	band      bool
	solo      bool
	countries []bool
}

type PageData struct {
	Order     string
	BandCheck bool
	SoloCheck bool
	CreMin    string
	CreMax    string
	FiAlMin   string
	FiAlMax   string
	PeMin     string
	PeMax     string
	Countries []bool
	Artists   []artistInfo
}

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
)

// artistInformation combines the API information from artists and relations
func artistInformation(artists []artist, rI relIndex) []artistInfo {
	artInfos := []artistInfo{}
	for i := 0; i < len(artists); i++ {
		artInfos = append(artInfos, getArtisInfo(artists[i], i, rI))

	}
	return artInfos
}

// pageDataValues formats the data to be sent to the template
func pageDataValues(f filter, ais []artistInfo) PageData {
	data := PageData{
		Order:     f.order,
		BandCheck: f.band,
		SoloCheck: f.solo,
		CreMin:    strconv.Itoa(f.created[0]),
		CreMax:    strconv.Itoa(f.created[1]),
		FiAlMin:   strconv.Itoa(f.firstAl[0]),
		FiAlMax:   strconv.Itoa(f.firstAl[1]),
		PeMin:     strconv.Itoa(f.recPerf[0]),
		PeMax:     strconv.Itoa(f.recPerf[1]),
		Countries: f.countries,
		Artists:   ais,
	}
	return data
}

// newFilter places the user's selections to a filter
func newFilter(r *http.Request) filter {
	ord := r.FormValue("order")
	showBand := r.FormValue("band") == "on"
	showSolo := r.FormValue("solo") == "on"
	formMin, _ := strconv.Atoi(r.FormValue("fomin"))
	formMax, _ := strconv.Atoi(r.FormValue("fomax"))
	if formMax < formMin {
		formMax = formMin
	}
	fAMin, _ := strconv.Atoi(r.FormValue("famin"))
	fAMax, _ := strconv.Atoi(r.FormValue("famax"))
	if fAMax < fAMin {
		fAMax = fAMin
	}
	peMin, _ := strconv.Atoi(r.FormValue("pemin"))
	peMax, _ := strconv.Atoi(r.FormValue("pemax"))
	if peMax < peMin {
		peMax = peMin
	}

	countries := make([]bool, len(allCountries))
	for i, c := range allCountries {
		countries[i] = (r.FormValue(c) == "on" || r.Method == http.MethodGet)
	}

	return filter{
		order:     ord,
		created:   [2]int{formMin, formMax},
		firstAl:   [2]int{fAMin, fAMax},
		recPerf:   [2]int{peMin, peMax},
		band:      showBand,
		solo:      showSolo,
		countries: countries,
	}
}

func defaultFilter(r *http.Request) filter {

	countries := make([]bool, len(allCountries))
	for i := range allCountries {
		countries[i] = true
	}

	ord := "namedown"
	showBand := true
	showSolo := true
	formMin := 1950
	formMax := 2024
	fAMin := 1950
	fAMax := 2024
	peMin := 1950
	peMax := 2024

	return filter{
		order:     ord,
		created:   [2]int{formMin, formMax},
		firstAl:   [2]int{fAMin, fAMax},
		recPerf:   [2]int{peMin, peMax},
		band:      showBand,
		solo:      showSolo,
		countries: countries,
	}
}

// filterBy removes all artists that don't pass the filter
func filterBy(fil filter, arInfos []artistInfo) []artistInfo {
	aisOut := []artistInfo{}

	for _, ai := range arInfos {
		passes := true
		if ai.CreDate < fil.created[0] || ai.CreDate > fil.created[1] {
			passes = false
		}
		if ai.FirstAlbum.Year() < fil.firstAl[0] || ai.FirstAlbum.Year() > fil.firstAl[1] {
			passes = false
		}
		if ai.Gigs[0].Date.Year() < fil.recPerf[0] || ai.Gigs[0].Date.Year() > fil.recPerf[1] {
			passes = false
		}
		if !fil.band && len(ai.Members) > 1 {
			passes = false
		}
		if !fil.solo && len(ai.Members) == 1 {
			passes = false
		}

		countryNames := []string{}
		for i := 0; i < len(allCountries); i++ {
			if fil.countries[i] {
				countryNames = append(countryNames, allCountries[i])
			}
		}
		found := false
		for _, cn := range countryNames {
			for _, g := range ai.Gigs {
				if g.Country == cn {
					found = true
					break // from inner loop
				}
			}
			if found {
				break //
			}
		}

		if passes && found {
			aisOut = append(aisOut, ai)
		}
	}

	sortArtists(&aisOut, fil.order)

	return aisOut
}

// compare compares two artistInfos according to an attribute specified in string s
func compare(a1, a2 artistInfo, s string) bool {
	switch s {
	case "namedown":
		return a1.Name <= a2.Name
	case "fodown":
		return a1.CreDate <= a2.CreDate
	case "fadown":
		return a1.FirstAlbum.Before(a2.FirstAlbum)
	case "perdown":
		return a1.Gigs[0].Date.Before(a2.Gigs[0].Date)
	case "nameup":
		return a1.Name > a2.Name
	case "foup":
		return a1.CreDate > a2.CreDate
	case "faup":
		return a1.FirstAlbum.After(a2.FirstAlbum)
	case "perup":
		return a1.Gigs[0].Date.After(a2.Gigs[0].Date)
	}
	return true
}

// sortArtists sorts a slice of artistInfo according to the instruction in a string ord
func sortArtists(as *[]artistInfo, ord string) {
	for i := 0; i < len(*as)-1; i++ {
		for j := i + 1; j < len(*as); j++ {
			if !compare((*as)[i], (*as)[j], ord) {
				(*as)[i], (*as)[j] = (*as)[j], (*as)[i]
			}
		}
	}
}

func readAPI(w http.ResponseWriter) {

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	apiData = fetchAPI(body)
	artists = fetchArtists(apiData.ArtistsUrl)
	relationIndex = fetchRelations(apiData.RelationUrl)
	artInfos = artistInformation(artists, relationIndex)
	fillAllCountries(artInfos)
}

// handler for the homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {

	if firstLoad {
		readAPI(w)
		flt = defaultFilter(r)
		firstLoad = false
	}

	if r.Method == http.MethodPost && r.FormValue("reset") != "rd" {
		flt = newFilter(r)
	}

	if r.FormValue("reset") == "rd" {
		flt = defaultFilter(r)
	}

	toDisplay := filterBy(flt, artInfos)
	data := pageDataValues(flt, toDisplay)

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
	//readAPI(w)		// Only necessary on first Get?

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
