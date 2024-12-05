package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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

type artisPageData struct {
	Artist artistInfo
	Gigs   [][2]string
}

type errorPageData struct {
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

	var dataAP artisPageData
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

	var status int
	errorMsg := ""
	dataAP.Gigs, status, errorMsg = getGigs(arti)
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
	http.HandleFunc("/groupie-tracker/about", homeHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
