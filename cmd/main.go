package cmd

import (
	"log"
	"net/http"
	"shadowflade/timers/global"
	"shadowflade/timers/pkg/handlers"
	"html/template"
)

func main() {

	timerHandler := handlers.TimerHandler{}
	http.Handle("/createTimer", http.HandlerFunc(timerHandler.Create))
	http.Handle("/", http.HandlerFunc()
	log.Fatal(http.ListenAndServe(":"+global.PORT, nil))
}
