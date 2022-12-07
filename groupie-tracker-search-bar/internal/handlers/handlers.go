package handlers

import (
	"html/template"
	"net/http"

	checkRequest "groupie/internal/checkRequest"
	"groupie/internal/errors"
	"groupie/internal/parseJson"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	statusCheckMethod := checkRequest.CheckMethod(r)
	if statusCheckMethod != http.StatusOK {
		check := checkRequest.CheckStatus(statusCheckMethod)
		w.WriteHeader(statusCheckMethod)
		errors.CheckErrors(w, check)
		return
	}

	statusCheckPath := checkRequest.CheckPath(r, "/")
	if statusCheckPath != http.StatusOK {
		check := checkRequest.CheckStatus(statusCheckPath)
		w.WriteHeader(statusCheckPath)
		errors.CheckErrors(w, check)
		return
	}

	statusCheckParseFiles := checkRequest.CheckArtists(w, "./ui/templates/index.html")
	if statusCheckParseFiles != http.StatusOK {
		check := checkRequest.CheckStatus(statusCheckParseFiles)
		w.WriteHeader(statusCheckParseFiles)
		errors.CheckErrors(w, check)
		return
	}
}

func PageTwo(w http.ResponseWriter, r *http.Request) {
	statusCheckMethod := checkRequest.CheckMethod(r)
	if statusCheckMethod != http.StatusOK {
		check := checkRequest.CheckStatus(statusCheckMethod)
		w.WriteHeader(statusCheckMethod)
		errors.CheckErrors(w, check)
		return
	}

	statusCheckParseFiles := checkRequest.CheckArtist(w, r, "./ui/templates/pageTwo.html")
	if statusCheckParseFiles != http.StatusOK {
		check := checkRequest.CheckStatus(statusCheckParseFiles)
		w.WriteHeader(statusCheckParseFiles)
		errors.CheckErrors(w, check)
		return
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		if r.Method == http.MethodGet {
			check := checkRequest.CheckStatus(http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
			errors.CheckErrors(w, check)
			return
		}
		check := checkRequest.CheckStatus(http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		errors.CheckErrors(w, check)
		return
	}

	statusErr := serchArtists(w, r)
	if statusErr != http.StatusOK {
		check := checkRequest.CheckStatus(statusErr)
		w.WriteHeader(http.StatusMethodNotAllowed)
		errors.CheckErrors(w, check)
	}
}

func serchArtists(w http.ResponseWriter, r *http.Request) int {
	tmpl, _ := template.ParseFiles("./ui/templates/index.html")
	search_word := r.FormValue("search-bar")
	res := Search(search_word)
	var Alls parseJson.All
	Alls.AllArtists = parseJson.Artists
	Alls.FoundArtists = res

	err := tmpl.Execute(w, Alls)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
