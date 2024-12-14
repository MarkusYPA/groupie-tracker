package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type countrySelect struct {
	Name     string
	Selected bool
}

type localeSelect struct {
	Name     string
	Selected bool
	Display  string // Show checkbox in filter or don't
}

type memNumSelect struct {
	Name     string
	Selected bool
}

// Filter selections and artist informations for home template
type homePageData struct {
	Order           string
	StartMin        string
	StartMax        string
	AlbumMin        string
	AlbumMax        string
	ShowMax         string
	ShowYearMarkers []int
	Countries       []countrySelect
	Locales         []localeSelect
	MemNums         []memNumSelect
	Artists         []artistInfo
	MinMax          [6]int
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

var (
	tmplArtist  *template.Template
	tmplIndex   *template.Template
	tmplError   *template.Template
	tmplAbout   *template.Template
	apiReadTime time.Time
)

func init() {
	var err error
	tmplArtist, err = template.ParseFiles("templates/artistpage.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
	tmplIndex, err = template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
	tmplError, err = template.ParseFiles("templates/errorpage.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
	tmplAbout, err = template.ParseFiles("templates/about.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

// pageDataValues formats the data to be executed with the home template so filter
// selections are retained and filtered artists displayed
func homePageDataValues(f filter, ais []artistInfo) homePageData {

	couSels := []countrySelect{}
	for i, boo := range f.countries {
		couSels = append(couSels, countrySelect{allCountries[i], boo})
	}

	selectedCountries := []string{}
	for i := 0; i < len(allCountries); i++ {
		if f.countries[i] {
			selectedCountries = append(selectedCountries, allCountries[i])
		}
	}

	// numbers of members to display
	memNumSels := []memNumSelect{}
	for i, boo := range f.numbsOfMembers {
		memNumSels = append(memNumSels, memNumSelect{strconv.Itoa(allMemberNumbers[i]), boo})
	}

	locSels := []localeSelect{}
	for i, boo := range f.locales {
		display := ""
		// display only locales from selected countries
		if isLocaleInCountries(allLocales[i], selectedCountries) {
			display = `style="display: initial`
		} else {
			display = `style="display: none`
		}
		locSels = append(locSels, localeSelect{allLocales[i], boo, display})
	}

	// Create markers for the slider, rounded to divisible by 10
	firstMark := minmaxLimits[4] + (10-(minmaxLimits[4]%10))%10 // round up
	marks := []int{}
	for i := firstMark; i <= minmaxLimits[5]; i += 10 {
		marks = append(marks, i)
	}

	data := homePageData{
		Order:           f.order,
		StartMin:        strconv.Itoa(f.created[0]),
		StartMax:        strconv.Itoa(f.created[1]),
		AlbumMin:        strconv.Itoa(f.firstAl[0]),
		AlbumMax:        strconv.Itoa(f.firstAl[1]),
		ShowMax:         strconv.Itoa(f.recShow),
		ShowYearMarkers: marks,
		Countries:       couSels,
		Locales:         locSels,
		MemNums:         memNumSels,
		Artists:         ais,
		MinMax:          minmaxLimits,
	}
	return data
}

// handler for the homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Allow only GET and POST methods
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		goToErrorPage(http.StatusMethodNotAllowed, "Method Not Allowed", "Only GET and POST methods are allowed", w) // Error 405
		return
	}

	if r.URL.Path != "/" && r.URL.Path != "/groupie-tracker" && r.URL.Path != "/groupie-tracker/about" {
		goToErrorPage(http.StatusNotFound, "Not Found", "Page doesn't exist", w) // Error 404
		fmt.Println("Bad URL path:", r.URL.Path)
		return
	}

	if r.URL.Path == "/groupie-tracker/about" {
		if tmplAbout != nil {
			err := tmplAbout.Execute(w, nil)
			if err != nil {
				log.Printf("Error executing %v", err)
				goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Error executing HTML template", w) // Error 500
				return
			}
			return
		} else {
			goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Error executing HTML template", w) // Error 500
			return
		}
	}

	var err error
	if apiReadTime.Before(time.Now().Add(-5 * time.Minute)) {
		fmt.Println("Loading API")
		err = readAPI(w)
		apiReadTime = time.Now()
	}
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

	if tmplIndex != nil {
		err = tmplIndex.Execute(w, data)
		if err != nil {
			fmt.Println("Trying to execute home page")
			log.Printf("Error executing %v", err)
			goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Error executing HTML template", w) // Error 500
			return
		}
	} else {
		goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Error executing HTML template", w) // Error 500
		return
	}
}

// artistHandler serves a site for a specific artist
func artistHandler(w http.ResponseWriter, r *http.Request) {
	// Allow only GET method
	if r.Method != http.MethodGet {
		goToErrorPage(http.StatusMethodNotAllowed, "Method Not Allowed", "Only GET method is allowed", w) // Error 405
		return
	}

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

	if artistID > 0 && artistID <= len(artInfos) {
		dataAP.Artist = artInfos[artistID-1]
	} else {
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

	if tmplArtist != nil {
		err = tmplArtist.Execute(w, dataAP)
		if err != nil {
			log.Printf("Error executing %v", err)
			goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Error executing HTML template", w) // Error 500
			return
		}
	} else {
		goToErrorPage(http.StatusInternalServerError, "Internal Server Error", "Error executing HTML template", w) // Error 500
		return
	}
}

// goToErrorPage handles errors by loading an error page to the user
func goToErrorPage(errorN int, m1 string, m2 string, w http.ResponseWriter) {

	fmt.Printf("%d %s, %s\n", errorN, m1, m2)

	if tmplError != nil {

		w.WriteHeader(errorN)
		epd := errorPageData{uint(errorN), m1, m2}

		err := tmplError.Execute(w, epd)
		if err != nil {
			log.Printf("Error executing template 'errorpage.html': %v", err)
			fmt.Fprintf(w, "%d %s\n%s", errorN, m1, m2) // Error 500 plaintext
			return
		}
	} else {
		http.Error(w, strconv.Itoa(errorN)+" "+m1+": "+m2, errorN)
		return
	}
}

func main() {

	apiReadTime = time.Now().Add(-6 * time.Minute)

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
