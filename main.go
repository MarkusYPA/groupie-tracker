package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

	if r.URL.Path != "/" {
		goToErrorPage(http.StatusNotFound, "Not Found", `Page doesn't exist`, w) // Error 404
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
	t := template.Must(template.ParseFiles("templates/index.html", "templates/modebuttonscript.html", "templates/footer.html"))
	t.Execute(w, data)
}

// artistHandler serves a site for a specific artist
func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/artist/"):]
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

	t := template.Must(template.ParseFiles("templates/artistpage.html", "templates/modebuttonscript.html", "templates/footer.html"))
	t.Execute(w, dataAP)
}

// goToErrorPage handles errors by loading an error page to the user
func goToErrorPage(errorN int, m1 string, m2 string, w http.ResponseWriter) {
	errorTemplate := template.Must(template.ParseFiles("templates/errorpage.html", "templates/modebuttonscript.html", "templates/footer.html"))
	w.WriteHeader(errorN)
	epd := ErrorPageData{uint(errorN), m1, m2}
	errorTemplate.Execute(w, epd)
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/static/styles.css", http.StripPrefix("/static/", fileServer))
	http.Handle("/static/sad.jpg", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist/", artistHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
