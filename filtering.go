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
	recPerf   [2]int
	band      bool
	solo      bool
	countries []bool
}

var (
	minmaxFirst [6]int
)

// getMinMax retrieves the minimun and maximum values for three ranges in the filter
func getMinMax() (int, int, int, int, int, int) {
	formMin, formMax, fAMin, fAMax, peMin, peMax := 1950, 2024, 1950, 2024, 1950, 2024
	if len(artInfos) > 0 {
		formMin, formMax, fAMin, fAMax = artInfos[0].CreDate, artInfos[0].CreDate, artInfos[0].FirstAlbum.Year(), artInfos[0].FirstAlbum.Year()
		peMin, peMax = artInfos[0].Gigs[0].Date.Year(), artInfos[0].Gigs[0].Date.Year()
	}
	for _, ai := range artInfos {
		if ai.CreDate < formMin {
			formMin = ai.CreDate
		}
		if ai.CreDate > formMax {
			formMax = ai.CreDate
		}
		if ai.FirstAlbum.Year() < fAMin {
			fAMin = ai.FirstAlbum.Year()
		}
		if ai.FirstAlbum.Year() > fAMax {
			fAMax = ai.FirstAlbum.Year()
		}
		for _, gig := range ai.Gigs {
			if gig.Date.Year() < peMin {
				peMin = ai.Gigs[0].Date.Year()
			}
			if gig.Date.Year() > peMax {
				peMax = ai.Gigs[0].Date.Year()
			}
		}
	}
	return formMin, formMax, fAMin, fAMax, peMin, peMax
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
	formMin, formMax, fAMin, fAMax, peMin, peMax := getMinMax()
	minmaxFirst = [6]int{formMin, formMax, fAMin, fAMax, peMin, peMax}

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
