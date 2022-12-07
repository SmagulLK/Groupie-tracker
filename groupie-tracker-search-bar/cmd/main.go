package main

import (
	"log"
	"net/http"

	"groupie/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomePage)
	mux.HandleFunc("/pageTwo/", handlers.PageTwo)
	mux.HandleFunc("/search", handlers.SearchHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))

	log.Print("Start server http://127.0.0.1:4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
