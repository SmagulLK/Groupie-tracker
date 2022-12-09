package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"groupie/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomePage)
	mux.HandleFunc("/pageTwo/", handlers.PageTwo)
	mux.HandleFunc("/search", handlers.SearchHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))

	log.Print("Start server http://127.0.0.1:8000")
	go func() {
		log.Fatal(http.ListenAndServe(":8000", mux))
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Println("Shutting down")
}
