package main

import (
	"net/http"
	"strconv"
)

// Contains user selections
type filter struct {
	order     string
	created   [2]int
	firstAl   [2]int
	recShow   [2]int
	band      bool
	solo      bool
	countries []bool
}

var (
	minmaxLimits [6]int
)

// getMinMaxLimits retrieves the minimun and maximum values for three ranges in the filter
func getMinMaxLimits() [6]int {
	startMin, startMax, albumMin, albumMax, showMin, showMax := 1950, 2024, 1950, 2024, 1950, 2024
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

	ord := "namedown"
	showBand := true
	showSolo := true
	minmaxLimits = getMinMaxLimits()

	return filter{
		order:     ord,
		created:   [2]int{minmaxLimits[0], minmaxLimits[1]},
		firstAl:   [2]int{minmaxLimits[2], minmaxLimits[3]},
		recShow:   [2]int{minmaxLimits[4], minmaxLimits[5]},
		band:      showBand,
		solo:      showSolo,
		countries: countries,
	}
}

// newFilter places the user's selections to a filter
func newFilter(r *http.Request) filter {
	ord := r.FormValue("order")
	showBand := r.FormValue("band") == "on"
	showSolo := r.FormValue("solo") == "on"
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
	showMin, _ := strconv.Atoi(r.FormValue("showmin"))
	showMax, _ := strconv.Atoi(r.FormValue("showmax"))
	if showMax < showMin {
		showMax = showMin
	}

	countries := make([]bool, len(allCountries))
	for i, c := range allCountries {
		countries[i] = (r.FormValue(c) == "on" || r.Method == http.MethodGet)
	}

	return filter{
		order:     ord,
		created:   [2]int{startMin, startMax},
		firstAl:   [2]int{albumMin, albumMax},
		recShow:   [2]int{showMin, showMax},
		band:      showBand,
		solo:      showSolo,
		countries: countries,
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

// filterBy returns artistInfos that pass the filter, sorted by the rule in it
func filterBy(fil filter, arInfos []artistInfo) []artistInfo {
	aisOut := []artistInfo{}
	for _, ai := range arInfos {
		passes := true
		if ai.StartDate < fil.created[0] || ai.StartDate > fil.created[1] {
			passes = false
		}
		if ai.FirstAlbum.Year() < fil.firstAl[0] || ai.FirstAlbum.Year() > fil.firstAl[1] {
			passes = false
		}
		if ai.Gigs[0].Date.Year() < fil.recShow[0] || ai.Gigs[0].Date.Year() > fil.recShow[1] {
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
				break // from outer loop
			}
		}

		if passes && found {
			aisOut = append(aisOut, ai)
		}
	}
	sortArtists(&aisOut, fil.order)
	return aisOut
}
