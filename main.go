package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Holds the name of a country and if it's been selected
type countryInfo struct {
	Name     string
	Selected bool
}

// Filter selections and artist informations for home template
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
	Artist artistInfo
	Gigs   [][2]string
}

type ErrorPageData struct {
	Error    uint
	Message1 string
	Message2 string
}

var (
	firstLoad bool = true
	flt       filter
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// pageDataValues formats the data to be sent to the home template
func homePageDataValues(f filter, ais []artistInfo) HomePageData {

	cInfos := []countryInfo{}
	for i, boo := range f.countries {
		cInfos = append(cInfos, countryInfo{allCountries[i], boo})
	}

	data := HomePageData{
		Order:     f.order,
		BandCheck: f.band,
		SoloCheck: f.solo,
		CreMin:    strconv.Itoa(f.created[0]),
		CreMax:    strconv.Itoa(f.created[1]),
		FiAlMin:   strconv.Itoa(f.firstAl[0]),
		FiAlMax:   strconv.Itoa(f.firstAl[1]),
		PeMin:     strconv.Itoa(f.recPerf[0]),
		PeMax:     strconv.Itoa(f.recPerf[1]),
		Countries: cInfos,
		Artists:   ais,
		MinMax:    minmaxFirst,
	}
	return data
}

// handler for the homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/groupie-tracker" && r.URL.Path != "/groupie-tracker/about" {
		goToErrorPage(http.StatusNotFound, "Not Found", `Page doesn't exist`, w) // Error 404
		fmt.Println("Bad url path:", r.URL.Path)
		return
	}

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
	if r.URL.Path == "/groupie-tracker/about" {
		tmpl.ExecuteTemplate(w, "about.html", nil)
	} else {
		tmpl.ExecuteTemplate(w, "index.html", data)
	}

}

// artistHandler serves a site for a specific artist
func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.String()[strings.LastIndex(r.URL.String(), "=")+1:]
	artistID, err := strconv.Atoi(id)
	if err != nil {
		goToErrorPage(http.StatusBadRequest, "Bad Request", "Invalid artist ID: "+err.Error(), w) // Error 400
		return
	}

	if len(artInfos) == 0 { // In case someone navigates to an artist page directly
		readAPI(w)
	}

	var dataAP ArtisPageData
	var found1 bool
	for _, ai := range artInfos {
		if ai.Id == artistID {
			dataAP.Artist = ai
			found1 = true
			break
		}
	}
	if !found1 {
		goToErrorPage(http.StatusNotFound, "Not Found", "Artist "+id+` doesn't exist`, w) // Error 404
		return
	}

	var arti artist
	var found2 bool
	for _, a := range artists {
		if a.Id == artistID {
			arti = a
			found2 = true
			break
		}
	}
	if !found2 {
		goToErrorPage(http.StatusNotFound, "Not Found", "Artist "+id+` doesn't exist`, w) // Error 404
		return
	}

	dataAP.Gigs, err = getGigs(arti)
	if err != nil {
		goToErrorPage(http.StatusBadRequest, "Bad Request", "Failed to fetch data from API: "+err.Error(), w) // Error 400
		return
	}

	tmpl.ExecuteTemplate(w, "artistpage.html", dataAP)
}

// goToErrorPage handles errors by loading an error page to the user
func goToErrorPage(errorN int, m1 string, m2 string, w http.ResponseWriter) {
	w.WriteHeader(errorN)
	epd := ErrorPageData{uint(errorN), m1, m2}
	fmt.Printf("%d %s, %s\n", errorN, m1, m2)
	tmpl.ExecuteTemplate(w, "errorpage.html", epd)
}

func main() {
	//http.Handle("/static/", http.FileServer(http.Dir(".")))

	fileServer := http.FileServer(http.Dir("."))
	http.Handle("/static/styles.css", fileServer)
	http.Handle("/static/sad.jpg", fileServer)
	http.Handle("/static/guitarbrown.png", fileServer)
	http.Handle("/static/home-functions.js", fileServer)
	http.Handle("/static/ui-functions.js", fileServer)
	http.Handle("/favicon.ico", fileServer)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/groupie-tracker/artist/", artistHandler)
	http.HandleFunc("/groupie-tracker/about", homeHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
