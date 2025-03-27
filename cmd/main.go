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
	mux.HandleFunc("/createTimer", timerHandler.CreateTimer)
	mux.HandleFunc("/", timerHandler.RenderUserTimers)
	mux.Handle("/assets/", http.StripPrefix("/assets", fileServer))

	log.Fatal(http.ListenAndServe(":"+global.PORT, mux))
}
