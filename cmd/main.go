package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofor-little/env"

	"shadowflade/timers/pkg/global"
	"shadowflade/timers/pkg/handlers"
)

func main() {
	timerHandler := handlers.TimerHandler{}

	mux := http.NewServeMux()
	fmt.Println("slkdjflsdkjf")
	mux.HandleFunc("/assets/", assetsHandler)
	mux.HandleFunc("/createTimer", timerHandler.CreateTimer)
	mux.HandleFunc("/pauseTimer", timerHandler.PauseTimer)
	mux.HandleFunc("/stopTimer", timerHandler.StopTimer)
	mux.HandleFunc("/startTimer", timerHandler.StartTimer)
	mux.HandleFunc("/deleteTimer", timerHandler.DeleteTimer)
	mux.HandleFunc("/createUser", timerHandler.CreateUser)
	mux.HandleFunc("/updateTimerTittle", timerHandler.UpdateTimerTitle)
	mux.HandleFunc("/updateTimerColor", timerHandler.AddUpdateTimerColor)
	mux.HandleFunc("/", timerHandler.RenderUserTimers)

	log.Fatal(http.ListenAndServe(":"+global.PORT, mux))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	time.Local = time.UTC
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	if err := env.Load("./.env"); err != nil {
		fmt.Println("error")
		panic(err)
	}

	path := strings.TrimPrefix(r.URL.Path, "/assets/")

	if path == "" || strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}

	if strings.Contains(path, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	default:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}

	// Кеширование
	cacheControlHeader := ""
	if env.Get("IS_PROD", "Y") == "Y" {
		cacheControlHeader = "public, max-age=" + env.Get("MAX_AGE_CACHE_TIME", "86400")
	} else {
		cacheControlHeader = "no-cache, no-store, must-revalidate"
	}
	if ext == ".css" || ext == ".js" {
		w.Header().Set("Cache-Control", cacheControlHeader)
	}

	http.ServeFile(w, r, "./assets/"+path)
}
