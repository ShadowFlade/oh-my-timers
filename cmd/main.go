package main

import (
	"log"
	"net/http"
	"shadowflade/timers/pkg/global"
	"shadowflade/timers/pkg/handlers"
	"time"
)

func main() {
	timerHandler := handlers.TimerHandler{}
	fileServer := http.FileServer(http.Dir("./assets/"))

	mux := http.NewServeMux()
	mux.HandleFunc("/createTimer", timerHandler.CreateTimer)
	mux.HandleFunc("/pauseTimer", timerHandler.PauseTimer)
	mux.HandleFunc("/startTimer", timerHandler.StartTimer)
	mux.HandleFunc("/deleteTimer", timerHandler.DeleteTimer)
	mux.HandleFunc("/createUser", timerHandler.CreateUser)
	mux.HandleFunc("/updateTimerTittle", timerHandler.UpdateTimerTitle)
	mux.HandleFunc("/", timerHandler.RenderUserTimers)
	mux.Handle("/assets/", http.StripPrefix("/assets", fileServer))

	log.Fatal(http.ListenAndServe(":"+global.PORT, mux))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	time.Local = time.UTC
}
