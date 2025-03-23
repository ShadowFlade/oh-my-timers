package main

import (
	"log"
	"net/http"
	"shadowflade/timers/global"
	"shadowflade/timers/pkg/handlers"
)

func main() {
	timerHandler := handlers.TimerHandler{}
	fileServer := http.FileServer(http.Dir("./assets/"))
	mux := http.NewServeMux()
	mux.Handle("/createTimer", http.HandlerFunc(timerHandler.CreateTimer))
	mux.Handle("/", http.HandlerFunc(timerHandler.RenderUserTimers))
	mux.Handle("/assets/", http.StripPrefix("/assets", fileServer))

	log.Fatal(http.ListenAndServe(":"+global.PORT, mux))
}
