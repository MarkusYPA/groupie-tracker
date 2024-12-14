package main

import (
	"net/http"
	"strconv"
)

// Contains user selections
type filter struct {
	order          string
	created        [2]int
	firstAl        [2]int
	recShow        int
	numbsOfMembers []bool
	countries      []bool
	locales        []bool
}

var (
	minmaxLimits [6]int
)

// getMinMaxLimits retrieves the minimun and maximum values for three ranges in the filter
func getMinMaxLimits() [6]int {
	startMin, startMax, albumMin, albumMax, showMin, showMax := 1900, 2050, 1900, 2050, 1900, 2050
	if len(artInfos) > 0 {
		startMin, startMax, albumMin, albumMax = artInfos[0].StartDate, artInfos[0].StartDate, artInfos[0].FirstAlbum.Year(), artInfos[0].FirstAlbum.Year()
		showMin, showMax = artInfos[0].Gigs[0].Date.Year(), artInfos[0].Gigs[0].Date.Year()
	}
	for _, ai := range artInfos {
		if ai.StartDate < startMin {
			startMin = ai.StartDate
		}
		if ai.StartDate > startMax {
			startMax = ai.StartDate
		}
		if ai.FirstAlbum.Year() < albumMin {
			albumMin = ai.FirstAlbum.Year()
		}
		if ai.FirstAlbum.Year() > albumMax {
			albumMax = ai.FirstAlbum.Year()
		}
		for _, gig := range ai.Gigs {
			if gig.Date.Year() < showMin {
				showMin = ai.Gigs[0].Date.Year()
			}
			if gig.Date.Year() > showMax {
				showMax = ai.Gigs[0].Date.Year()
			}
		}
	}

	return [6]int{startMin, startMax, albumMin, albumMax, showMin, showMax}
}

// defaultFilter sets the filter values to default
func defaultFilter() filter {
	countries := make([]bool, len(allCountries))
	for i := range allCountries {
		countries[i] = true
	}

	locales := make([]bool, len(allLocales))
	for i := range allLocales {
		locales[i] = true
	}

	memNumbs := make([]bool, len(allMemberNumbers))
	for i := range allMemberNumbers {
		memNumbs[i] = true
	}

	ord := "namedown"
	numbsOfMembers := memNumbs
	minmaxLimits = getMinMaxLimits()

	return filter{
		order:          ord,
		created:        [2]int{minmaxLimits[0], minmaxLimits[1]},
		firstAl:        [2]int{minmaxLimits[2], minmaxLimits[3]},
		recShow:        minmaxLimits[5],
		numbsOfMembers: numbsOfMembers,
		countries:      countries,
		locales:        locales,
	}
}

// newFilter places the user's selections to a filter
func newFilter(r *http.Request) filter {
	ord := r.FormValue("order")
	startMin, _ := strconv.Atoi(r.FormValue("startmin"))
	startMax, _ := strconv.Atoi(r.FormValue("startmax"))
	if startMax < startMin {
		startMax = startMin
	}
	albumMin, _ := strconv.Atoi(r.FormValue("albummin"))
	albumMax, _ := strconv.Atoi(r.FormValue("albummax"))
	if albumMax < albumMin {
		albumMax = albumMin
	}
	showMax, _ := strconv.Atoi(r.FormValue("showmax"))

	selectedCountries := []string{}
	countries := make([]bool, len(allCountries))
	for i, c := range allCountries {
		countries[i] = (r.FormValue(c) == "on" || r.Method == http.MethodGet)
		if countries[i] {
			selectedCountries = append(selectedCountries, c)
		}
	}

	locales := make([]bool, len(allLocales))
	for i, l := range allLocales {
		locales[i] = (r.FormValue(l) == "on" || r.Method == http.MethodGet)
		if !isLocaleInCountries(l, selectedCountries) { // unselect places that aren't displayed so they don't interfere with user's selections
			locales[i] = false
		}
	}

	// Get values from form about what member numbers are selected
	memNumbs := make([]bool, len(allMemberNumbers))
	for i := range allMemberNumbers {
		memNumbs[i] = (r.FormValue(strconv.Itoa(i+1)) == "on" || r.Method == http.MethodGet) // Name checkboxes just numbers?
	}

	return filter{
		order:          ord,
		created:        [2]int{startMin, startMax},
		firstAl:        [2]int{albumMin, albumMax},
		recShow:        showMax,
		numbsOfMembers: memNumbs,
		countries:      countries,
		locales:        locales,
	}
}

// compare compares two artistInfos according to an attribute specified in string s
func compare(a1, a2 artistInfo, s string) bool {
	switch s {
	case "namedown":
		return a1.Name <= a2.Name
	case "startdown":
		return a1.StartDate <= a2.StartDate
	case "albumdown":
		return a1.FirstAlbum.Before(a2.FirstAlbum)
	case "showdown":
		return a1.Gigs[0].Date.Before(a2.Gigs[0].Date)
	case "nameup":
		return a1.Name > a2.Name
	case "startup":
		return a1.StartDate > a2.StartDate
	case "albumup":
		return a1.FirstAlbum.After(a2.FirstAlbum)
	case "showup":
		return a1.Gigs[0].Date.After(a2.Gigs[0].Date)
	}
	return true
}

// sortArtists sorts a slice of artistInfo according to the instruction in string ord
func sortArtists(as *[]artistInfo, ord string) {
	for i := 0; i < len(*as)-1; i++ {
		for j := i + 1; j < len(*as); j++ {
			if !compare((*as)[i], (*as)[j], ord) {
				(*as)[i], (*as)[j] = (*as)[j], (*as)[i]
			}
		}
	}
}

// isLocaleInCountries tells if a locale exists in one of the countries on the given list
func isLocaleInCountries(locale string, selCountries []string) bool {
	for _, cou := range selCountries {
		for _, pair := range allCountryLocalePairs {
			if pair[0] == cou && pair[1] == locale {
				return true
			}
		}
	}
	return false
}

// filterBy returns artistInfos that pass the filter, sorted by the rule in it
func filterBy(fil filter, arInfos []artistInfo) []artistInfo {
	aisOut := []artistInfo{}

	selectedCountries := []string{}
	for i := 0; i < len(allCountries); i++ {
		if fil.countries[i] {
			selectedCountries = append(selectedCountries, allCountries[i])
		}
	}

	selectedLocales := []string{}
	for i := 0; i < len(allLocales); i++ {
		if fil.locales[i] {
			selectedLocales = append(selectedLocales, allLocales[i])
		}
	}

	selectedMemberNums := []int{}
	for i := 0; i < len(allMemberNumbers); i++ {
		if fil.numbsOfMembers[i] {
			selectedMemberNums = append(selectedMemberNums, allMemberNumbers[i])
		}
	}

	for _, ai := range arInfos {
		if ai.StartDate < fil.created[0] || ai.StartDate > fil.created[1] {
			continue
		}
		if ai.FirstAlbum.Year() < fil.firstAl[0] || ai.FirstAlbum.Year() > fil.firstAl[1] {
			continue
		}
		if ai.Gigs[0].Date.Year() > fil.recShow {
			continue
		}

		foundCountry := false
		for _, cn := range selectedCountries {
			for _, g := range ai.Gigs {
				if g.Country == cn {
					foundCountry = true
					break // from inner loop
				}
			}
			if foundCountry {
				break // from outer loop
			}
		}

		foundLocale := false
		for _, ln := range selectedLocales {
			for _, g := range ai.Gigs {
				if g.Locale == ln {
					foundLocale = true
					break // from inner loop
				}
			}
			if foundLocale {
				break // from outer loop
			}
		}

		foundMembers := false
		for _, mn := range selectedMemberNums {
			if len(ai.Members) == mn {
				foundMembers = true
			}
		}

		if foundCountry && foundLocale && foundMembers {
			aisOut = append(aisOut, ai)
		}
	}
	sortArtists(&aisOut, fil.order)
	return aisOut
}
