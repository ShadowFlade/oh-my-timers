package main

import (
	"log"
	"net/http"
	"shadowflade/timers/global"
	"shadowflade/timers/pkg/handlers"
)

func main() {
	timerHandler := handlers.TimerHandler{}
	http.Handle("/createTimer", http.HandlerFunc(timerHandler.Create))
	http.Handle("/", http.HandlerFunc(timerHandler.RenderUserTimers))
	log.Fatal(http.ListenAndServe(":"+global.PORT, nil))
}
