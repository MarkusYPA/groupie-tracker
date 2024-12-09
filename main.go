package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type countryInfo struct {
	Name     string
	Selected bool
}

// Filter selections and artist informations for home template
type homePageData struct {
	Order     string
	BandCheck bool
	SoloCheck bool
	StartMin  string
	StartMax  string
	AlbumMin  string
	AlbumMax  string
	ShowMin   string
	ShowMax   string
	Countries []countryInfo
	Artists   []artistInfo
	MinMax    [6]int
}

type artistPageData struct {
	Artist artistInfo
	Gigs   [][2]string
}

type errorPageData struct {
	Error    uint
	Message1 string
	Message2 string
}

/* var (
	firstLoad bool = true
	flt       filter
) */

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// pageDataValues formats the data to be sent to the home template
func homePageDataValues(f filter, ais []artistInfo) homePageData {

	cInfos := []countryInfo{}
	for i, boo := range f.countries {
		cInfos = append(cInfos, countryInfo{allCountries[i], boo})
	}

	data := homePageData{
		Order:     f.order,
		BandCheck: f.band,
		SoloCheck: f.solo,
		StartMin:  strconv.Itoa(f.created[0]),
		StartMax:  strconv.Itoa(f.created[1]),
		AlbumMin:  strconv.Itoa(f.firstAl[0]),
		AlbumMax:  strconv.Itoa(f.firstAl[1]),
		ShowMin:   strconv.Itoa(f.recShow[0]),
		ShowMax:   strconv.Itoa(f.recShow[1]),
		Countries: cInfos,
		Artists:   ais,
		MinMax:    minmaxLimits,
	}
	return data
}

// handler for the homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/groupie-tracker" && r.URL.Path != "/groupie-tracker/about" {
		goToErrorPage(http.StatusNotFound, "Not Found", "Page doesn't exist", w) // Error 404
		fmt.Println("Bad URL path:", r.URL.Path)
		return
	}

	if r.URL.Path == "/groupie-tracker/about" {
		tmpl.ExecuteTemplate(w, "about.html", nil)
		return
	}

	err := readAPI(w)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var flt filter
	if r.Method == http.MethodPost && r.FormValue("reset") != "resetfilter" && minmaxLimits[0] != 0 {
		flt = newFilter(r)
	} else {
		flt = defaultFilter()
	}

	artistsToDisplay := filterBy(flt, artInfos)
	data := homePageDataValues(flt, artistsToDisplay)

	tmpl.ExecuteTemplate(w, "index.html", data)
}

// artistHandler serves a site for a specific artist
func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/groupie-tracker/artist/"):]
	artistID, err := strconv.Atoi(id)
	if err != nil {
		goToErrorPage(http.StatusBadRequest, "Bad Request", "Invalid artist ID: "+id, w)
		return
	}

	if len(artInfos) == 0 { // When navigating to an artist page directly
		readAPI(w)
	}

	var dataAP artistPageData
	var foundId bool
	for _, ai := range artInfos {
		if ai.Id == artistID {
			foundId = true
			dataAP.Artist = ai
			break
		}
	}

	if !foundId {
		goToErrorPage(http.StatusNotFound, "Not Found", "Artist "+id+` doesn't exist`, w)
		return
	}

	var status int
	var errorMsg string
	dataAP.Gigs, status, errorMsg = getGigs(dataAP.Artist)
	if status != http.StatusOK {
		goToErrorPage(http.StatusBadRequest, errorMsg, "Failed to fetch data from API", w) // Error 400
		return
	}

	tmpl.ExecuteTemplate(w, "artistpage.html", dataAP)
}

// goToErrorPage handles errors by loading an error page to the user
func goToErrorPage(errorN int, m1 string, m2 string, w http.ResponseWriter) {
	w.WriteHeader(errorN)
	epd := errorPageData{uint(errorN), m1, m2}
	fmt.Printf("%d %s, %s\n", errorN, m1, m2)
	tmpl.ExecuteTemplate(w, "errorpage.html", epd)
}

func main() {
	fileServer := http.FileServer(http.Dir("."))
	http.Handle("/static/css/styles.css", fileServer)
	http.Handle("/static/css/homepage.css", fileServer)
	http.Handle("/static/css/headerfooter.css", fileServer)
	http.Handle("/static/css/artistpage.css", fileServer)
	http.Handle("/static/css/darkmode.css", fileServer)
	http.Handle("/static/sad.jpg", fileServer)
	http.Handle("/static/guitar2.png", fileServer)
	http.Handle("/static/home-functions.js", fileServer)
	http.Handle("/static/ui-functions.js", fileServer)
	http.Handle("/favicon.ico", fileServer)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/groupie-tracker/artist/", artistHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
